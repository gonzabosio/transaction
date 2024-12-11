package gateway

import (
	"net/http"

	"github.com/gonzabosio/transaction/gateway/utils"
)

func (gw *Gateway) PaymentGateway(w http.ResponseWriter, r *http.Request) {
	// gw.mq.RunCheckoutTasks()
	utils.WriteJSON(w, map[string]interface{}{
		"message": "Payment is being processed ⏳",
	}, http.StatusAccepted)
}
