package trainalarm

import (
	"time"

	"github.com/google/uuid"
)

// TrainAlarm represent a alarm from a user
type TrainAlarm struct {
	id                    string
	identifyer            string
	plattform             string
	trainName             string
	stationEVA            int
	stationDate           int64
	finalArrivalAt        time.Time
	finalDestinationName  string
	delayThresholdMinutes int
	lastNotificationAt    *time.Time
	lastDelayMinutes      int
}

// NewTrainAlarm returns a new train alarm
func NewTrainAlarm(
	identifyer string,
	plattform string,
	trainName string,
	stationEVA int,
	stationDate int64,
	finalArrivalAt time.Time,
	finalDestinationName string,
) (*TrainAlarm, error) {
	return &TrainAlarm{
		id:                    uuid.New().String(),
		identifyer:            identifyer,
		plattform:             plattform,
		trainName:             trainName,
		stationEVA:            stationEVA,
		stationDate:           stationDate,
		finalArrivalAt:        finalArrivalAt,
		finalDestinationName:  finalDestinationName,
		delayThresholdMinutes: 0,
		lastDelayMinutes:      0,
	}, nil
}

// NewTrainAlarmFromRepository returns a new train alarm
func NewTrainAlarmFromRepository(
	id string,
	identifyer string,
	plattform string,
	trainName string,
	stationEVA int,
	stationDate int64,
	finalArrivalAt time.Time,
	finalDestinationName string,
	delayThresholdMinutes int,
	lastNotificationAt *time.Time,
	lastDelayMinutes int,
) (*TrainAlarm, error) {
	return &TrainAlarm{
		id:                    id,
		identifyer:            identifyer,
		plattform:             plattform,
		trainName:             trainName,
		stationEVA:            stationEVA,
		stationDate:           stationDate,
		finalArrivalAt:        finalArrivalAt,
		finalDestinationName:  finalDestinationName,
		delayThresholdMinutes: delayThresholdMinutes,
		lastNotificationAt:    lastNotificationAt,
		lastDelayMinutes:      lastDelayMinutes,
	}, nil
}

// SetSuccessfulNotification add the current timestamp
func (t *TrainAlarm) SetSuccessfulNotification() {
	now := time.Now()
	t.lastNotificationAt = &now
}

// GetID from domain
func (t *TrainAlarm) GetID() string {
	return t.id
}

// GetIdentifyer from user
func (t *TrainAlarm) GetIdentifyer() string {
	return t.identifyer
}

// GetPlattform from user
func (t *TrainAlarm) GetPlattform() string {
	return t.plattform
}

// GetTrainName from train
func (t *TrainAlarm) GetTrainName() string {
	return t.trainName
}

// GetStationEVA from train
func (t *TrainAlarm) GetStationEVA() int {
	return t.stationEVA
}

// GetStationDate from train
func (t *TrainAlarm) GetStationDate() int64 {
	return t.stationDate
}

// GetFinalArrivalAt from train
func (t *TrainAlarm) GetFinalArrivalAt() time.Time {
	return t.finalArrivalAt
}

// GetFinalDestinationName from train
func (t *TrainAlarm) GetFinalDestinationName() string {
	return t.finalDestinationName
}

// GetLastNotificationAt returns nil if no notification is send yet
func (t *TrainAlarm) GetLastNotificationAt() *time.Time {
	return t.lastNotificationAt
}

// GetLastDelay returns last crawled delay in minutes
func (t *TrainAlarm) GetLastDelay() int {
	return t.lastDelayMinutes
}

// SetLastDelay in minutes
func (t *TrainAlarm) SetLastDelay(delayMinutes int) {
	t.lastDelayMinutes = delayMinutes
}

// SetDelayThresholdMinutes for train
func (t *TrainAlarm) SetDelayThresholdMinutes(minutes int) {
	if minutes < 0 {
		minutes = 0
	}

	t.delayThresholdMinutes = minutes
}

// GetDelayThresholdMinutes before alarm
func (t *TrainAlarm) GetDelayThresholdMinutes() int {
	return t.delayThresholdMinutes
}

// Compare two train alarms without id
func (t *TrainAlarm) Compare(other *TrainAlarm) bool {
	return t.GetIdentifyer() == other.GetIdentifyer() &&
		t.GetPlattform() == other.GetPlattform() &&
		t.GetTrainName() == other.GetTrainName() &&
		t.GetStationEVA() == other.GetStationEVA() &&
		t.GetStationDate() == other.GetStationDate()
}
