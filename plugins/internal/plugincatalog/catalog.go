package plugincatalog

import (
	"strings"

	"github.com/byte-v-forge/gpt/pkg/gptplugin"
)

func N8NAction(owner string, actionID string, displayName string, workflowKey string, workflowIDPrefix string, startPath string, actionScope string, webhookPath string, actionPathPrefix string, requestProto string, responseProto string, buttonLabel string, placement string, staleSteps []string) gptplugin.ActionDefinition {
	return gptplugin.ActionDefinition{
		ActionID:      actionID,
		DisplayName:   displayName,
		Owner:         owner,
		Engine:        gptplugin.EngineN8N,
		RequestProto:  requestProto,
		ResponseProto: responseProto,
		Visibility:    "account",
		Workflow: gptplugin.WorkflowDefinition{
			Key:              workflowKey,
			IDPrefix:         workflowIDPrefix,
			StartPath:        startPath,
			N8NActionScope:   actionScope,
			N8NWebhookPath:   webhookPath,
			ActionPathPrefix: actionPathPrefix,
		},
		UIButtons:              []gptplugin.UIButton{{ID: actionButtonID(actionID), Label: buttonLabel, Placement: placement}},
		StaleSteps:             append([]string(nil), staleSteps...),
		BlockedAccountStatuses: defaultBlockedStatuses(),
	}
}

func WithUIButton(def gptplugin.ActionDefinition, label string, placement string) gptplugin.ActionDefinition {
	def.UIButtons = append(def.UIButtons, gptplugin.UIButton{
		ID:        actionButtonID(def.ActionID) + "-" + actionButtonID(placement),
		Label:     label,
		Placement: placement,
	})
	return def
}

func WithRequiredStatuses(def gptplugin.ActionDefinition, statuses ...string) gptplugin.ActionDefinition {
	def.RequiredAccountStatuses = append(def.RequiredAccountStatuses, statuses...)
	return def
}

func WithBlockedStatuses(def gptplugin.ActionDefinition, statuses ...string) gptplugin.ActionDefinition {
	def.BlockedAccountStatuses = append(def.BlockedAccountStatuses, statuses...)
	return def
}

func WithRequiredFields(def gptplugin.ActionDefinition, fields ...string) gptplugin.ActionDefinition {
	def.RequiredFields = append(def.RequiredFields, fields...)
	return def
}

func WithCapabilities(def gptplugin.ActionDefinition, capabilities ...string) gptplugin.ActionDefinition {
	def.Capabilities = append(def.Capabilities, capabilities...)
	return def
}

func Field(key string, label string, kind gptplugin.ConfigFieldKind, defaultValue string) gptplugin.ConfigField {
	return gptplugin.ConfigField{Key: key, Label: label, Kind: kind, DefaultValue: defaultValue}
}

func actionButtonID(actionID string) string {
	return strings.ToLower(strings.ReplaceAll(actionID, "_", "-"))
}

func defaultBlockedStatuses() []string {
	return []string{"DEACTIVATED", "EMAIL_ALREADY_EXISTS", "USER_ALREADY_EXISTS"}
}
