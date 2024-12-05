package gateway

import (
	"net/http"
)

type gateway struct{}

func NewAPIGateway() *gateway {
	return &gateway{}
}

func (gw *gateway) PaymentGateway(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Trigger payment process"))
}
