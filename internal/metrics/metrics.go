package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MetricsCollector struct {
	ProcessedOpsTotalCounter prometheus.Counter
}

func New(prefix string, ps *prometheus.Registry) MetricsCollector {
	mc := MetricsCollector{
		ProcessedOpsTotalCounter: promauto.NewCounter(prometheus.CounterOpts{
			Name: prefix + "_processed_ops_total",
			Help: "The total number of processed events",
		}),
	}

	ps.MustRegister(mc.ProcessedOpsTotalCounter)

	return mc
}
