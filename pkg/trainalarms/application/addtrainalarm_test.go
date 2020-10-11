package application

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAddTrainAlarm(t *testing.T) {
	app, repo := createTestCase(true)
	ctx := context.Background()

	cmd := AddTrainAlarmCmd{
		Identifyer:  uuid.New().String(),
		Plattform:   "telegram",
		TrainName:   "ice 4",
		StationEVA:  8503000,
		StationDate: int64(1595797980000),
	}

	err := app.AddTrainAlarm(ctx, cmd)
	assert.Nil(t, err)
	assert.Len(t, repo.TrainAlarms, 1)
}

func TestAddTrainAlarmNotFound(t *testing.T) {
	app, repo := createTestCase(false)
	ctx := context.Background()

	cmd := AddTrainAlarmCmd{
		Identifyer:  uuid.New().String(),
		Plattform:   "telegram",
		TrainName:   "ice 4",
		StationEVA:  8503000,
		StationDate: int64(1595797980000),
	}

	err := app.AddTrainAlarm(ctx, cmd)
	assert.NotNil(t, err)
	assert.Len(t, repo.TrainAlarms, 0)
}
