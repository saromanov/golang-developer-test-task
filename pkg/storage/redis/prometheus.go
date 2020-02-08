package redis

import (
	"github.com/prometheus/client_golang/prometheus"
)

var totalReads = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "total_reads",
	})

var totalWrites = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "failed_writes",
	})

func initPrometheus() {
	prometheus.MustRegister(totalReads)
	prometheus.MustRegister(totalWrites)
}
