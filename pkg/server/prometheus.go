package server

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
)

var totalRequests = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "total_requests",
	})

var failedRequests = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "failed_requests",
	})

var statusCodes = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Processing of status codes.",
	},
	[]string{"code", "method"},
)

func writeStatusCode(code int, method string) {
	statusCodes.WithLabelValues(fmt.Sprintf("%d", code), method).Inc()
}
func initPrometheus() {
	prometheus.MustRegister(totalRequests)
	prometheus.MustRegister(failedRequests)
	prometheus.MustRegister(statusCodes)
}
