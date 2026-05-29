package appsvc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/byte-v-forge/common-lib/hashx"
	"github.com/byte-v-forge/common-lib/httpx"
	"github.com/byte-v-forge/common-lib/redactx"
	"github.com/byte-v-forge/common-lib/stringx"
)

const (
	goPayProxyPurpose              = "gopay_app"
	goPayProxyCountryCode          = "ID"
	goPayProxyLeaseTTL             = "600s"
	goPayProxyPreflightMaxAttempts = 10
)

var goPayProxyConnectivityTargets = []string{
	"https://accounts.goto-products.com/",
	"https://customer.gopayapi.com/",
	"https://gwa.gopayapi.com/",
	"https://api.gojekapi.com/",
	"https://pin-web-client.gopayapi.com/",
}

type proxyRuntimeAcquireOptions struct {
	AccountID     string
	CountryCode   string
	ForceNew      bool
	SkipPreflight bool
}

type proxyRuntimeLeaseResponse struct {
	Lease struct {
		LeaseID           string               `json:"lease_id"`
		ProviderAccountID string               `json:"provider_account_id"`
		ExpiresAt         string               `json:"expires_at"`
		Session           proxyRuntimeSession  `json:"session"`
		Egress            proxyRuntimeEndpoint `json:"egress"`
		Listener          proxyRuntimeListener `json:"listener"`
		ChainPlan         map[string]any       `json:"chain_plan"`
		ErrorMessage      string               `json:"error_message"`
		Labels            map[string]string    `json:"labels"`
	} `json:"lease"`
	Egress    proxyRuntimeEndpoint `json:"egress"`
	ChainPlan map[string]any       `json:"chain_plan"`
	Pool      struct {
		Endpoints []map[string]any `json:"endpoints"`
	} `json:"pool"`
}

type proxyRuntimeSession struct {
	SessionID  string `json:"session_id"`
	ProviderID string `json:"provider_id"`
}

type proxyRuntimeEndpoint struct {
	ID           string            `json:"id"`
	Protocol     string            `json:"protocol"`
	Host         string            `json:"host"`
	Port         int               `json:"port"`
	SessionID    string            `json:"session_id"`
	ProviderID   string            `json:"provider_id"`
	UpstreamKind string            `json:"upstream_kind"`
	RotationMode string            `json:"rotation_mode"`
	Labels       map[string]string `json:"labels"`
}

type proxyRuntimeListener struct {
	ListenerID string `json:"listener_id"`
	Kind       string `json:"kind"`
	ListenAddr string `json:"listen_addr"`
	Protocol   string `json:"protocol"`
	RouteID    string `json:"route_id"`
	Managed    bool   `json:"managed"`
}

type proxyRuntimeExitIPResponse struct {
	ProxyExitIP struct {
		IP           string `json:"ip"`
		ErrorMessage string `json:"error_message"`
	} `json:"proxy_exit_ip"`
}

type proxyRuntimeGeoResponse struct {
	ProxyExitGeo struct {
		IP           string `json:"ip"`
		CountryCode  string `json:"country_code"`
		Region       string `json:"region"`
		City         string `json:"city"`
		ErrorMessage string `json:"error_message"`
	} `json:"proxy_exit_geo"`
}

type proxyRuntimeFraudResponse struct {
	Check struct {
		IP             string   `json:"ip"`
		NetworkKind    string   `json:"network_kind"`
		AnonymizerKind string   `json:"anonymizer_kind"`
		RiskLevel      string   `json:"risk_level"`
		RiskScore      float64  `json:"risk_score"`
		RiskSignals    []string `json:"risk_signals"`
		CountryCode    string   `json:"country_code"`
		Region         string   `json:"region"`
		City           string   `json:"city"`
		ErrorMessage   string   `json:"error_message"`
	} `json:"check"`
}

type proxyRuntimeConnectivityResponse struct {
	Check struct {
		TargetURL    string `json:"target_url"`
		Host         string `json:"host"`
		Reachable    bool   `json:"reachable"`
		StatusCode   int    `json:"status_code"`
		LatencyMS    int    `json:"latency_ms"`
		ErrorMessage string `json:"error_message"`
	} `json:"check"`
}

func (s *Server) ensureProxyRuntimeSession(ctx context.Context, state stateMap, options proxyRuntimeAcquireOptions) error {
	if state == nil {
		return fmt.Errorf("gopay state missing")
	}
	identity := s.bindGoPayAccountIdentity(state, options.AccountID)
	if identity == "" {
		return fmt.Errorf("gopay account identity missing")
	}
	state["_gopay_country_code"] = normalizeGoPayProxyCountryCode(options.CountryCode)
	if !options.ForceNew && proxyRuntimeLeaseActive(state) && stateString(state, "_gopay_proxy") != "" && stateString(state, "_proxy_runtime_listener_id") != "" {
		return nil
	}
	sessionData, err := s.createProxyRuntimeSession(ctx, state, proxyRuntimeAcquireOptions{AccountID: identity, CountryCode: options.CountryCode, ForceNew: options.ForceNew, SkipPreflight: options.SkipPreflight})
	if err != nil {
		return err
	}
	for key, value := range sessionData {
		state[key] = value
	}
	return nil
}

func (s *Server) createProxyRuntimeSession(ctx context.Context, state stateMap, options proxyRuntimeAcquireOptions) (map[string]any, error) {
	baseURL := proxyRuntimeAPIBase(s.cfg.ProxyRuntimeHTTPAddr)
	if baseURL == "" {
		return nil, fmt.Errorf("PROXY_RUNTIME_HTTP_ADDR is required for GoPay dynamic IP proxy")
	}
	identity := s.bindGoPayAccountIdentity(state, options.AccountID)
	if identity == "" {
		return nil, fmt.Errorf("gopay account identity missing")
	}
	countryCode := normalizeGoPayProxyCountryCode(options.CountryCode)
	if options.SkipPreflight {
		return s.acquireProxyRuntimeSession(ctx, baseURL, state, identity, countryCode, options.ForceNew, 1, true)
	}
	var lastErr error
	for attempt := 1; attempt <= goPayProxyPreflightMaxAttempts; attempt++ {
		attemptData, err := s.acquireAndPreflightProxyRuntimeSession(ctx, baseURL, state, identity, countryCode, options.ForceNew || attempt > 1, attempt)
		if err == nil {
			attemptData["_proxy_runtime_preflight_attempts"] = attempt
			return attemptData, nil
		}
		lastErr = err
	}
	return nil, fmt.Errorf("gopay dynamic IP preflight failed after %d attempts: %w", goPayProxyPreflightMaxAttempts, lastErr)
}

func (s *Server) acquireAndPreflightProxyRuntimeSession(ctx context.Context, baseURL string, state stateMap, identity string, countryCode string, forceNew bool, attempt int) (map[string]any, error) {
	out, listenerID, exitIP, err := s.acquireProxyRuntimeSessionWithListener(ctx, baseURL, state, identity, countryCode, forceNew, attempt, false)
	if err != nil {
		return nil, err
	}
	geo, err := s.checkProxyRuntimeGeo(ctx, baseURL, exitIP)
	if err != nil {
		return nil, err
	}
	if geo.ProxyExitGeo.CountryCode != "" && !strings.EqualFold(geo.ProxyExitGeo.CountryCode, countryCode) {
		return nil, fmt.Errorf("proxy exit country mismatch: got %s want %s", geo.ProxyExitGeo.CountryCode, countryCode)
	}
	fraud, err := s.checkProxyRuntimeFraud(ctx, baseURL, exitIP)
	if err != nil {
		return nil, err
	}
	if proxyIPFraudRiskRejected(fraud.Check.RiskLevel) {
		return nil, fmt.Errorf("proxy IP fraud risk rejected: level=%s score=%.0f", fraud.Check.RiskLevel, fraud.Check.RiskScore)
	}
	connectivity, err := s.checkGoPayProxyRuntimeConnectivity(ctx, baseURL, listenerID)
	if err != nil {
		return nil, err
	}
	out["_proxy_runtime_preflight_country_code"] = countryCode
	out["_proxy_runtime_preflight_attempt"] = attempt
	out["_proxy_runtime_exit_ip_hash"] = hashx.ShortSHA256(exitIP, 12)
	out["_proxy_runtime_exit_country_code"] = geo.ProxyExitGeo.CountryCode
	out["_proxy_runtime_exit_region"] = geo.ProxyExitGeo.Region
	out["_proxy_runtime_ip_fraud_risk_level"] = normalizeProxyRiskLevel(fraud.Check.RiskLevel)
	out["_proxy_runtime_ip_fraud_risk_score"] = fraud.Check.RiskScore
	out["_proxy_runtime_ip_network_kind"] = shortProxyEnum(fraud.Check.NetworkKind)
	out["_proxy_runtime_ip_anonymizer_kind"] = shortProxyEnum(fraud.Check.AnonymizerKind)
	out["_proxy_runtime_connectivity_reachable"] = true
	out["_proxy_runtime_connectivity_targets"] = connectivity
	return out, nil
}

func (s *Server) acquireProxyRuntimeSession(ctx context.Context, baseURL string, state stateMap, identity string, countryCode string, forceNew bool, attempt int, skipPreflight bool) (map[string]any, error) {
	out, _, _, err := s.acquireProxyRuntimeSessionWithListener(ctx, baseURL, state, identity, countryCode, forceNew, attempt, skipPreflight)
	return out, err
}

func (s *Server) acquireProxyRuntimeSessionWithListener(ctx context.Context, baseURL string, state stateMap, identity string, countryCode string, forceNew bool, attempt int, skipPreflight bool) (map[string]any, string, string, error) {
	accountID := proxyRuntimeAccountID(identity)
	leaseReq := map[string]any{
		"account_id": accountID,
		"purpose":    goPayProxyPurpose,
		"force_new":  forceNew,
		"policy": map[string]any{
			"mode":          "PROXY_SESSION_MODE_STICKY",
			"region":        countryCode,
			"sticky_ttl":    goPayProxyLeaseTTL,
			"upstream_kind": "PROXY_UPSTREAM_KIND_DYNAMIC_IP",
			"rotation_mode": "PROXY_ROTATION_MODE_STICKY_SESSION",
			"labels": map[string]string{
				"purpose": "gopay_app",
				"country": countryCode,
			},
		},
		"chain_policy": map[string]any{
			"country_code":                 countryCode,
			"purpose":                      goPayProxyPurpose,
			"strategy":                     "PROXY_CHAIN_STRATEGY_REGION_AWARE",
			"require_dynamic_exit":         true,
			"allow_direct_dynamic_gateway": true,
			"prefer_line_proxy":            true,
		},
	}
	var lease proxyRuntimeLeaseResponse
	if err := postProxyRuntime(ctx, baseURL+"/leases/acquire", leaseReq, &lease, 30*time.Second); err != nil {
		return nil, "", "", err
	}
	egress := lease.Egress
	if egress.Host == "" {
		egress = lease.Lease.Egress
	}
	proxyURL, err := proxyRuntimeProxyURL(egress)
	if err != nil {
		return nil, "", "", err
	}
	listenerID := strings.TrimSpace(lease.Lease.Listener.ListenerID)
	if listenerID == "" {
		return nil, "", "", fmt.Errorf("proxy-runtime returned lease without listener")
	}

	chainPlan := firstMap(lease.ChainPlan, lease.Lease.ChainPlan)
	out := map[string]any{
		"_gopay_proxy":                      proxyURL,
		"_gopay_account_id":                 identity,
		"_gopay_country_code":               countryCode,
		"_proxy_runtime_account_id":         accountID,
		"_proxy_runtime_lease_id":           lease.Lease.LeaseID,
		"_proxy_runtime_provider_account":   lease.Lease.ProviderAccountID,
		"_proxy_runtime_lease_expires_at":   lease.Lease.ExpiresAt,
		"_proxy_runtime_listener_id":        listenerID,
		"_proxy_runtime_listener_kind":      lease.Lease.Listener.Kind,
		"_proxy_runtime_session_started_at": time.Now().Unix(),
		"_proxy_runtime_pool_endpoints":     len(lease.Pool.Endpoints),
		"_proxy_runtime_session_rotated":    forceNew,
		"_proxy_runtime_preflight_skipped":  skipPreflight,
	}
	if lease.Lease.Session.SessionID != "" {
		out["_proxy_runtime_session_hash"] = hashx.ShortSHA256(lease.Lease.Session.SessionID, 12)
	}
	if routeLabel := chainRouteLabel(chainPlan); routeLabel != "" {
		out["_proxy_runtime_chain_route"] = routeLabel
	}
	exitIP := ""
	if !skipPreflight {
		exitIP, err = s.checkProxyRuntimeExitIP(ctx, baseURL, listenerID)
		if err != nil {
			return nil, "", "", err
		}
	}
	return out, listenerID, exitIP, nil
}

func (s *Server) checkProxyRuntimeExitIP(ctx context.Context, baseURL string, listenerID string) (string, error) {
	var parsed proxyRuntimeExitIPResponse
	if err := postProxyRuntime(ctx, baseURL+"/proxy_exit_ip", map[string]any{"listener_id": listenerID}, &parsed, 20*time.Second); err != nil {
		return "", fmt.Errorf("proxy exit ip check: %w", err)
	}
	if parsed.ProxyExitIP.ErrorMessage != "" {
		return "", fmt.Errorf("proxy exit ip check: %s", parsed.ProxyExitIP.ErrorMessage)
	}
	if parsed.ProxyExitIP.IP == "" {
		return "", fmt.Errorf("proxy exit ip check returned empty ip")
	}
	return parsed.ProxyExitIP.IP, nil
}

func (s *Server) checkProxyRuntimeGeo(ctx context.Context, baseURL string, ip string) (proxyRuntimeGeoResponse, error) {
	var parsed proxyRuntimeGeoResponse
	if err := postProxyRuntime(ctx, baseURL+"/proxy_exit_geo", map[string]any{"ip": ip}, &parsed, 20*time.Second); err != nil {
		return parsed, fmt.Errorf("proxy exit geo check: %w", err)
	}
	if parsed.ProxyExitGeo.ErrorMessage != "" {
		return parsed, fmt.Errorf("proxy exit geo check: %s", parsed.ProxyExitGeo.ErrorMessage)
	}
	return parsed, nil
}

func (s *Server) checkProxyRuntimeFraud(ctx context.Context, baseURL string, ip string) (proxyRuntimeFraudResponse, error) {
	var parsed proxyRuntimeFraudResponse
	if err := postProxyRuntime(ctx, baseURL+"/ip_fraud_check", map[string]any{"ip": ip}, &parsed, 25*time.Second); err != nil {
		return parsed, fmt.Errorf("proxy IP fraud check: %w", err)
	}
	if parsed.Check.ErrorMessage != "" && !proxyIPFraudUnsupported(parsed.Check.RiskLevel, parsed.Check.ErrorMessage) {
		return parsed, fmt.Errorf("proxy IP fraud check: %s", parsed.Check.ErrorMessage)
	}
	return parsed, nil
}

func (s *Server) checkGoPayProxyRuntimeConnectivity(ctx context.Context, baseURL string, listenerID string) ([]map[string]any, error) {
	results := make([]map[string]any, 0, len(goPayProxyConnectivityTargets))
	for _, target := range goPayProxyConnectivityTargets {
		var parsed proxyRuntimeConnectivityResponse
		if err := postProxyRuntime(ctx, baseURL+"/target_connectivity_check", map[string]any{"listener_id": listenerID, "target_url": target}, &parsed, 20*time.Second); err != nil {
			return results, fmt.Errorf("target connectivity check %s: %w", target, err)
		}
		result := map[string]any{
			"target_url":  target,
			"reachable":   parsed.Check.Reachable,
			"status_code": parsed.Check.StatusCode,
			"latency_ms":  parsed.Check.LatencyMS,
		}
		if parsed.Check.ErrorMessage != "" {
			result["error_message"] = parsed.Check.ErrorMessage
		}
		results = append(results, result)
		if !parsed.Check.Reachable {
			return results, fmt.Errorf("target connectivity failed: %s", target)
		}
	}
	return results, nil
}

func postProxyRuntime(ctx context.Context, endpoint string, payload map[string]any, out any, timeout time.Duration) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	reqCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	raw, err := httpx.ReadLimited(resp.Body, 1<<20)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("status %d %s", resp.StatusCode, redactx.Snippet(redactx.Text(string(raw)), 300))
	}
	if out == nil {
		return nil
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return fmt.Errorf("parse response: %w", err)
	}
	return nil
}

func proxyRuntimeAPIBase(value string) string {
	value = strings.TrimRight(strings.TrimSpace(value), "/")
	if value == "" {
		return ""
	}
	if strings.HasSuffix(value, "/api/proxy-runtime") || strings.HasSuffix(value, "/proxy") {
		return value
	}
	return value + "/api/proxy-runtime"
}

func proxyRuntimeAccountID(identity string) string {
	identity = strings.TrimSpace(identity)
	if identity == "" {
		return ""
	}
	if strings.HasPrefix(identity, "gopay-app-") {
		return identity
	}
	return "gopay-app-" + hashx.ShortSHA256(identity, 20)
}

func (s *Server) bindGoPayAccountIdentity(state stateMap, explicit string) string {
	identity := strings.TrimSpace(explicit)
	if identity == "" {
		identity = stateString(state, "_gopay_account_id")
	}
	if identity == "" {
		identity = stateString(state, "_gopay_user_id")
	}
	if identity == "" {
		identity = "local"
	}
	if state != nil {
		state["_gopay_account_id"] = identity
		state["_proxy_runtime_account_id"] = proxyRuntimeAccountID(identity)
	}
	return identity
}

func normalizeGoPayProxyCountryCode(value string) string {
	value = strings.ToUpper(strings.TrimSpace(value))
	value = strings.TrimPrefix(value, "+")
	switch value {
	case "", "62", "ID", "IDN", "INDONESIA":
		return goPayProxyCountryCode
	default:
		return value
	}
}

func proxyRuntimeProxyURL(endpoint proxyRuntimeEndpoint) (string, error) {
	if endpoint.Host == "" || endpoint.Port <= 0 {
		return "", fmt.Errorf("proxy-runtime returned invalid lease egress")
	}
	scheme := "http"
	switch endpoint.Protocol {
	case "PROXY_PROTOCOL_SOCKS5", "SOCKS5", "socks5":
		scheme = "socks5"
	}
	return (&url.URL{Scheme: scheme, Host: fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)}).String(), nil
}

func proxyRuntimeLeaseActive(state stateMap) bool {
	if stateString(state, "_proxy_runtime_lease_id") == "" {
		return false
	}
	expiresAt := stateString(state, "_proxy_runtime_lease_expires_at")
	if expiresAt == "" {
		return false
	}
	parsed, err := time.Parse(time.RFC3339Nano, expiresAt)
	if err != nil {
		return false
	}
	return time.Now().Add(30 * time.Second).Before(parsed)
}

func proxyIPFraudRiskRejected(level string) bool {
	switch normalizeProxyRiskLevel(level) {
	case "HIGH", "CRITICAL":
		return true
	default:
		return false
	}
}

func proxyIPFraudUnsupported(level string, message string) bool {
	if normalizeProxyRiskLevel(level) == "UNSUPPORTED" {
		return true
	}
	return strings.Contains(strings.ToLower(strings.TrimSpace(message)), "unsupported")
}

func normalizeProxyRiskLevel(level string) string {
	level = strings.ToUpper(strings.TrimSpace(level))
	level = strings.TrimPrefix(level, "PROXY_IP_FRAUD_RISK_LEVEL_")
	return level
}

func shortProxyEnum(value string) string {
	value = strings.ToUpper(strings.TrimSpace(value))
	for _, prefix := range []string{"PROXY_IP_NETWORK_KIND_", "PROXY_IP_ANONYMIZER_KIND_", "PROXY_PROTOCOL_", "PROXY_UPSTREAM_KIND_", "PROXY_ROTATION_MODE_", "EGRESS_LISTENER_KIND_"} {
		value = strings.TrimPrefix(value, prefix)
	}
	return value
}

func firstMap(values ...map[string]any) map[string]any {
	for _, value := range values {
		if len(value) > 0 {
			return value
		}
	}
	return nil
}

func chainRouteLabel(chainPlan map[string]any) string {
	if len(chainPlan) == 0 {
		return ""
	}
	line := nestedMap(chainPlan["line"])
	dynamicGateway := nestedMap(chainPlan["dynamic_gateway"])
	lineName := stringx.FirstNonEmpty(anyString(line["display_name"]), anyString(line["node_id"]), anyString(line["source_id"]))
	gatewayName := stringx.FirstNonEmpty(anyString(dynamicGateway["display_name"]), anyString(dynamicGateway["gateway_id"]), anyString(dynamicGateway["provider_id"]))
	switch {
	case lineName != "" && gatewayName != "":
		return lineName + " -> " + gatewayName
	case gatewayName != "":
		return gatewayName
	default:
		return lineName
	}
}
