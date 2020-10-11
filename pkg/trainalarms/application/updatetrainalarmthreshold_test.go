package application

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/pkuebler/bahn-bot/pkg/trainalarms/domain"
)

func TestUpdateTrainAlarmThreshold(t *testing.T) {
	app, repo := createTestCase(true)

	ctx := context.Background()
	alarm, err := domain.NewTrainAlarm(
		uuid.New().String(),
		"telegram",
		"ice 4",
		8503000,
		int64(1595797980000),
		time.Now(),
		"Berlin Ostbahnhof",
	)
	assert.Nil(t, err)
	assert.NotNil(t, alarm)

	repo.TrainAlarms[alarm.GetID()] = alarm

	// Alarm found
	cmd := UpdateTrainAlarmThresholdCmd{
		AlarmID:          alarm.GetID(),
		ThresholdMinutes: 10,
	}

	err = app.UpdateTrainAlarmThreshold(ctx, cmd)
	assert.Nil(t, err)
	assert.Equal(t, repo.TrainAlarms[alarm.GetID()].GetDelayThresholdMinutes(), 10)

	// Alarm not found
	cmd = UpdateTrainAlarmThresholdCmd{
		AlarmID:          "123",
		ThresholdMinutes: 10,
	}

	err = app.UpdateTrainAlarmThreshold(ctx, cmd)
	assert.NotNil(t, err)
}
