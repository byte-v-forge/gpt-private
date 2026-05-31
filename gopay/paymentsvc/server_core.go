package paymentsvc

import "github.com/byte-v-forge/gpt-private/gopay/pb"

type Server struct {
	pb.UnimplementedPaymentServiceServer
	cfg   Config
	flows *flowStore
}

func NewServer(cfg Config) *Server {
	return &Server{cfg: cfg, flows: &flowStore{items: map[string]*pendingFlow{}}}
}

type pendingFlow struct {
	charger         *charger
	state           map[string]any
	useAccountToken bool
}

func (f *pendingFlow) close() {
	if f != nil && f.charger != nil {
		f.charger.close()
	}
}
