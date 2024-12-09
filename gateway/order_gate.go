package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gonzabosio/transaction/gateway/utils"
	"github.com/gonzabosio/transaction/model"
	"github.com/gonzabosio/transaction/services/proto/inventory"
)

func (gw *gateway) OrderGateway(w http.ResponseWriter, r *http.Request) {
	reqBody := new(model.ProductRequest)
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		utils.WriteJSON(w, map[string]string{
			"message":    "Failed to read request",
			"error_info": err.Error(),
		}, http.StatusBadRequest)
		return
	}
	if err := utils.ValidateStruct(reqBody); err != nil {
		utils.WriteJSON(w, map[string]string{
			"message":    "Failed to validate request body",
			"error_info": err.Error(),
		}, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	res, err := gw.svs.Inventory.GetStock(ctx, &inventory.ProductRequest{Name: reqBody.Name})
	if err != nil {
		if err == context.DeadlineExceeded {
			utils.WriteJSON(w, map[string]string{
				"message":    "Getting stock process timed out",
				"error_info": err.Error(),
			}, http.StatusInternalServerError)
			return
		} else if err == context.Canceled {
			utils.WriteJSON(w, map[string]string{
				"message":    "Process to get stock was cancelled",
				"error_info": err.Error(),
			}, http.StatusInternalServerError)
			return
		} else {
			utils.WriteJSON(w, map[string]string{
				"message":    "Failed to get product stock",
				"error_info": err.Error(),
			}, http.StatusNotFound)
			return
		}
	}
	if !res.IsAvailable {
		utils.WriteJSON(w, map[string]string{
			"message": "Product is not available",
		}, http.StatusUnprocessableEntity)
		return
	}

	// 2. create order

	utils.WriteJSON(w, map[string]interface{}{
		"message": fmt.Sprintf("Order for %s was created", reqBody.Name),
		"stock":   res.Stock,
	}, http.StatusCreated)
}
