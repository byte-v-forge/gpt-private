//go:build private_plugins

package api

import "strings"

func goPayAppAccountID(value string) string {
	return strings.TrimSpace(value)
}
