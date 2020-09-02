package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// MetricNamePrefix of all metric names
	MetricNamePrefix = "bahn_bot_"

	// telegram
	metricTelegramPrefix         = MetricNamePrefix + "telegram_"
	metricTelegramInPrefix       = metricTelegramPrefix + "in_"
	metricTelegramOutPrefix      = metricTelegramPrefix + "out_"
	telegramInUpdatesTotalName   = metricTelegramInPrefix + "updates_total"
	telegramInMessagesTotalName  = metricTelegramInPrefix + "messages_total"
	telegramInCommandsTotalName  = metricTelegramInPrefix + "commands_total"
	telegramInQueriesTotalName   = metricTelegramInPrefix + "queries_total"
	telegramOutMessagesTotalName = metricTelegramOutPrefix + "messages_total"
	telegramRequestDurationName  = metricTelegramPrefix + "duration_seconds"
)

// PrometheusRegistry to collect all metrics
type PrometheusRegistry struct {
	TelegramInUpdatesTotal   prometheus.Counter
	TelegramInMessagesTotal  prometheus.Counter
	TelegramInCommandsTotal  prometheus.Counter
	TelegramInQueriesTotal   prometheus.Counter
	TelegramOutMessagesTotal prometheus.Counter
	TelegramRequestDuration  prometheus.Summary
}

// NewPrometheusMetric configuration
func NewPrometheusMetric() *PrometheusRegistry {
	telegramInUpdatesTotal := newCounter(telegramInUpdatesTotalName, "help...")
	telegramInMessagesTotal := newCounter(telegramInMessagesTotalName, "help...")
	telegramInCommandsTotal := newCounter(telegramInCommandsTotalName, "help...")
	telegramInQueriesTotal := newCounter(telegramInQueriesTotalName, "help...")
	telegramOutMessagesTotal := newCounter(telegramOutMessagesTotalName, "help...")
	telegramRequestDuration := newSummary(telegramRequestDurationName, "help...")

	metricsRegistry := &PrometheusRegistry{
		TelegramInUpdatesTotal:   telegramInUpdatesTotal,
		TelegramInMessagesTotal:  telegramInMessagesTotal,
		TelegramInCommandsTotal:  telegramInCommandsTotal,
		TelegramInQueriesTotal:   telegramInQueriesTotal,
		TelegramOutMessagesTotal: telegramOutMessagesTotal,
		TelegramRequestDuration:  telegramRequestDuration,
	}

	prometheus.MustRegister(metricsRegistry.TelegramInUpdatesTotal)
	prometheus.MustRegister(metricsRegistry.TelegramInMessagesTotal)
	prometheus.MustRegister(metricsRegistry.TelegramInCommandsTotal)
	prometheus.MustRegister(metricsRegistry.TelegramInQueriesTotal)
	prometheus.MustRegister(metricsRegistry.TelegramOutMessagesTotal)
	prometheus.MustRegister(metricsRegistry.TelegramRequestDuration)

	return metricsRegistry
}

func newSummary(name string, help string) prometheus.Summary {
	return prometheus.NewSummary(prometheus.SummaryOpts{
		Name: name,
		Help: help,
	})
}

func newCounter(name string, help string) prometheus.Counter {
	return prometheus.NewCounter(prometheus.CounterOpts{
		Name: name,
		Help: help,
	})
}

func newHistogram(name string, help string) prometheus.Histogram {
	return prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: name,
		Help: help,
	})
}
