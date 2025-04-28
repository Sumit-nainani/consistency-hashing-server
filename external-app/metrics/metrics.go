package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pod_request_total",
			Help: "Total number of requests handled by this pod",
		},
		[]string{"pod_ip"},
	)
)

func Init() {
	prometheus.MustRegister(RequestCount)
}
