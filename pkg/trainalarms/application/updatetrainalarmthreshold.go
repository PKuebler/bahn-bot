package application

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"

	"github.com/pkuebler/bahn-bot/pkg/trainalarms/domain"
)

// UpdateTrainAlarmThresholdCmd for UpdateTrainAlarmThreshold
type UpdateTrainAlarmThresholdCmd struct {
	AlarmID          string
	ThresholdMinutes int
}

// UpdateTrainAlarmThreshold at database
func (a *Application) UpdateTrainAlarmThreshold(ctx context.Context, cmd UpdateTrainAlarmThresholdCmd) error {
	log := a.log.WithFields(logrus.Fields{
		"alarmID":          cmd.AlarmID,
		"thresholdMinutes": cmd.ThresholdMinutes,
	})

	// search alarm at repository
	err := a.repo.UpdateTrainAlarm(ctx, cmd.AlarmID, func(alarm *domain.TrainAlarm) (*domain.TrainAlarm, error) {
		alarm.SetDelayThresholdMinutes(cmd.ThresholdMinutes)
		return alarm, nil
	})
	if err != nil {
		log.Trace("not found")
		return errors.New("not found")
	}

	log.Trace("alarm updated")
	return nil
}
