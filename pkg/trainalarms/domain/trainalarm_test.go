package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewTrainAlarm(t *testing.T) {
	alarm, err := NewTrainAlarm("identifyer", "telegram", "ice 6", 2342342, int64(234234234234), time.Now(), "Berlin Ostbahnhof")
	assert.Nil(t, err)
	assert.NotNil(t, alarm)
}

func TestTrainAlarmGetters(t *testing.T) {
	identifyer := uuid.New().String()
	plattform := "telegram"
	trainName := "ice 6"
	stationEVA := 2342342
	stationDate := int64(234234234234)
	finalArrivalAt := time.Now()
	finalDestinationName := "Berlin Ostbahnhof"
	alarm, err := NewTrainAlarm(identifyer, plattform, trainName, stationEVA, stationDate, finalArrivalAt, finalDestinationName)
	assert.Nil(t, err)
	assert.NotNil(t, alarm)

	assert.Equal(t, identifyer, alarm.GetIdentifyer())
	assert.Equal(t, plattform, alarm.GetPlattform())
	assert.Equal(t, trainName, alarm.GetTrainName())
	assert.Equal(t, stationEVA, alarm.GetStationEVA())
	assert.Equal(t, stationDate, alarm.GetStationDate())
	assert.Equal(t, finalArrivalAt, alarm.GetFinalArrivalAt())
	assert.Equal(t, finalDestinationName, alarm.GetFinalDestinationName())
	assert.Nil(t, alarm.GetLastNotificationAt())
	alarm.SetSuccessfulNotification()
	assert.NotNil(t, alarm.GetLastNotificationAt())

	assert.Equal(t, 0, alarm.GetDelayThresholdMinutes())
	alarm.SetDelayThresholdMinutes(-20)
	assert.Equal(t, 0, alarm.GetDelayThresholdMinutes())
	alarm.SetDelayThresholdMinutes(20)
	assert.Equal(t, 20, alarm.GetDelayThresholdMinutes())
}
