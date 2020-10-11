package application

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/pkuebler/bahn-bot/pkg/infrastructure/marudor"
	"github.com/pkuebler/bahn-bot/pkg/trainalarms/domain"
)

// NotifyUsers check train delay threshold and call the notifyFn
func (a *Application) NotifyUsers(ctx context.Context, notifyFn func(ctx context.Context, alarm *domain.TrainAlarm, train marudor.HafasTrain, diff time.Duration) error) error {
	// sort alarms by lastNotificationAt
	alarms, err := a.repo.GetTrainAlarmsSortByLastNotificationAt(ctx, 20)
	if err != nil {
		a.log.Error(err)
		return errors.New("internal server error")
	}

	for _, alarm := range alarms {
		go a.notifyUser(ctx, notifyFn, *alarm)
	}

	return nil
}

func (a *Application) notifyUser(ctx context.Context, notifyFn func(ctx context.Context, alarm *domain.TrainAlarm, train marudor.HafasTrain, diff time.Duration) error, alarm domain.TrainAlarm) error {
	log := a.log.WithFields(logrus.Fields{
		"alarmID":    alarm.GetID(),
		"identifyer": alarm.GetIdentifyer(),
		"plattform":  alarm.GetPlattform(),
		"trainname":  alarm.GetTrainName(),
	})

	// UpdateTrainAlarm use a transaction
	err := a.repo.UpdateTrainAlarm(ctx, alarm.GetID(), func(t *domain.TrainAlarm) (*domain.TrainAlarm, error) {
		// search train
		// todo: use cache
		train, err := a.hafas.GetTrainByStation(ctx, t.GetTrainName(), t.GetStationEVA(), t.GetStationDate())
		if err != nil {
			log.Error(err)
			return nil, errors.New("not found")
		}

		if train.CurrentStop == nil {
			log.Trace("current stop empty")
			t.SetSuccessfulNotification()
			return t, nil
		}

		var scheduledTime time.Time
		var currentTime time.Time

		if train.CurrentStop.Arrival != nil {
			scheduledTime = train.CurrentStop.Arrival.GoScheduledTime()
			currentTime = train.CurrentStop.Arrival.GoTime()
		} else if train.CurrentStop.Departure != nil {
			scheduledTime = train.CurrentStop.Departure.GoScheduledTime()
			currentTime = train.CurrentStop.Departure.GoTime()
		} else {
			log.Trace("current stop arrival and departure empty")
			t.SetSuccessfulNotification()
			return t, nil
		}

		diff := currentTime.Sub(scheduledTime)

		if int(diff.Minutes()) == t.GetLastDelay() {
			log.Trace("same delay")
			t.SetSuccessfulNotification()
			return t, nil
		}

		minutes := int(diff.Minutes())

		if alarm.GetDelayThresholdMinutes() > minutes {
			if alarm.GetLastDelay() < alarm.GetDelayThresholdMinutes() {
				t.SetLastDelay(minutes)
				t.SetSuccessfulNotification()
				return t, nil
			}
		}

		t.SetLastDelay(minutes)

		if err = notifyFn(ctx, t, *train, diff); err != nil {
			log.Error(err)
			return nil, errors.New("can't notify user")
		}

		t.SetSuccessfulNotification()
		return t, nil
	})

	if err != nil {
		log.Error(err)
		return err
	}

	log.Trace("train alarm updated")
	return nil
}
