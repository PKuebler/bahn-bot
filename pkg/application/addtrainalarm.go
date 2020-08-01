package application

import (
	"context"
	"errors"
	"time"

	"github.com/pkuebler/bahn-bot/pkg/domain/trainalarm"
	"github.com/sirupsen/logrus"
)

// AddTrainAlarmCmd for AddTrainAlarm
type AddTrainAlarmCmd struct {
	Identifyer  string
	Plattform   string
	TrainName   string
	StationEVA  int
	StationDate int64
}

// AddTrainAlarm to database
func (a *Application) AddTrainAlarm(ctx context.Context, cmd AddTrainAlarmCmd) error {
	log := a.log.WithFields(logrus.Fields{
		"identifyer": cmd.Identifyer,
		"plattform":  cmd.Plattform,
		"trainname":  cmd.TrainName,
	})

	// search train by marudor
	_, err := a.hafas.GetTrainByStation(ctx, cmd.TrainName, cmd.StationEVA, cmd.StationDate)
	if err != nil {
		log.Error(err)
		return errors.New("not found")
	}

	// create train alarm with final arrival
	finalArrival := time.Now().AddDate(0, 0, -3) // todo: use train
	alarm, err := trainalarm.NewTrainAlarm(
		cmd.Identifyer,
		cmd.Plattform,
		cmd.TrainName,
		cmd.StationEVA,
		cmd.StationDate,
		finalArrival,
	)

	_, err = a.repo.GetOrCreateTrainAlarm(ctx, alarm)
	if err != nil {
		log.Error(err)
		return errors.New("internal server error")
	}

	log.Trace("alarm created")
	return nil
}
