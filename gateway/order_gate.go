package gateway

import (
	"net/http"

	"github.com/gonzabosio/transaction/gateway/utils"
)

func (gw *gateway) OrderGateway(w http.ResponseWriter, r *http.Request) {
	// 1. query inventory
	// 2. create order
	utils.WriteJSON(w, map[string]interface{}{
		"message": "Order was created successfully",
	}, http.StatusCreated)
}
