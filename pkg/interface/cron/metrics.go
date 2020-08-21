package cron

import "github.com/prometheus/client_golang/prometheus"

// CronMetrics Registry
type CronMetrics struct {
	TrainAlertNotificationsTotal *prometheus.CounterVec
}

// NewCronMetrics return a new metric registry
func NewCronMetrics() *CronMetrics {
	cronPrefix := "bahn_bot_cron_"

	total := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: cronPrefix + "notifications_total",
		Help: "How many train alert notifications sended, partitoned by train name",
	}, []string{"trainname"})

	register := &CronMetrics{
		TrainAlertNotificationsTotal: total,
	}

	prometheus.MustRegister(register.TrainAlertNotificationsTotal)

	return register
}
