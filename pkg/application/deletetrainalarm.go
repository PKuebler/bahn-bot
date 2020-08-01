package application

import (
	"context"
	"errors"

	"github.com/pkuebler/bahn-bot/pkg/domain/trainalarm"
	"github.com/sirupsen/logrus"
)

// DeleteTrainAlarmCmd for DeleteTrainAlarm
type DeleteTrainAlarmCmd struct {
	AlarmID string
}

// DeleteTrainAlarm at database
func (a *Application) DeleteTrainAlarm(ctx context.Context, cmd DeleteTrainAlarmCmd) (*trainalarm.TrainAlarm, error) {
	log := a.log.WithFields(logrus.Fields{
		"alarmID": cmd.AlarmID,
	})

	// search alarm at repository
	alarm, err := a.repo.GetTrainAlarm(ctx, cmd.AlarmID)
	if err != nil {
		log.Error(err)
		return nil, errors.New("internal server error")
	}
	if alarm == nil {
		log.Trace("not found")
		return nil, errors.New("not found")
	}

	// delete alarm at repository
	if err := a.repo.DeleteTrainAlarm(ctx, alarm.GetID()); err != nil {
		log.Error(err)
		return nil, errors.New("internal server error")
	}

	log.Trace("alarm deleted")
	return alarm, nil
}
