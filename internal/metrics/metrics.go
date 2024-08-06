package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MetricsCollector struct {
	OpsProcessed prometheus.Counter
}

func New(prefix string, ps *prometheus.Registry) MetricsCollector {
	return MetricsCollector{
		OpsProcessed: promauto.NewCounter(prometheus.CounterOpts{
			Name: prefix + "_processed_ops_total",
			Help: "The total number of processed events",
		}),
	}
}
