package server

import (
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


func initPrometheus() {
	prometheus.MustRegister(totalRequests)
	prometheus.MustRegister(failedRequests)
}
