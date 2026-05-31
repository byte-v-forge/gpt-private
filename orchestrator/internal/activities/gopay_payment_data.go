//go:build private_plugins

package activities

import (
	"strings"

	"google.golang.org/protobuf/types/known/structpb"
)

func protoData(data map[string]any) *structpb.Struct {
	if data == nil {
		return &structpb.Struct{}
	}
	out, err := structpb.NewStruct(data)
	if err != nil {
		return &structpb.Struct{}
	}
	return out
}

func normalizeGoPayWorkflowStateJSON(stateJSON string) string {
	stateJSON = strings.TrimSpace(stateJSON)
	if stateJSON == "" {
		return "{}"
	}
	return stateJSON
}

func normalizeIndonesiaPhone(phone string) string {
	phone = strings.TrimSpace(phone)
	phone = strings.TrimPrefix(phone, "+")
	if strings.HasPrefix(phone, "0") {
		return "62" + strings.TrimPrefix(phone, "0")
	}
	return phone
}
