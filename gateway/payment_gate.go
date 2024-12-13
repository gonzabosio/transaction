package gateway

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gonzabosio/transaction/gateway/utils"
)

type PaymentRequest struct {
	OrderId string `json:"order_id" validate:"required"`
}

func (gw *Gateway) PaymentGateway(w http.ResponseWriter, r *http.Request) {
	payload := &PaymentRequest{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteJSON(w, map[string]string{
			"message":    "Failed to read order payload",
			"error_info": err.Error(),
		}, http.StatusBadRequest)
		return
	}
	if err := utils.ValidateStruct(payload); err != nil {
		utils.WriteJSON(w, map[string]string{
			"message":    "Invalid payment payload",
			"error_info": err.Error(),
		}, http.StatusBadRequest)
		return
	}
	h := r.Header.Get("Authorization")
	accessToken := strings.TrimPrefix(h, "Bearer ")
	if err := gw.mq.RunPaymentCheckoutTasks(payload.OrderId, accessToken); err != nil {
		utils.WriteJSON(w, map[string]string{
			"message":    "Payment could not be processed",
			"error_info": err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, map[string]interface{}{
		"message": "Payment is being processed ⏳",
	}, http.StatusAccepted)
}
