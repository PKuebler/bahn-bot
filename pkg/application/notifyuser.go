package application

import (
	"context"
	"errors"
	"time"

	"github.com/pkuebler/bahn-bot/pkg/domain/trainalarm"
	"github.com/pkuebler/bahn-bot/pkg/infrastructure/marudor"
	"github.com/sirupsen/logrus"
)

// NotifyUsers check train delay threshold and call the notifyFn
func (a *Application) NotifyUsers(ctx context.Context, notifyFn func(ctx context.Context, alarm *trainalarm.TrainAlarm, train marudor.HafasTrain, diff time.Duration) error) error {
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

func (a *Application) notifyUser(ctx context.Context, notifyFn func(ctx context.Context, alarm *trainalarm.TrainAlarm, train marudor.HafasTrain, diff time.Duration) error, alarm trainalarm.TrainAlarm) error {
	log := a.log.WithFields(logrus.Fields{
		"alarmID":    alarm.GetID(),
		"identifyer": alarm.GetIdentifyer(),
		"plattform":  alarm.GetPlattform(),
		"trainname":  alarm.GetTrainName(),
	})

	// UpdateTrainAlarm use a transaction
	err := a.repo.UpdateTrainAlarm(ctx, alarm.GetID(), func(t *trainalarm.TrainAlarm) (*trainalarm.TrainAlarm, error) {
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
			scheduledTime = time.Unix(0, train.CurrentStop.Arrival.ScheduledTime*int64(time.Millisecond))
			currentTime = time.Unix(0, train.CurrentStop.Arrival.Time*int64(time.Millisecond))
		} else if train.CurrentStop.Departure != nil {
			scheduledTime = time.Unix(0, train.CurrentStop.Departure.ScheduledTime*int64(time.Millisecond))
			currentTime = time.Unix(0, train.CurrentStop.Departure.Time*int64(time.Millisecond))
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
