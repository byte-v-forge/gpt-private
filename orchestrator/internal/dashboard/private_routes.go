//go:build private_plugins

package dashboard

import (
	"fmt"
	"net/http"
	"strings"

	"orchestrator/internal/contracts"
)

func (s *server) privateRouteBindings() []routeBinding {
	return nil
}

func (s *server) goPayActionsEnabled() bool {
	if s == nil || s.actionRegistry == nil {
		return false
	}
	return s.actionRegistry.HasCapability(contracts.CapabilityGoPay)
}

func (s *server) handlePrivateJobAction(w http.ResponseWriter, r *http.Request, jobID string, parts []string) bool {
	if len(parts) < 2 {
		return false
	}
	switch parts[1] {
	case "gopay-payment":
		if len(parts) != 3 || parts[2] != "confirm" {
			writeError(w, http.StatusNotFound, fmt.Errorf("unsupported job gopay-payment action: %s", strings.Join(parts[1:], "/")))
			return true
		}
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return true
		}
		s.confirmManualGoPayPayment(w, r, jobID)
		return true
	default:
		return false
	}
}
