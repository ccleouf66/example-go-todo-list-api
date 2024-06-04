package api

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metricsHandler struct {
	opsProcessed prometheus.Counter
}

func NewLetricsHandler() metricsHandler {
	return metricsHandler{
		opsProcessed: promauto.NewCounter(prometheus.CounterOpts{
			Name: "todoapi_query_total",
			Help: "The total number of query",
		}),
	}
}

func (h *metricsHandler) GetMetrics() http.HandlerFunc {
	return promhttp.Handler().ServeHTTP
}

func (h *metricsHandler) IncrementTotalQueryMetric(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.opsProcessed.Inc()
		next.ServeHTTP(w, r)
	})
}
