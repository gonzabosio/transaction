package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gonzabosio/transaction/gateway"
)

func NewRouter(gw *gateway.Gateway) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/order", gw.OrderGateway)
	r.Post("/payment", gw.PaymentGateway)
	return r
}
