package paymentsvc

import (
	"strings"

	"github.com/byte-v-forge/gpt-private/gopay/pb"
)

type credential struct {
	sessionToken string
	accessToken  string
	profile      requestProfile
	hasProfile   bool
}

func requestCredential(value *pb.ChatGPTCredential) credential {
	if value == nil {
		return credential{}
	}
	profile, hasProfile := requestProfileFromProto(value.GetRequestProfile())
	return credential{
		sessionToken: strings.TrimSpace(value.GetSessionToken()),
		accessToken:  strings.TrimSpace(value.GetAccessToken()),
		profile:      profile,
		hasProfile:   hasProfile,
	}
}

func (c credential) empty() bool {
	return c.sessionToken == "" && c.accessToken == ""
}

func (c credential) chatGPTProfile(fallback requestProfile) requestProfile {
	if !c.hasProfile {
		return fallback
	}
	return c.profile.withDefaults(fallback)
}

func requestProfileFromProto(value *pb.ChatGPTRequestProfile) (requestProfile, bool) {
	if value == nil {
		return requestProfile{}, false
	}
	profile := requestProfile{
		Name:           "account",
		ProxyURL:       strings.TrimSpace(value.GetProxyUrl()),
		TLSProfile:     strings.TrimSpace(value.GetTlsProfile()),
		UserAgent:      strings.TrimSpace(value.GetUserAgent()),
		SecCHUA:        strings.TrimSpace(value.GetSecChUa()),
		SecCHPlatform:  strings.TrimSpace(value.GetSecChUaPlatform()),
		AcceptLanguage: strings.TrimSpace(value.GetAcceptLanguage()),
		OAILanguage:    strings.TrimSpace(value.GetOaiLanguage()),
		Locale:         strings.TrimSpace(value.GetLocale()),
		DeviceID:       strings.TrimSpace(value.GetDeviceId()),
		Platform:       strings.TrimSpace(value.GetPlatform()),
	}
	return profile, true
}
