package paymentsvc

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/byte-v-forge/common-lib/stringx"
)

func (s *Server) probeTier(ctx context.Context, cred credential) (tierProbe, error) {
	if cred.accessToken != "" {
		return s.probeTierAccessToken(ctx, cred.accessToken, cred), nil
	}
	return s.probePlusActiveSessionToken(ctx, cred.sessionToken, cred), nil
}

func (s *Server) probePlusTrial(ctx context.Context, cred credential) (trialProbe, error) {
	var sessionProbe tierProbe
	if cred.sessionToken != "" {
		sessionProbe = s.probePlusActiveSessionToken(ctx, cred.sessionToken, cred)
		if sessionProbe.Checked && sessionProbe.PlusActive {
			return trialProbe{Checked: true, PlusTrialEligible: false, PlusActive: true, Source: sessionProbe.Source, PlanType: stringx.FirstNonEmpty(sessionProbe.Tier, sessionProbe.PlanType)}, nil
		}
		if sessionProbe.ErrorMessage != "" && cred.accessToken == "" {
			return trialProbe{Checked: sessionProbe.Checked, PlusActive: false, Source: sessionProbe.Source, PlanType: stringx.FirstNonEmpty(sessionProbe.Tier, sessionProbe.PlanType), ErrorMessage: sessionProbe.ErrorMessage}, nil
		}
	}
	ch, err := s.newCharger(ctx, cred, "", "", "", defaultTokenization)
	if err != nil {
		return trialProbe{}, err
	}
	defer ch.close()
	result, err := ch.probePlusTrialCheckout(ctx)
	if err != nil {
		return trialProbe{}, err
	}
	result.PlusActive = sessionProbe.PlusActive
	result.PlanType = stringx.FirstNonEmpty(sessionProbe.Tier, sessionProbe.PlanType)
	return result, nil
}

func (s *Server) probeTierAccessToken(ctx context.Context, accessToken string, cred credential) tierProbe {
	if strings.TrimSpace(accessToken) == "" {
		return tierProbe{Source: "wham_usage", ErrorMessage: "access_token is required"}
	}
	profile := cred.chatGPTProfile(s.cfg.CheckoutProfile)
	client, err := s.newGptClient(ctx, cred, profile)
	if err != nil {
		return tierProbe{Source: "wham_usage", ErrorMessage: err.Error()}
	}
	defer client.close()
	client.setAuthorization("Bearer " + accessToken)
	client.setHeader("Accept", "application/json")
	if accountID := accessTokenAccountID(accessToken); accountID != "" {
		client.setHeader("ChatGPT-Account-Id", accountID)
	}
	resp, err := client.request(ctx, http.MethodGet, "https://chatgpt.com/backend-api/wham/usage", requestOptions{})
	if err != nil {
		return tierProbe{Source: "wham_usage", ErrorMessage: "wham usage probe failed: " + err.Error()}
	}
	if resp.status != http.StatusOK {
		return tierProbe{Source: "wham_usage", ErrorMessage: fmt.Sprintf("wham usage returned status %d", resp.status)}
	}
	if result := tierResult(stringx.FirstNonEmpty(stringAt(resp.json, "plan_type"), stringAt(resp.json, "planType")), "wham_usage.plan_type"); result.Checked {
		return result
	}
	return tierProbe{Checked: true, Source: "wham_usage", ErrorMessage: "wham usage returned no plan_type"}
}

func (s *Server) probePlusActiveSessionToken(ctx context.Context, sessionToken string, cred credential) tierProbe {
	if strings.TrimSpace(sessionToken) == "" {
		return tierProbe{Source: "auth_session", ErrorMessage: "session_token is required"}
	}
	profile := cred.chatGPTProfile(s.cfg.CheckoutProfile)
	client, err := s.newGptClient(ctx, cred, profile)
	if err != nil {
		return tierProbe{Source: "auth_session", ErrorMessage: err.Error()}
	}
	defer client.close()
	client.setHeader("Accept", "application/json")
	resp, err := client.request(ctx, http.MethodGet, "https://chatgpt.com/api/auth/session", requestOptions{})
	if err != nil {
		return tierProbe{Source: "auth_session", ErrorMessage: "auth session probe failed: " + err.Error()}
	}
	if resp.status != http.StatusOK {
		return tierProbe{Source: "auth_session", ErrorMessage: fmt.Sprintf("auth session returned status %d", resp.status)}
	}
	accessToken := stringAt(resp.json, "accessToken")
	if accessToken != "" {
		wham := s.probeTierAccessToken(ctx, accessToken, cred)
		if wham.Checked && stringx.FirstNonEmpty(wham.Tier, wham.PlanType) != "" {
			return wham
		}
		if result := accessTokenTier(accessToken, "accessToken.auth"); result.Checked {
			return result
		}
	}
	result := detectPlusActiveFromSessionPayload(resp.json)
	if !result.Checked {
		result.ErrorMessage = "auth session returned no authenticated user"
	}
	return result
}
