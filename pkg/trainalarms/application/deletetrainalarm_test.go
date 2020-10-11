package application

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/pkuebler/bahn-bot/pkg/trainalarms/domain"
)

func TestDeleteTrainAlarm(t *testing.T) {
	app, repo := createTestCase(true)
	ctx := context.Background()

	alarm, err := domain.NewTrainAlarm(uuid.New().String(), "telegram", "ice 4", 8503000, int64(1595797980000), time.Now(), "Berlin Ostbahnhof")
	assert.Nil(t, err)
	assert.NotNil(t, alarm)

	cmd := DeleteTrainAlarmCmd{
		AlarmID: alarm.GetID(),
	}

	// entry not found
	a, err := app.DeleteTrainAlarm(ctx, cmd)
	assert.NotNil(t, err)
	assert.Nil(t, a)

	// add entry
	repo.TrainAlarms[alarm.GetID()] = alarm

	// entry deleted
	a, err = app.DeleteTrainAlarm(ctx, cmd)
	assert.Nil(t, err)
	assert.NotNil(t, a)
	assert.Len(t, repo.TrainAlarms, 0)
}
