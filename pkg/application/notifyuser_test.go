package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/pkuebler/bahn-bot/pkg/domain/trainalarm"
	"github.com/pkuebler/bahn-bot/pkg/infrastructure/marudor"
)

func TestNotifyUser(t *testing.T) {
	type TestCase struct {
		Description   string
		HafasFound    bool
		DatabaseFound bool
		NotifySuccess bool
	}

	testCases := []TestCase{
		{
			Description:   "not at database",
			HafasFound:    true,
			DatabaseFound: false,
			NotifySuccess: true,
		},
		{
			Description:   "not at hafas",
			HafasFound:    false,
			DatabaseFound: true,
			NotifySuccess: true,
		},
		{
			Description:   "notify falied",
			HafasFound:    true,
			DatabaseFound: true,
			NotifySuccess: false,
		},
		{
			Description:   "founded",
			HafasFound:    true,
			DatabaseFound: true,
			NotifySuccess: true,
		},
	}

	ctx := context.Background()
	for _, testCase := range testCases {
		app, repo := createTestCase(testCase.HafasFound)

		alarm, err := trainalarm.NewTrainAlarm(
			uuid.New().String(),
			"telegram",
			"ice 4",
			8503000,
			int64(1595797980000),
			time.Now(),
		)
		alarm.SetLastDelay(-100)
		assert.Nil(t, err, testCase.Description)
		assert.NotNil(t, alarm, testCase.Description)

		if testCase.DatabaseFound {
			repo.TrainAlarms[alarm.GetID()] = alarm
		}

		isRunning := false
		err = app.notifyUser(ctx, func(ctx context.Context, alarm *trainalarm.TrainAlarm, train marudor.HafasTrain, diff time.Duration) error {
			isRunning = true

			if !testCase.NotifySuccess {
				isRunning = false
				return errors.New("notification failed")
			}

			return nil
		}, *alarm)

		if !testCase.HafasFound || !testCase.DatabaseFound || !testCase.NotifySuccess {
			assert.NotNil(t, err, testCase.Description)
			assert.False(t, isRunning, testCase.Description)
		} else {
			assert.Nil(t, err, testCase.Description)
			assert.True(t, isRunning, testCase.Description)
		}
	}

}

func TestNotifyUsers(t *testing.T) {
	app, repo := createTestCase(true)
	ctx := context.Background()

	for i := 0; i < 4; i++ {
		alarm, err := trainalarm.NewTrainAlarm(
			uuid.New().String(),
			"telegram",
			"ice 4",
			8503000,
			int64(1595797980000),
			time.Now(),
		)
		assert.Nil(t, err)
		assert.NotNil(t, alarm)
		repo.TrainAlarms[alarm.GetID()] = alarm
	}

	err := app.NotifyUsers(ctx, func(ctx context.Context, alarm *trainalarm.TrainAlarm, train marudor.HafasTrain, diff time.Duration) error {
		return nil
	})
	assert.Nil(t, err)
}
