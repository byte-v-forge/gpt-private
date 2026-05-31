package paymentsvc

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/byte-v-forge/common-lib/stringx"
)

const (
	sessionCookieName         = "__Secure-next-auth.session-token"
	sessionCookieFallbackName = "next-auth.session-token"
	sessionCookieChunkSize    = 4096 - 163
)

var sessionCookieChunkSuffixRE = regexp.MustCompile(`^\d+$`)

func chatGPTCookieHeader(sessionToken, deviceID string) string {
	return cookieHeaderWithDeviceID(sessionCookieParts(sessionToken), deviceID)
}

func splitCookieHeader(value string) []string {
	var parts []string
	seen := map[string]bool{}
	for _, raw := range strings.Split(value, ";") {
		part := strings.TrimSpace(raw)
		if part == "" || !strings.Contains(part, "=") {
			continue
		}
		name := strings.TrimSpace(strings.SplitN(part, "=", 2)[0])
		if name == "" || seen[name] {
			continue
		}
		seen[name] = true
		parts = append(parts, part)
	}
	return parts
}

func cookieHeaderWithDeviceID(parts []string, deviceID string) string {
	deviceID = strings.TrimSpace(deviceID)
	if deviceID == "" {
		return strings.Join(parts, "; ")
	}
	out := make([]string, 0, len(parts)+1)
	found := false
	seen := map[string]bool{}
	for _, part := range parts {
		if !strings.Contains(part, "=") {
			continue
		}
		name := strings.TrimSpace(strings.SplitN(part, "=", 2)[0])
		if name == "" || seen[name] {
			continue
		}
		seen[name] = true
		if name == "oai-did" {
			out = append(out, "oai-did="+deviceID)
			found = true
			continue
		}
		out = append(out, part)
	}
	if !found {
		out = append(out, "oai-did="+deviceID)
	}
	return strings.Join(out, "; ")
}

func sessionCookieParts(value string) []string {
	raw := strings.TrimSpace(value)
	if raw == "" {
		return nil
	}
	if strings.HasPrefix(strings.ToLower(raw), "cookie:") {
		raw = strings.TrimSpace(raw[strings.Index(raw, ":")+1:])
	}
	if token := sessionTokenFromJSON(raw); token != "" {
		raw = token
	}
	if strings.Contains(raw, "=") {
		var parts []string
		foundSession := false
		seen := map[string]bool{}
		for _, chunk := range strings.Split(raw, ";") {
			part := strings.Trim(strings.TrimSpace(chunk), `'"`)
			if !strings.Contains(part, "=") {
				continue
			}
			name := strings.TrimSpace(strings.SplitN(part, "=", 2)[0])
			if name == "" || seen[name] {
				continue
			}
			seen[name] = true
			parts = append(parts, part)
			if sessionCookieNameMatches(name) {
				foundSession = true
			}
		}
		if foundSession {
			return parts
		}
	}
	token := strings.Trim(raw, `'"`)
	if token == "" {
		return nil
	}
	if len(token) <= sessionCookieChunkSize {
		return []string{sessionCookieName + "=" + token}
	}
	var out []string
	for idx, offset := 0, 0; offset < len(token); idx, offset = idx+1, offset+sessionCookieChunkSize {
		end := offset + sessionCookieChunkSize
		if end > len(token) {
			end = len(token)
		}
		out = append(out, fmt.Sprintf("%s.%d=%s", sessionCookieName, idx, token[offset:end]))
	}
	return out
}

func sessionTokenFromJSON(value string) string {
	value = strings.TrimSpace(value)
	if !strings.HasPrefix(value, "{") {
		return ""
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(value), &payload); err != nil {
		return ""
	}
	return stringx.FirstNonEmpty(stringAt(payload, "sessionToken"), stringAt(payload, "session_token"))
}

func sessionCookieNameMatches(name string) bool {
	name = strings.TrimSpace(name)
	if name == sessionCookieName || name == sessionCookieFallbackName {
		return true
	}
	for _, base := range []string{sessionCookieName, sessionCookieFallbackName} {
		if strings.HasPrefix(name, base+".") && sessionCookieChunkSuffixRE.MatchString(strings.TrimPrefix(name, base+".")) {
			return true
		}
	}
	return false
}
