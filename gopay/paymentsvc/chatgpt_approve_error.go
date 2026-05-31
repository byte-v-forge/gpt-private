package paymentsvc

import (
	"errors"
	"fmt"
	"strings"
)

type chatGPTApproveBlockedError struct {
	status int
	body   string
}

func (e chatGPTApproveBlockedError) Error() string {
	return fmt.Sprintf("chatgpt approve: result=\"blocked\" body=%s", e.body)
}

func isChatGPTApproveBlocked(err error) bool {
	if err == nil {
		return false
	}
	var blocked chatGPTApproveBlockedError
	if errors.As(err, &blocked) {
		return true
	}
	return strings.Contains(strings.ToLower(err.Error()), "chatgpt approve") && strings.Contains(strings.ToLower(err.Error()), "blocked")
}
