package repository

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"

	"github.com/pkuebler/bahn-bot/pkg/domain/trainalarm"
)

func TestNewSQLDatabase(t *testing.T) {
	db := NewSQLDatabase(os.Getenv("DB_DIALECT"), os.Getenv("DB_PATH"))
	assert.NotNil(t, db)

	db.Close()
}

func TestConvertModel(t *testing.T) {
	alarm, err := trainalarm.NewTrainAlarm("mychatID", "telegram", "ICE 7", 123456, 1234567890, time.Now())
	assert.NotNil(t, alarm)
	assert.Nil(t, err)

	alarm.SetDelayThresholdMinutes(10)
	alarm.SetSuccessfulNotification()
	alarm.SetLastDelay(5)

	// convert
	model := NewSQLTrainAlarmModel(alarm)
	assert.NotNil(t, model)

	assert.Equal(t, alarm.GetID(), model.ID)
	assert.Equal(t, alarm.GetIdentifyer(), model.Identifyer)
	assert.Equal(t, alarm.GetPlattform(), model.Plattform)
	assert.Equal(t, alarm.GetTrainName(), model.TrainName)
	assert.Equal(t, alarm.GetStationEVA(), model.StationEVA)
	assert.Equal(t, alarm.GetStationDate(), model.StationDate)
	assert.Equal(t, alarm.GetFinalArrivalAt(), model.FinalArrivalAt)
	assert.Equal(t, alarm.GetDelayThresholdMinutes(), model.DelayThresholdMinutes)
	assert.Equal(t, alarm.GetLastNotificationAt(), model.LastNotificationAt)
	assert.Equal(t, alarm.GetLastDelay(), model.LastDelayMinutes)

	// convert back
	back := model.TrainAlarm()
	assert.NotNil(t, back)

	assert.Equal(t, model.ID, back.GetID())
	assert.Equal(t, model.Identifyer, back.GetIdentifyer())
	assert.Equal(t, model.Plattform, back.GetPlattform())
	assert.Equal(t, model.TrainName, back.GetTrainName())
	assert.Equal(t, model.StationEVA, back.GetStationEVA())
	assert.Equal(t, model.StationDate, back.GetStationDate())
	assert.Equal(t, model.FinalArrivalAt, back.GetFinalArrivalAt())
	assert.Equal(t, model.DelayThresholdMinutes, back.GetDelayThresholdMinutes())
	assert.Equal(t, model.LastNotificationAt, back.GetLastNotificationAt())
	assert.Equal(t, model.LastDelayMinutes, back.GetLastDelay())
}

func TestSQLGetOrCreateTrainAlarm(t *testing.T) {
	ctx := context.Background()
	db := NewSQLDatabase(os.Getenv("DB_DIALECT"), os.Getenv("DB_PATH"))
	assert.NotNil(t, db)
	defer db.Close()
	err := db.db.Delete(SQLTrainAlarmModel{}).Error
	assert.Nil(t, err)

	query := createTestTrainAlarm(time.Now())

	// no entry found -> create
	alarm, err := db.GetOrCreateTrainAlarm(ctx, query)
	assert.Nil(t, err)
	assert.NotNil(t, alarm)
	var dbAlarm SQLTrainAlarmModel
	err = db.db.Where("id = ?", alarm.GetID()).First(&dbAlarm).Error
	assert.Nil(t, err)
	assert.Equal(t, alarm.GetID(), dbAlarm.ID)

	// entry found -> return old entry
	query = createTestTrainAlarm(time.Now())
	oldAlarm, err := db.GetOrCreateTrainAlarm(ctx, query)
	assert.Nil(t, err)
	assert.NotNil(t, alarm)
	assert.Equal(t, alarm.GetID(), oldAlarm.GetID())
}

func TestSQLGetTrainAlarm(t *testing.T) {
	ctx := context.Background()
	db := NewSQLDatabase(os.Getenv("DB_DIALECT"), os.Getenv("DB_PATH"))
	assert.NotNil(t, db)
	defer db.Close()
	err := db.db.Delete(SQLTrainAlarmModel{}).Error
	assert.Nil(t, err)

	query := createTestTrainAlarm(time.Now())

	// not found
	alarm, err := db.GetTrainAlarm(ctx, query.GetID())
	assert.Nil(t, err)
	assert.Nil(t, alarm)

	// founded
	err = db.db.Create(NewSQLTrainAlarmModel(query)).Error
	assert.Nil(t, err)

	alarm, err = db.GetTrainAlarm(ctx, query.GetID())
	assert.Nil(t, err)
	assert.NotNil(t, alarm)
}

func TestSQLGetTrainAlarms(t *testing.T) {
	ctx := context.Background()
	db := NewSQLDatabase(os.Getenv("DB_DIALECT"), os.Getenv("DB_PATH"))
	assert.NotNil(t, db)
	defer db.Close()
	err := db.db.Delete(SQLTrainAlarmModel{}).Error
	assert.Nil(t, err)

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

		err = db.db.Create(NewSQLTrainAlarmModel(alarm)).Error
		assert.Nil(t, err)
	}

	alarms, err = db.GetTrainAlarms(ctx, "1234", "telegram")
	assert.Nil(t, err)
	assert.Len(t, alarms, 4)
}

func TestSQLGetTrainAlarmsSortByLastNotificationAt(t *testing.T) {
	ctx := context.Background()
	db := NewSQLDatabase(os.Getenv("DB_DIALECT"), os.Getenv("DB_PATH"))
	assert.NotNil(t, db)
	defer db.Close()
	err := db.db.Delete(SQLTrainAlarmModel{}).Error
	assert.Nil(t, err)

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

		err = db.db.Create(NewSQLTrainAlarmModel(alarm)).Error
		assert.Nil(t, err)
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

func TestSQLDeleteTrainAlarm(t *testing.T) {
	ctx := context.Background()
	db := NewSQLDatabase(os.Getenv("DB_DIALECT"), os.Getenv("DB_PATH"))
	assert.NotNil(t, db)
	defer db.Close()
	err := db.db.Delete(SQLTrainAlarmModel{}).Error
	assert.Nil(t, err)

	alarm := createTestTrainAlarm(time.Now())

	// not found
	err = db.DeleteTrainAlarm(ctx, alarm.GetID())
	assert.NotNil(t, err)

	// delete
	err = db.db.Create(NewSQLTrainAlarmModel(alarm)).Error
	assert.Nil(t, err)

	err = db.DeleteTrainAlarm(ctx, alarm.GetID())
	assert.Nil(t, err)

	var dbAlarm SQLTrainAlarmModel
	err = db.db.Where("id = ?", alarm.GetID()).First(&dbAlarm).Error
	assert.NotNil(t, err)
}

func TestSQLUpdateTrainAlarm(t *testing.T) {
	ctx := context.Background()
	db := NewSQLDatabase(os.Getenv("DB_DIALECT"), os.Getenv("DB_PATH"))
	assert.NotNil(t, db)
	defer db.Close()
	err := db.db.Delete(SQLTrainAlarmModel{}).Error
	assert.Nil(t, err)

	alarm := createTestTrainAlarm(time.Now())

	// not found
	notRunning := true
	err = db.UpdateTrainAlarm(ctx, alarm.GetID(), func(alarm *trainalarm.TrainAlarm) (*trainalarm.TrainAlarm, error) {
		notRunning = false
		return nil, nil
	})
	assert.True(t, notRunning)
	assert.NotNil(t, err)

	// update with error
	err = db.db.Create(NewSQLTrainAlarmModel(alarm)).Error
	assert.Nil(t, err)

	running := false
	err = db.UpdateTrainAlarm(ctx, alarm.GetID(), func(alarm *trainalarm.TrainAlarm) (*trainalarm.TrainAlarm, error) {
		running = true
		return nil, errors.New("error")
	})
	assert.True(t, running)
	assert.NotNil(t, err)

	var dbAlarm SQLTrainAlarmModel
	err = db.db.Where("id = ?", alarm.GetID()).First(&dbAlarm).Error
	assert.Nil(t, err)
	assert.Equal(t, alarm.GetID(), dbAlarm.ID)

	// update
	err = db.db.Create(NewSQLTrainAlarmModel(alarm)).Error
	assert.Nil(t, err)

	running = false
	err = db.UpdateTrainAlarm(ctx, alarm.GetID(), func(alarm *trainalarm.TrainAlarm) (*trainalarm.TrainAlarm, error) {
		running = true
		alarm.SetSuccessfulNotification()
		return alarm, nil
	})
	assert.True(t, running)
	assert.Nil(t, err)

	err = db.db.Where("id = ?", alarm.GetID()).First(&dbAlarm).Error
	assert.Nil(t, err)
	assert.Equal(t, alarm.GetID(), dbAlarm.ID)
	assert.NotNil(t, alarm.GetLastNotificationAt())
}
