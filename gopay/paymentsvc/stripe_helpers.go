package paymentsvc

import (
	"strings"

	"github.com/google/uuid"
)

func containsString(value any, needle string) bool {
	items, ok := value.([]any)
	if !ok {
		return false
	}
	for _, item := range items {
		if stringAt(map[string]any{"value": item}, "value") == needle {
			return true
		}
	}
	return false
}

func uuidHex() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "")
}

func newElementsSessionID() string {
	value := uuidHex()
	if len(value) > 11 {
		value = value[:11]
	}
	return "elements_session_" + value
}
