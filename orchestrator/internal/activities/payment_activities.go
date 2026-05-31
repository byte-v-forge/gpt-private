//go:build private_plugins

package activities

import (
	"context"
)

func (s *Server) GoPayPaymentPrepareActivity(ctx context.Context, input GoPayActivityInput) (GoPayPaymentPrepareOutput, error) {
	output := GoPayPaymentPrepareOutput{}
	stopHeartbeat := startActivityHeartbeat(ctx, input.GetJobId(), stepGoPayPaymentPrepare, "preparing gopay payment", protoData(goPayPaymentHeartbeatFields(input)))
	defer stopHeartbeat()

	step, err := s.startActivityStep(ctx, input.GetJobId(), stepGoPayPaymentPrepare, false, true)
	if err != nil {
		return output, err
	}
	account, err := s.paymentActivityAccount(ctx, &input)
	if err != nil {
		return output, step.complete(protoData(map[string]any{"error_message": err.Error()}), err)
	}

	output, err = s.prepareGoPayPayment(ctx, step, input, account)
	if err != nil {
		return output, step.complete(output.GetData(), err)
	}
	return output, step.complete(output.GetData(), nil)
}

func (s *Server) GoPayPaymentPrepareCheckoutActivity(ctx context.Context, input GoPayActivityInput) (GoPayPaymentPrepareOutput, error) {
	output := GoPayPaymentPrepareOutput{}
	stopHeartbeat := startActivityHeartbeat(ctx, input.GetJobId(), stepGoPayPaymentPrepareCheckout, "creating gopay payment checkout", protoData(goPayPaymentHeartbeatFields(input)))
	defer stopHeartbeat()

	step, err := s.startActivityStep(ctx, input.GetJobId(), stepGoPayPaymentPrepareCheckout, false, true)
	if err != nil {
		return output, err
	}
	account, err := s.paymentActivityAccount(ctx, &input)
	if err != nil {
		return output, step.complete(protoData(map[string]any{"error_message": err.Error()}), err)
	}

	output, err = s.prepareGoPayPaymentCheckout(ctx, step, input, account)
	if err != nil {
		return output, step.complete(output.GetData(), err)
	}
	return output, step.complete(output.GetData(), nil)
}

func (s *Server) GoPayPaymentPrepareRefreshActivity(ctx context.Context, input GoPayActivityInput) (GoPayPaymentPrepareOutput, error) {
	output := GoPayPaymentPrepareOutput{}
	stopHeartbeat := startActivityHeartbeat(ctx, input.GetJobId(), stepGoPayPaymentPrepareRefresh, "refreshing gopay payment checkout", protoData(goPayPaymentHeartbeatFields(input)))
	defer stopHeartbeat()

	step, err := s.startActivityStep(ctx, input.GetJobId(), stepGoPayPaymentPrepareRefresh, false, true)
	if err != nil {
		return output, err
	}
	output, err = s.refreshGoPayPaymentCheckout(ctx, step, input)
	if err != nil {
		return output, step.complete(output.GetData(), err)
	}
	return output, step.complete(output.GetData(), nil)
}

func (s *Server) GoPayPaymentPrepareLinkActivity(ctx context.Context, input GoPayActivityInput) (GoPayPaymentPrepareOutput, error) {
	output := GoPayPaymentPrepareOutput{}
	stopHeartbeat := startActivityHeartbeat(ctx, input.GetJobId(), stepGoPayPaymentPrepareLink, "linking gopay payment checkout", protoData(goPayPaymentHeartbeatFields(input)))
	defer stopHeartbeat()

	step, err := s.startActivityStep(ctx, input.GetJobId(), stepGoPayPaymentPrepareLink, false, true)
	if err != nil {
		return output, err
	}
	output, err = s.prepareGoPayPaymentLink(ctx, step, input)
	if err != nil {
		return output, step.complete(output.GetData(), err)
	}
	return output, step.complete(output.GetData(), nil)
}
