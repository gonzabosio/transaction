package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gonzabosio/transaction/gateway/utils"
	"github.com/gonzabosio/transaction/model"
	"github.com/gonzabosio/transaction/services/proto/inventory"
	"github.com/gonzabosio/transaction/services/proto/order"
)

func (gw *Gateway) OrderGateway(w http.ResponseWriter, r *http.Request) {
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

	// inventory service
	res, err := gw.svs.Inventory.GetStock(ctx, &inventory.ProductRequest{ProductId: reqBody.ProductId})
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
			"message":    "Client cannot make an order of the product",
			"error_info": "Product is not available",
		}, http.StatusUnprocessableEntity)
		return
	}

	// order service
	accessToken, err := gw.svs.Order.NewAccessToken(ctx, &order.Client{ClientAuth: gw.clientAuth})
	if err != nil {
		utils.WriteJSON(w, map[string]string{
			"message":    "Failed to generate access token",
			"error_info": err.Error(),
		}, http.StatusBadRequest)
		return
	}
	price := strconv.FormatFloat(float64(res.Price), 'f', -1, 32)
	newOrder, err := gw.svs.Order.NewOrder(ctx, &order.Order{AccessToken: accessToken.Value, Currency: "USD", Amount: price})
	if err != nil {
		utils.WriteJSON(w, map[string]string{
			"message":    "Failed to create new order",
			"error_info": err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	orderDetails, err := gw.svs.Order.GetOrderDetails(ctx, &order.OrderDetailsRequest{Id: newOrder.OrderId, AccessToken: accessToken.Value})
	if err != nil {
		utils.WriteJSON(w, map[string]string{
			"message":    "Failed to get order details",
			"error_info": err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	if err := gw.cache.SaveAccessToken(newOrder.OrderId, accessToken.Value); err != nil {
		utils.WriteJSON(w, map[string]string{
			"message":    "Failed to save access token in cache memory",
			"error_info": err.Error(),
		}, http.StatusBadRequest)
		return
	}

	utils.WriteJSON(w, map[string]interface{}{
		"message":       fmt.Sprintf("Order for %s was created", res.Name),
		"stock":         res.Stock,
		"order_details": orderDetails,
	}, http.StatusCreated)
}
