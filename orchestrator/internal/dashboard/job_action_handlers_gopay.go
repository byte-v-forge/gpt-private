//go:build private_plugins

package dashboard

import (
	"context"
	"errors"
	"net/http"

	"orchestrator/pb"
)

type goPayManualPaymentConfirmer interface {
	ConfirmManualGoPayPayment(ctx context.Context, req *pb.ConfirmManualGoPayPaymentRequest) (*pb.ConfirmManualGoPayPaymentResponse, error)
}

func (s *server) confirmManualGoPayPayment(w http.ResponseWriter, r *http.Request, jobID string) {
	confirmer, ok := s.n8nActions.(goPayManualPaymentConfirmer)
	if !ok {
		writeError(w, http.StatusServiceUnavailable, errors.New("gopay manual payment confirmation is not configured"))
		return
	}
	resp, err := confirmer.ConfirmManualGoPayPayment(r.Context(), &pb.ConfirmManualGoPayPaymentRequest{
		JobId: jobID,
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}
	if resp.GetErrorMessage() != "" {
		writeError(w, http.StatusBadRequest, errors.New(resp.GetErrorMessage()))
		return
	}
	writeJSON(w, http.StatusOK, resp)
}
