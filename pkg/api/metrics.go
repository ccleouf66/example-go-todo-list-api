package api

import (
	"math/rand"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metricsHandler struct {
	opsProcessed  prometheus.Counter
	desiredPodNum prometheus.Gauge
}

func NewLetricsHandler() metricsHandler {
	return metricsHandler{
		opsProcessed: promauto.NewCounter(prometheus.CounterOpts{
			Name: "todoapi_query_total",
			Help: "The total number of query",
		}),
		desiredPodNum: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "todoapi_desired_pod",
			Help: "The desired number of api pod",
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

func (h *metricsHandler) RandDesiredPodNumber() {
	r := 1 + rand.Float64()*(20-1)
	h.desiredPodNum.Set(r)
}
