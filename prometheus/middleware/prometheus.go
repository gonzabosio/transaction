package middleware

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	RequestCount   *prometheus.CounterVec
	ErrorCount     *prometheus.CounterVec
	RequestLatency *prometheus.HistogramVec
}

func MetricsSetup() *Metrics {
	return &Metrics{
		RequestCount: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "transaction_requests_total",
				Help: "Count of request processed by the transaction server",
			},
			[]string{"path", "method", "status"},
		),
		// sum by(path, method, status)(transaction_requests_total{path!="/metrics", path!="/app-worker.js"})
		ErrorCount: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "transaction_requests_errors_total",
				Help: "Count of error requests processed by the transaction server",
			},
			[]string{"path", "method", "status"},
		),
		// sum by(path, method, status)(transaction_requests_errors_total{path!="/metrics", path!="/app-worker.js"})
		RequestLatency: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "transaction_request_latency_seconds",
				Help: "Latency of requests processed by the transaction server",
			},
			[]string{"path", "method", "status"},
		),
		// rate(transaction_request_latency_seconds_sum{path!="/metrics", path!="/app-worker.js"}[5m]) / ignoring(instance, job) rate(transaction_request_latency_seconds_count{path!="metrics", path!="/app-worker.js"}[5m])
	}
}

func (m *Metrics) PrometheusInit() {
	prometheus.MustRegister(
		m.RequestCount,
		m.ErrorCount,
		m.RequestLatency,
	)
}

type ResponseWriterTracker struct {
	http.ResponseWriter
	StatusCode int
}

func (cwr *ResponseWriterTracker) WriteHeader(code int) {
	cwr.StatusCode = code
	cwr.ResponseWriter.WriteHeader(cwr.StatusCode)
}

func (m *Metrics) TrackMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		crw := &ResponseWriterTracker{w, http.StatusOK}

		next.ServeHTTP(crw, r)
		// total requests - errors
		path := r.URL.Path
		status := crw.StatusCode
		m.RequestCount.WithLabelValues(path, r.Method, http.StatusText(status)).Inc()
		if status >= 400 {
			m.ErrorCount.WithLabelValues(path, r.Method, http.StatusText(status)).Inc()
		}
		// latency
		m.RequestLatency.WithLabelValues(path, r.Method, http.StatusText(status)).
			Observe(time.Since(start).Seconds())
	})
}
