package cron

import (
	"context"
	"fmt"
	"time"

	"github.com/pkuebler/bahn-bot/pkg/application"
	"github.com/pkuebler/bahn-bot/pkg/domain/trainalarm"
	"github.com/pkuebler/bahn-bot/pkg/infrastructure/marudor"
	"github.com/pkuebler/bahn-bot/pkg/infrastructure/telegramconversation"
	"github.com/sirupsen/logrus"
)

// CronJob triggers applications
type CronJob struct {
	application   *application.Application
	notifications chan telegramconversation.TContext
	log           *logrus.Entry
	metrics       *CronMetrics
}

// NewCronJob service
func NewCronJob(log *logrus.Entry, application *application.Application, metrics bool) *CronJob {
	c := &CronJob{
		application:   application,
		notifications: make(chan telegramconversation.TContext),
		log:           log,
	}

	if metrics {
		c.metrics = NewCronMetrics()
	}

	return c
}

// Start ticker
func (c *CronJob) Start(ctx context.Context) {
	clearDatabaseTicker := time.NewTicker(1 * time.Hour)
	notifyTicker := time.NewTicker(1 * time.Minute)

	c.log.Info("start cronjob")

	c.ClearDatabase(ctx)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-clearDatabaseTicker.C:
				c.ClearDatabase(ctx)
			case <-notifyTicker.C:
				c.log.Info("notify users")
				c.NotifyUsers(ctx)
			}
		}
	}()
}

// ClearDatabase delete old trainalarms from database
func (c *CronJob) ClearDatabase(ctx context.Context) {
	c.log.Info("Clear database...")
	c.application.DeleteOldTrainAlarms(ctx)
	c.application.DeleteOldStates(ctx)
}

// NotifyUsers about train delays
func (c *CronJob) NotifyUsers(ctx context.Context) {
	c.application.NotifyUsers(ctx, func(ctx context.Context, alarm *trainalarm.TrainAlarm, train marudor.HafasTrain, diff time.Duration) error {
		tctx := telegramconversation.NewTContext(alarm.GetIdentifyer())

		if c.metrics != nil {
			c.metrics.TrainAlertNotificationsTotal.WithLabelValues(alarm.GetTrainName()).Inc()
		}

		txt := fmt.Sprintf("Zug %s hat `%s` Verspätung.", alarm.GetTrainName(), diff.String())

		tctx.Send(txt)

		c.notifications <- tctx

		return nil
	})
}

// NotificationChannel returns the channel with telegram notifications
func (c *CronJob) NotificationChannel() chan telegramconversation.TContext {
	return c.notifications
}
