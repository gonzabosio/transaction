package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gonzabosio/transaction/gateway"
	middleprom "github.com/gonzabosio/transaction/prometheus/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(gw *gateway.Gateway, m *middleprom.Metrics) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(m.TrackMetrics)

	// expose metrics
	r.Handle("/metrics", promhttp.Handler())

	r.Post("/order", gw.OrderGateway)
	r.Post("/payment", gw.PaymentGateway)
	return r
}
