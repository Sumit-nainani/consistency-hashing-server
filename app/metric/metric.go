package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	UniqueIPsGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "http_unique_ips_last_1m",
		Help: "Number of unique client IPs in the last 1 minute",
	})
)

func Init() {
	prometheus.MustRegister(UniqueIPsGauge)
}
