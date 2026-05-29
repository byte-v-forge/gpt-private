package appsvc

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/byte-v-forge/common-lib/hashx"
	"github.com/byte-v-forge/common-lib/httpjson"
	"github.com/byte-v-forge/common-lib/randx"
	"github.com/byte-v-forge/common-lib/stringx"
	"github.com/byte-v-forge/gpt-private/gopay/pb"
	gopayapp "github.com/byte-v-forge/gpt-private/gopay/protocol/app"
)

type Server struct {
	pb.UnimplementedGopayAppServiceServer
	cfg   Config
	store *StateStore
}

func NewServer(cfg Config) (*Server, error) {
	store, err := NewStateStore(context.Background(), cfg.StateRedisURL, cfg.StateKeyPrefix, cfg.StateTTL)
	if err != nil {
		return nil, err
	}
	return &Server{cfg: cfg, store: store}, nil
}

func (s *Server) Close() error {
	if s == nil || s.store == nil {
		return nil
	}
	return s.store.Close()
}

func (s *Server) GetGoPayState(ctx context.Context, req *pb.GetGoPayStateRequest) (*pb.GetGoPayStateResponse, error) {
	key, err := NormalizeStateKey(req.GetUserId())
	if err != nil {
		return &pb.GetGoPayStateResponse{Success: false, ErrorMessage: err.Error()}, nil
	}
	raw, err := s.store.Load(ctx, key)
	if err != nil {
		return &pb.GetGoPayStateResponse{Success: false, ErrorMessage: err.Error()}, nil
	}
	state := s.parseRequestState(raw)
	s.bindGoPayAccountIdentity(state, key)
	raw = stateJSON(state)
	return &pb.GetGoPayStateResponse{Success: true, UserId: key, StateJson: raw}, nil
}

func (s *Server) UpsertGoPayState(ctx context.Context, req *pb.UpsertGoPayStateRequest) (*pb.UpsertGoPayStateResponse, error) {
	key, err := NormalizeStateKey(req.GetUserId())
	if err != nil {
		return &pb.UpsertGoPayStateResponse{Success: false, ErrorMessage: err.Error()}, nil
	}
	state := s.parseRequestState(stringx.FirstNonEmpty(req.GetStateJson(), "{}"))
	s.bindGoPayAccountIdentity(state, key)
	raw, err := s.store.Save(ctx, key, stateJSON(state))
	if err != nil {
		return &pb.UpsertGoPayStateResponse{Success: false, ErrorMessage: err.Error()}, nil
	}
	return &pb.UpsertGoPayStateResponse{Success: true, UserId: key, StateJson: raw}, nil
}

func (s *Server) DeleteGoPayState(ctx context.Context, req *pb.DeleteGoPayStateRequest) (*pb.DeleteGoPayStateResponse, error) {
	key, err := NormalizeStateKey(req.GetUserId())
	if err != nil {
		return &pb.DeleteGoPayStateResponse{Success: false, ErrorMessage: err.Error()}, nil
	}
	if err := s.store.Delete(ctx, key); err != nil {
		return &pb.DeleteGoPayStateResponse{Success: false, ErrorMessage: err.Error()}, nil
	}
	return &pb.DeleteGoPayStateResponse{Success: true}, nil
}

func (s *Server) parseRequestState(raw string) stateMap {
	state, err := parseState(raw)
	if err != nil {
		return stateMap{"last_error": err.Error()}
	}
	return state
}

func (s *Server) authBody(extra map[string]any) map[string]any {
	body := map[string]any{}
	for key, value := range extra {
		body[key] = value
	}
	body["client_id"] = s.cfg.GotoClientID
	body["client_secret"] = s.cfg.GotoClientSecret
	return body
}

func (s *Server) pin(value string) string {
	return strings.TrimSpace(value)
}

func (s *Server) signupProfile(phone, name, email string) (string, string) {
	resolvedName := strings.TrimSpace(name)
	resolvedEmail := strings.TrimSpace(email)
	if resolvedName != "" {
		return resolvedName, resolvedEmail
	}
	return signupNameFromSeed(signupSeed(phone)), resolvedEmail
}

func (s *Server) signupBasicAuthorization() string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(s.cfg.SignupAuthUUID))
}

func (s *Server) newClient(ctx context.Context, token string, proxyURL string, device gopayapp.DeviceFingerprint) (*gopayapp.ClientSet, error) {
	cfg := gopayapp.ConfigFromEnv(token)
	cfg.ProxyURL = proxyURL
	cfg.Timeout = 30 * time.Second
	cfg.Device = device
	cfg.Logger = func(ctx context.Context, message string, fields map[string]any) {
		fmt.Printf("[gopay-app] %s %v\n", message, fields)
	}
	return gopayapp.NewClientSet(cfg)
}

func (s *Server) clientForState(ctx context.Context, state stateMap) (*gopayapp.ClientSet, error) {
	refresh := s.ensureAccessToken(ctx, state, s.cfg.TokenRefreshMinTTL, false)
	if !anyBool(refresh["success"]) && !tokenUsable(state, "token", 0) {
		return nil, fmt.Errorf("%s", stringx.FirstNonEmpty(anyString(refresh["error"]), "token refresh failed"))
	}
	if err := s.ensureProxyRuntimeSession(ctx, state, proxyRuntimeAcquireOptions{}); err != nil {
		return nil, err
	}
	device, err := s.ensureDevice(state)
	if err != nil {
		return nil, err
	}
	return s.newClient(ctx, stateString(state, "token"), s.proxyForState(state), device)
}

func (s *Server) tmpClientForState(ctx context.Context, state stateMap) (*gopayapp.ClientSet, error) {
	token := stateString(state, "_tmp_token")
	if token == "" {
		return nil, fmt.Errorf("temporary account token missing")
	}
	if !tmpTokenUsable(state, 0) {
		expiresAt := firstNonZero(jwtExpiresAt(token), stateInt(state, "_tmp_token_expires_at"))
		return nil, fmt.Errorf("temporary account token expired: expires_at=%d", expiresAt)
	}
	if err := s.ensureProxyRuntimeSession(ctx, state, proxyRuntimeAcquireOptions{}); err != nil {
		return nil, err
	}
	device, err := s.ensureDevice(state)
	if err != nil {
		return nil, err
	}
	return s.newClient(ctx, token, s.proxyForState(state), device)
}

func (s *Server) rotateLoginAttemptIdentity(ctx context.Context, state stateMap) error {
	if state == nil {
		return nil
	}
	if err := s.ensureProxyRuntimeSession(ctx, state, proxyRuntimeAcquireOptions{ForceNew: true}); err != nil {
		return err
	}
	state["_proxy_runtime_session_rotated_for_login"] = true
	return nil
}

func (s *Server) proxyForState(state stateMap) string {
	return stateString(state, "_gopay_proxy")
}

func (s *Server) generateDeviceProxyState(ctx context.Context, accountID string, countryCode string, forceNew bool, skipPreflight bool, ephemeralProfile bool) (stateMap, error) {
	identity := strings.TrimSpace(accountID)
	if identity == "" {
		identity = "local"
	}
	state := stateMap{}
	if ephemeralProfile {
		identity = "ephemeral:" + randomProfileID()
		state["_gopay_profile_ephemeral"] = true
	} else {
		state = s.loadAccountProfile(ctx, identity)
	}
	s.bindGoPayAccountIdentity(state, identity)
	state["_gopay_country_code"] = normalizeGoPayProxyCountryCode(countryCode)
	var err error
	if ephemeralProfile {
		_, err = ensureRandomDevice(state)
	} else {
		_, err = s.ensureDevice(state)
	}
	if err != nil {
		return state, err
	}
	if err := s.ensureProxyRuntimeSession(ctx, state, proxyRuntimeAcquireOptions{AccountID: identity, CountryCode: countryCode, ForceNew: forceNew, SkipPreflight: skipPreflight}); err != nil {
		return state, err
	}
	if !ephemeralProfile {
		_ = s.saveAccountProfile(ctx, identity, state)
	}
	return state, nil
}

func (s *Server) loadAccountProfile(ctx context.Context, identity string) stateMap {
	key := accountProfileStateKey(identity)
	if key == "" || s.store == nil {
		return stateMap{}
	}
	raw, err := s.store.Load(ctx, key)
	if err != nil {
		return stateMap{}
	}
	state, err := parseState(raw)
	if err != nil {
		return stateMap{}
	}
	return state
}

func (s *Server) saveAccountProfile(ctx context.Context, identity string, state stateMap) error {
	key := accountProfileStateKey(identity)
	if key == "" || s.store == nil {
		return nil
	}
	profile := stateMap{}
	if device := nestedMap(state["device"]); len(device) > 0 {
		profile["device"] = device
	}
	if accountID := stateString(state, "_gopay_account_id"); accountID != "" {
		profile["_gopay_account_id"] = accountID
	}
	if countryCode := stateString(state, "_gopay_country_code"); countryCode != "" {
		profile["_gopay_country_code"] = countryCode
	}
	if proxyAccountID := stateString(state, "_proxy_runtime_account_id"); proxyAccountID != "" {
		profile["_proxy_runtime_account_id"] = proxyAccountID
	}
	if len(profile) == 0 {
		return nil
	}
	_, err := s.store.Save(ctx, key, stateJSON(profile))
	return err
}

func accountProfileStateKey(identity string) string {
	identity = strings.TrimSpace(identity)
	if identity == "" {
		return ""
	}
	return "profile:" + hashx.ShortSHA256(identity, 24)
}

func (s *Server) deviceProxyDiagnostics(state stateMap) map[string]any {
	data := map[string]any{}
	proxyURL := stateString(state, "_gopay_proxy")
	if proxyURL != "" {
		data["proxy_hash"] = hashx.ShortSHA256(proxyURL, 12)
	}
	if hash := stateString(state, "_proxy_runtime_session_hash"); hash != "" {
		data["proxy_runtime_session_hash"] = hash
		data["proxy_runtime_pool_endpoints"] = anyInt(state["_proxy_runtime_pool_endpoints"])
	}
	if rotated, ok := state["_proxy_runtime_session_rotated"].(bool); ok {
		data["proxy_runtime_session_rotated"] = rotated
	}
	for _, key := range []string{
		"_gopay_account_id",
		"_gopay_country_code",
		"_proxy_runtime_account_id",
		"_proxy_runtime_lease_id",
		"_proxy_runtime_lease_expires_at",
		"_proxy_runtime_listener_id",
		"_proxy_runtime_chain_route",
		"_proxy_runtime_exit_ip_hash",
		"_proxy_runtime_exit_country_code",
		"_proxy_runtime_exit_region",
		"_proxy_runtime_ip_fraud_risk_level",
		"_proxy_runtime_ip_fraud_risk_score",
		"_proxy_runtime_ip_network_kind",
		"_proxy_runtime_ip_anonymizer_kind",
		"_proxy_runtime_connectivity_reachable",
		"_proxy_runtime_connectivity_targets",
		"_proxy_runtime_preflight_skipped",
		"_proxy_runtime_preflight_attempts",
		"_gopay_profile_ephemeral",
	} {
		if value, ok := state[key]; ok {
			data[strings.TrimPrefix(key, "_")] = value
		}
	}
	if fp := deviceFingerprintForState(state); fp != "" {
		data["device_fingerprint"] = fp
	}
	return data
}

func deviceFingerprintForState(state stateMap) string {
	device := nestedMap(state["device"])
	if len(device) == 0 {
		return ""
	}
	out := []string{}
	addPlain := func(label, key string) {
		if value := anyString(device[key]); value != "" {
			out = append(out, label+"="+value)
		}
	}
	addHash := func(label, key string) {
		if value := anyString(device[key]); value != "" {
			out = append(out, label+"#"+hashx.ShortSHA256(value, 12))
		}
	}
	addPlain("profile", "profile_id")
	addPlain("make", "x-phonemake")
	addPlain("model", "x-phonemodel")
	addPlain("os", "x-deviceos")
	addPlain("screen", "m1_screen")
	addPlain("tls", "tls_profile")
	addHash("uid", "x-uniqueid")
	addHash("session", "x-session-id")
	addHash("tx", "transaction-id")
	addHash("d1", "d1")
	addHash("conn", "m1_connection_id")
	addHash("widevine", "m1_widevine_id")
	addHash("wifi", "m1_wifi_mac")
	addHash("ssid", "m1_wifi_ssid")
	addHash("sig", "m1_signature")
	addHash("sig_time", "m1_signature_time")
	addHash("firebase", "m1_firebase_app_instance_id")
	addHash("uuid", "m1_device_uuid")
	addHash("adid", "advertising_id")
	addHash("appset", "app_set_id")
	addHash("devtoken", "x-devicetoken")
	addHash("imei", "x-imei")
	addHash("ip", "x-ipaddress")
	if parsed := deviceFromMap(device); parsed.AppID != "" {
		out = append(out, "x_m1#"+hashx.ShortSHA256(parsed.XM1(), 12))
	}
	return strings.Join(out, "/")
}

func (s *Server) ensureDevice(state stateMap) (gopayapp.DeviceFingerprint, error) {
	raw := nestedMap(state["device"])
	if len(raw) > 0 {
		device := deviceFromMap(raw)
		if deviceNeedsBackfill(device) {
			next, err := gopayapp.NewDeviceFingerprint(gopayapp.DeviceConfigFromEnv())
			if err != nil {
				return gopayapp.DeviceFingerprint{}, err
			}
			device = mergeDevice(device, next)
		}
		state["device"] = deviceToMap(device)
		return device, nil
	}
	device, err := gopayapp.NewDeviceFingerprint(gopayapp.DeviceConfigFromEnv())
	if err != nil {
		return gopayapp.DeviceFingerprint{}, err
	}
	out := deviceToMap(device)
	out["profile_id"] = randomProfileID()
	out["profile_created_at"] = time.Now().Unix()
	state["device"] = out
	return device, nil
}

func ensureRandomDevice(state stateMap) (gopayapp.DeviceFingerprint, error) {
	device, err := gopayapp.NewDeviceFingerprint(gopayapp.DeviceConfig{})
	if err != nil {
		return gopayapp.DeviceFingerprint{}, err
	}
	out := deviceToMap(device)
	out["profile_id"] = randomProfileID()
	out["profile_created_at"] = time.Now().Unix()
	state["device"] = out
	return device, nil
}

func (s *Server) newLogonDevice() (gopayapp.DeviceFingerprint, map[string]any, error) {
	device, err := gopayapp.NewDeviceFingerprint(gopayapp.DeviceConfigFromEnv())
	if err != nil {
		return gopayapp.DeviceFingerprint{}, nil, err
	}
	out := deviceToMap(device)
	out["profile_id"] = randomProfileID()
	out["profile_created_at"] = time.Now().Unix()
	return device, out, nil
}

func randomProfileID() string {
	value, err := randx.Hex(8)
	if err != nil {
		return fmt.Sprintf("%x", time.Now().UnixNano())
	}
	return value
}

func deviceNeedsBackfill(device gopayapp.DeviceFingerprint) bool {
	return device.AppID == "" ||
		device.UniqueID == "" ||
		device.TLSProfileName == "" ||
		device.M1Hardware == "" ||
		device.IMEI == "" ||
		device.IPAddress == "" ||
		device.FirebaseID == "" ||
		device.AdvertisingID == "" ||
		device.AppSetID == "" ||
		device.M1SignatureTime == ""
}

func apiError(label string, resp *httpjson.Response) string {
	if resp == nil {
		return label + ": no response"
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return "AUTH_INVALID"
	}
	return fmt.Sprintf("%s: status %d %s", label, resp.StatusCode, compactErrorDetail(resp.Payload))
}

func responseErrors(resp *httpjson.Response) []any {
	if resp == nil {
		return nil
	}
	for _, source := range []any{resp.Payload["errors"], resp.Data()["errors"]} {
		if items, ok := source.([]any); ok {
			return items
		}
	}
	return nil
}

func responseText(resp *httpjson.Response) string {
	if resp == nil {
		return ""
	}
	return string(resp.Body)
}

func isRateLimited(resp *httpjson.Response) bool {
	if resp == nil {
		return false
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		return true
	}
	for _, err := range responseErrors(resp) {
		text := strings.ToLower(compactErrorDetail(err))
		if strings.Contains(text, "ratelimited") {
			return true
		}
	}
	return false
}

func loginMethodsInvalidUser(resp *httpjson.Response) bool {
	if resp == nil || resp.StatusCode != http.StatusUnauthorized {
		return false
	}
	for _, err := range responseErrors(resp) {
		text := strings.ToLower(compactErrorDetail(err))
		if strings.Contains(text, "invalid user") || strings.Contains(text, "could not find the user") {
			return true
		}
	}
	return false
}

func firstNonZero(values ...int64) int64 {
	for _, value := range values {
		if value != 0 {
			return value
		}
	}
	return 0
}
