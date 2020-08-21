package marudor

import "github.com/prometheus/client_golang/prometheus"

// APIMetrics Registry
type APIMetrics struct {
	RequestsTotal            *prometheus.CounterVec
	RequestDurationSeconds   *prometheus.HistogramVec
	RequestsByTrainNameTotal *prometheus.CounterVec
}

// NewAPIMetrics return a new metric registry
func NewAPIMetrics() *APIMetrics {
	marudorPrefix := "marudor_"

	buckets := []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}

	total := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: marudorPrefix + "requests_total",
		Help: "How many HTTP requests processed, partitoned by status code and endpoint",
	}, []string{"status_code", "endpoint"})

	duration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    marudorPrefix + "duration_seconds",
		Help:    "request duration",
		Buckets: buckets,
	}, []string{"status_code", "endpoint"})

	trainNames := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: marudorPrefix + "requests_by_trainname_total",
		Help: "How many requests processed, partitoned by status code and train name",
	}, []string{"status_code", "trainname"})

	register := &APIMetrics{
		RequestsTotal:            total,
		RequestDurationSeconds:   duration,
		RequestsByTrainNameTotal: trainNames,
	}

	prometheus.MustRegister(register.RequestsTotal)
	prometheus.MustRegister(register.RequestDurationSeconds)
	prometheus.MustRegister(register.RequestsByTrainNameTotal)

	return register
}
