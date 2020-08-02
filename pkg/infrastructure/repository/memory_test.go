package repository

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/pkuebler/bahn-bot/pkg/domain/trainalarm"
)

func TestNewMemoryDatabase(t *testing.T) {
	db := NewMemoryDatabase()
	assert.NotNil(t, db)
	assert.NotNil(t, db.Mutex)

	db.Mutex.Lock()
	db.Mutex.Unlock()
	db.TrainAlarms["test"] = createTestTrainAlarm(time.Now())
	assert.NotNil(t, db.TrainAlarms["test"])
}

func TestGetOrCreateTrainAlarm(t *testing.T) {
	db := NewMemoryDatabase()
	ctx := context.Background()

	query := createTestTrainAlarm(time.Now())

	// no entry found -> create
	alarm, err := db.GetOrCreateTrainAlarm(ctx, query)
	assert.Nil(t, err)
	assert.NotNil(t, alarm)
	_, ok := db.TrainAlarms[alarm.GetID()]
	assert.True(t, ok)
	assert.Equal(t, alarm.GetID(), db.TrainAlarms[alarm.GetID()].GetID())

	// entry found -> return old entry
	query = createTestTrainAlarm(time.Now())
	oldAlarm, err := db.GetOrCreateTrainAlarm(ctx, query)
	assert.Nil(t, err)
	assert.NotNil(t, alarm)
	assert.Equal(t, alarm.GetID(), oldAlarm.GetID())
}

func TestGetTrainAlarm(t *testing.T) {
	db := NewMemoryDatabase()
	ctx := context.Background()

	query := createTestTrainAlarm(time.Now())

	// not found
	alarm, err := db.GetTrainAlarm(ctx, query.GetID())
	assert.Nil(t, err)
	assert.Nil(t, alarm)

	// founded
	db.TrainAlarms[query.GetID()] = query
	alarm, err = db.GetTrainAlarm(ctx, query.GetID())
	assert.Nil(t, err)
	assert.NotNil(t, alarm)
}

func TestGetTrainAlarms(t *testing.T) {
	db := NewMemoryDatabase()
	ctx := context.Background()

	// empty database
	alarms, err := db.GetTrainAlarms(ctx, "1234", "telegram")
	assert.Nil(t, err)
	assert.Len(t, alarms, 0)

	for i := 0; i < 4; i++ {
		alarm, err := trainalarm.NewTrainAlarm("1234", "telegram", fmt.Sprintf("ice %d", i), i, int64(i), time.Now())
		assert.Nil(t, err)
		if i > 1 {
			alarm.SetSuccessfulNotification()
		}
		db.TrainAlarms[alarm.GetID()] = alarm
	}

	alarms, err = db.GetTrainAlarms(ctx, "1234", "telegram")
	assert.Nil(t, err)
	assert.Len(t, alarms, 4)
}

func TestGetTrainAlarmsSortByLastNotificationAt(t *testing.T) {
	db := NewMemoryDatabase()
	ctx := context.Background()

	// empty database
	alarms, err := db.GetTrainAlarmsSortByLastNotificationAt(ctx, 20)
	assert.Nil(t, err)
	assert.Len(t, alarms, 0)

	for i := 0; i < 4; i++ {
		alarm, err := trainalarm.NewTrainAlarm("identifyer", "telegram", fmt.Sprintf("ice %d", i), i, int64(i), time.Now())
		assert.Nil(t, err)
		if i > 1 {
			alarm.SetSuccessfulNotification()
		}
		db.TrainAlarms[alarm.GetID()] = alarm
	}

	alarms, err = db.GetTrainAlarmsSortByLastNotificationAt(ctx, 20)
	assert.Nil(t, err)
	assert.Len(t, alarms, 4)

	assert.Equal(t, 2+3, alarms[0].GetStationEVA()+alarms[1].GetStationEVA())
	assert.Equal(t, 0+1, alarms[2].GetStationEVA()+alarms[3].GetStationEVA())

	alarms, err = db.GetTrainAlarmsSortByLastNotificationAt(ctx, 2)
	assert.Nil(t, err)
	assert.Len(t, alarms, 2)
}

func TestDeleteTrainAlarm(t *testing.T) {
	db := NewMemoryDatabase()
	ctx := context.Background()

	alarm := createTestTrainAlarm(time.Now())

	// not found
	err := db.DeleteTrainAlarm(ctx, alarm.GetID())
	assert.NotNil(t, err)

	// delete
	db.TrainAlarms[alarm.GetID()] = alarm
	err = db.DeleteTrainAlarm(ctx, alarm.GetID())
	assert.Nil(t, err)
	_, ok := db.TrainAlarms[alarm.GetID()]
	assert.False(t, ok)
}

func TestUpdateTrainAlarm(t *testing.T) {
	db := NewMemoryDatabase()
	ctx := context.Background()

	alarm := createTestTrainAlarm(time.Now())

	// not found
	notRunning := true
	err := db.UpdateTrainAlarm(ctx, alarm.GetID(), func(alarm *trainalarm.TrainAlarm) (*trainalarm.TrainAlarm, error) {
		notRunning = false
		return nil, nil
	})
	assert.True(t, notRunning)
	assert.NotNil(t, err)

	// update with error
	db.TrainAlarms[alarm.GetID()] = alarm
	running := false
	err = db.UpdateTrainAlarm(ctx, alarm.GetID(), func(alarm *trainalarm.TrainAlarm) (*trainalarm.TrainAlarm, error) {
		running = true
		return nil, errors.New("error")
	})
	assert.True(t, running)
	assert.NotNil(t, err)
	assert.NotNil(t, db.TrainAlarms[alarm.GetID()])

	// update
	db.TrainAlarms[alarm.GetID()] = alarm
	running = false
	err = db.UpdateTrainAlarm(ctx, alarm.GetID(), func(alarm *trainalarm.TrainAlarm) (*trainalarm.TrainAlarm, error) {
		running = true
		alarm.SetSuccessfulNotification()
		return alarm, nil
	})
	assert.True(t, running)
	assert.Nil(t, err)
	assert.NotNil(t, db.TrainAlarms[alarm.GetID()])
	assert.NotNil(t, db.TrainAlarms[alarm.GetID()].GetLastNotificationAt())
}

func TestDeleteOldStates(t *testing.T) {
	db := NewMemoryDatabase()
	ctx := context.Background()

	db.CurrentState = append(db.CurrentState, &MemoryStateModel{
		Identifyer:    "true",
		Plattform:     "telegram",
		State:         "start",
		StatePlayload: "",
		UpdatedAt:     time.Now().AddDate(0, -1, 0),
	})

	db.CurrentState = append(db.CurrentState, &MemoryStateModel{
		Identifyer:    "false",
		Plattform:     "telegram",
		State:         "start",
		StatePlayload: "",
		UpdatedAt:     time.Now(),
	})

	err := db.DeleteOldStates(ctx, time.Now().AddDate(0, 0, -2))
	assert.Nil(t, err)
	assert.Len(t, db.CurrentState, 0)
}

func createTestTrainAlarm(finalArrival time.Time) *trainalarm.TrainAlarm {
	alarm, _ := trainalarm.NewTrainAlarm(
		"identifyer",
		"telegram",
		"ice 6",
		11111111,
		int64(12432342342334234),
		finalArrival,
	)
	return alarm
}
