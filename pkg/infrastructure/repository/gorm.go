package repository

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/pkuebler/bahn-bot/pkg/domain/trainalarm"
)

// SQLDatabase to persistence
type SQLDatabase struct {
	db *gorm.DB
}

// NewSQLDatabase by gorm
func NewSQLDatabase(dialect string, path string) *SQLDatabase {
	db, err := gorm.Open(dialect, path)
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&SQLTrainAlarmModel{}, &SQLStateModel{})

	return &SQLDatabase{
		db: db,
	}
}

// SQLTrainAlarmModel to save train alarm
type SQLTrainAlarmModel struct {
	ID                    string `gorm:"primary_key"`
	Identifyer            string `gorm:"index:identifyer_plattform"`
	Plattform             string `gorm:"index:identifyer_plattform"`
	TrainName             string
	StationEVA            int
	StationDate           int64
	FinalArrivalAt        time.Time
	DelayThresholdMinutes int
	LastNotificationAt    *time.Time
	LastDelayMinutes      int
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

// NewSQLTrainAlarmModel convert domain model to repository model
func NewSQLTrainAlarmModel(alarm *trainalarm.TrainAlarm) *SQLTrainAlarmModel {
	return &SQLTrainAlarmModel{
		ID:                    alarm.GetID(),
		Identifyer:            alarm.GetIdentifyer(),
		Plattform:             alarm.GetPlattform(),
		TrainName:             alarm.GetTrainName(),
		StationEVA:            alarm.GetStationEVA(),
		StationDate:           alarm.GetStationDate(),
		FinalArrivalAt:        alarm.GetFinalArrivalAt(),
		DelayThresholdMinutes: alarm.GetDelayThresholdMinutes(),
		LastNotificationAt:    alarm.GetLastNotificationAt(),
		LastDelayMinutes:      alarm.GetLastDelay(),
	}
}

// TrainAlarm convert to TrainAlarm domain model
func (s *SQLTrainAlarmModel) TrainAlarm() *trainalarm.TrainAlarm {
	alarm, _ := trainalarm.NewTrainAlarmFromRepository(
		s.ID,
		s.Identifyer,
		s.Plattform,
		s.TrainName,
		s.StationEVA,
		s.StationDate,
		s.FinalArrivalAt,
		s.DelayThresholdMinutes,
		s.LastNotificationAt,
		s.LastDelayMinutes,
	)

	return alarm
}

// SQLStateModel to save conversation state
type SQLStateModel struct {
	Identifyer   string `gorm:"primary_key"`
	Plattform    string `gorm:"primary_key"`
	State        string
	StatePayload string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// GetOrCreateTrainAlarm creates a new alarm if none exists yet
func (s *SQLDatabase) GetOrCreateTrainAlarm(ctx context.Context, alarm *trainalarm.TrainAlarm) (*trainalarm.TrainAlarm, error) {
	var oldModel SQLTrainAlarmModel
	if res := s.db.Where(
		"identifyer = ? AND plattform = ? AND train_name = ? AND station_eva = ? AND station_date = ?",
		alarm.GetIdentifyer(),
		alarm.GetPlattform(),
		alarm.GetTrainName(),
		alarm.GetStationEVA(),
		alarm.GetStationDate(),
	).Take(&oldModel); res.Error != nil {
		if !res.RecordNotFound() {
			return nil, res.Error
		}

		model := NewSQLTrainAlarmModel(alarm)
		if err := s.db.Create(model).Error; err != nil {
			return nil, err
		}
		return model.TrainAlarm(), nil
	}

	return oldModel.TrainAlarm(), nil
}

// GetTrainAlarm by id. returns NO error if nothing found
func (s *SQLDatabase) GetTrainAlarm(ctx context.Context, alarmID string) (*trainalarm.TrainAlarm, error) {
	var alarm SQLTrainAlarmModel
	if res := s.db.Where("id = ?", alarmID).Take(&alarm); res.Error != nil {
		if res.RecordNotFound() {
			return nil, nil
		}
		return nil, res.Error
	}

	return alarm.TrainAlarm(), nil
}

// GetTrainAlarms by identifyer and plattform. returns NO error if nothing found
func (s *SQLDatabase) GetTrainAlarms(ctx context.Context, identifyer string, plattform string) ([]*trainalarm.TrainAlarm, error) {
	var results []SQLTrainAlarmModel
	if res := s.db.Where("identifyer = ? AND plattform = ?", identifyer, plattform).Find(&results); res.Error != nil {
		if res.RecordNotFound() {
			return nil, nil
		}
		return nil, res.Error
	}

	// convert
	alarms := []*trainalarm.TrainAlarm{}
	for _, result := range results {
		alarms = append(alarms, result.TrainAlarm())
	}

	return alarms, nil
}

// GetTrainAlarmsSortByLastNotificationAt with a limit. returns NO error if nothing found
func (s *SQLDatabase) GetTrainAlarmsSortByLastNotificationAt(ctx context.Context, limit int) ([]*trainalarm.TrainAlarm, error) {
	var results []SQLTrainAlarmModel
	if res := s.db.Limit(limit).Order("last_notification_at").Find(&results); res.Error != nil {
		if res.RecordNotFound() {
			return nil, nil
		}
		return nil, res.Error
	}

	// convert
	alarms := []*trainalarm.TrainAlarm{}
	for _, result := range results {
		alarms = append(alarms, result.TrainAlarm())
	}

	return alarms, nil
}

// DeleteTrainAlarm returns an error if alarm not found
func (s *SQLDatabase) DeleteTrainAlarm(ctx context.Context, id string) error {
	err := s.db.Where("id = ?", id).Delete(SQLTrainAlarmModel{}).Error
	return err
}

// DeleteOldTrainAlarms returns NO error if nothing to do
func (s *SQLDatabase) DeleteOldTrainAlarms(ctx context.Context, threshold time.Time) error {
	err := s.db.Where("final_arrival_at < ?", threshold).Delete(SQLTrainAlarmModel{}).Error
	return err
}

// UpdateTrainAlarm create a transaction, find the model and save it after updateFn. returns a error if model not found or pipe the error from updateFn
func (s *SQLDatabase) UpdateTrainAlarm(ctx context.Context, alarmID string, updateFn func(alarm *trainalarm.TrainAlarm) (*trainalarm.TrainAlarm, error)) error {
	var model SQLTrainAlarmModel
	if res := s.db.Where("id = ?", alarmID).Take(&model); res.Error != nil {
		return res.Error
	}

	trainAlarm, err := updateFn(model.TrainAlarm())
	if err != nil {
		return err
	}

	newModel := NewSQLTrainAlarmModel(trainAlarm)
	model.CreatedAt = newModel.CreatedAt
	model.UpdatedAt = newModel.UpdatedAt

	return s.db.Save(newModel).Error
}

// GetState returns an error if no state exists
func (s *SQLDatabase) GetState(ctx context.Context, identifyer string, plattform string) (string, string, error) {
	var state SQLStateModel
	if res := s.db.Where("identifyer = ? AND plattform = ?", identifyer, plattform).Take(&state); res.Error != nil {
		return "", "", res.Error
	}

	return state.State, state.StatePayload, nil
}

// UpdateState create a new if no exists
func (s *SQLDatabase) UpdateState(ctx context.Context, identifyer string, plattform string, updateFn func(state string, payload string) (string, string, error)) error {
	var model SQLStateModel
	isNew := false
	if res := s.db.Where("identifyer = ? AND plattform = ?", identifyer, plattform).Take(&model); res.Error != nil {
		if !res.RecordNotFound() {
			return res.Error
		}
		isNew = true
	}

	state := ""
	statePayload := ""
	if !isNew {
		state = model.State
		statePayload = model.StatePayload
	}

	state, statePayload, err := updateFn(state, statePayload)
	if err != nil {
		return err
	}

	// create / update
	if isNew {
		if err := s.db.Create(&SQLStateModel{
			Identifyer:   identifyer,
			Plattform:    plattform,
			State:        state,
			StatePayload: statePayload,
		}).Error; err != nil {
			return err
		}

		return nil
	}

	model.State = state
	model.StatePayload = statePayload

	return s.db.Save(model).Error
}

// DeleteOldStates returns NO error if nothing to do
func (s *SQLDatabase) DeleteOldStates(ctx context.Context, threshold time.Time) error {
	err := s.db.Where("updated_at < ?", threshold).Delete(SQLStateModel{}).Error
	return err
}

// Close database
func (s *SQLDatabase) Close() {
	s.db.Close()
}
