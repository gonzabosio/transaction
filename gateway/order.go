package gateway

import (
	"net/http"
)

func (gw *gateway) OrderGateway(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Trigger order creation"))
}
