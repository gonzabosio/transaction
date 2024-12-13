package gateway

import (
	"encoding/json"
	"net/http"

	"github.com/gonzabosio/transaction/gateway/utils"
)

type PaymentRequest struct {
	OrderId   string `json:"order_id" validate:"required"`
	ProductId int64  `json:"product_id" validate:"required"`
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

	accessToken, err := gw.cache.GetAccessToken(payload.OrderId)
	if err != nil {
		utils.WriteJSON(w, map[string]string{
			"message":    "Failed to get access token from cache",
			"error_info": err.Error(),
		}, http.StatusBadRequest)
		return
	}

	if err := gw.mq.RunPaymentCheckoutTasks(payload.OrderId, accessToken, payload.ProductId); err != nil {
		utils.WriteJSON(w, map[string]string{
			"message":    "Payment could not be processed",
			"error_info": err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, map[string]string{
		"message": "Payment is being processed ⏳",
	}, http.StatusAccepted)
}
