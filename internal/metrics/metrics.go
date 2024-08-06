package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	buckets = []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}
)

type MetricsCollector struct {
	RequestDurationSecondsHistogram *prometheus.HistogramVec
}

func New(prefix string, ps *prometheus.Registry) MetricsCollector {
	mc := MetricsCollector{
		RequestDurationSecondsHistogram: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: prefix,
			Name:      "request_duration_seconds",
			Help:      "Histogram of response time for handler in seconds",
			Buckets:   buckets,
		}, []string{"route", "method", "status_code"}),
	}

	ps.MustRegister(
		mc.RequestDurationSecondsHistogram,
	)

	return mc
}
