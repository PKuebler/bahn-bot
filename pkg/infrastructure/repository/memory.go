package repository

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/pkuebler/bahn-bot/pkg/domain/trainalarm"
)

// MemoryStateModel save the current state
type MemoryStateModel struct {
	Identifyer    string
	Plattform     string
	State         string
	StatePlayload string
}

// MemoryDatabase to test without persistence
type MemoryDatabase struct {
	TrainAlarms  map[string]*trainalarm.TrainAlarm
	CurrentState []*MemoryStateModel
	Mutex        *sync.Mutex
}

// NewMemoryDatabase with mutex
func NewMemoryDatabase() *MemoryDatabase {
	return &MemoryDatabase{
		TrainAlarms:  map[string]*trainalarm.TrainAlarm{},
		CurrentState: []*MemoryStateModel{},
		Mutex:        &sync.Mutex{},
	}
}

// GetOrCreateTrainAlarm creates a new alarm if none exists yet
func (m *MemoryDatabase) GetOrCreateTrainAlarm(ctx context.Context, alarm *trainalarm.TrainAlarm) (*trainalarm.TrainAlarm, error) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	oldAlarm, err := m.findTrainAlarm(alarm)
	if err != nil {
		return nil, err
	}

	if oldAlarm == nil {
		m.TrainAlarms[alarm.GetID()] = alarm
		return alarm, nil
	}

	return oldAlarm, nil
}

// GetTrainAlarm by id. returns NO error if nothing found
func (m *MemoryDatabase) GetTrainAlarm(ctx context.Context, alarmID string) (*trainalarm.TrainAlarm, error) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	if _, ok := m.TrainAlarms[alarmID]; !ok {
		return nil, nil
	}

	return m.TrainAlarms[alarmID], nil
}

// GetTrainAlarms by identifyer and plattform. returns NO error if nothing found
func (m *MemoryDatabase) GetTrainAlarms(ctx context.Context, identifyer string, plattform string) ([]*trainalarm.TrainAlarm, error) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	alarms := []*trainalarm.TrainAlarm{}

	for _, alarm := range m.TrainAlarms {
		if identifyer == alarm.GetIdentifyer() && plattform == alarm.GetPlattform() {
			alarms = append(alarms, alarm)
		}
	}

	return alarms, nil
}

// GetTrainAlarmsSortByLastNotificationAt with a limit. returns NO error if nothing found
func (m *MemoryDatabase) GetTrainAlarmsSortByLastNotificationAt(ctx context.Context, limit int) ([]*trainalarm.TrainAlarm, error) {
	results := []*trainalarm.TrainAlarm{}

	// convert map to slice
	alarms := []*trainalarm.TrainAlarm{}

	func() {
		m.Mutex.Lock()
		defer m.Mutex.Unlock()

		for _, alarm := range m.TrainAlarms {
			alarms = append(alarms, alarm)
		}
	}()

	// sort map, nil is bigger then date
	sort.Slice(alarms, func(i, j int) bool {
		// nil is bigger than other
		if alarms[i].GetLastNotificationAt() == nil {
			if alarms[j].GetLastNotificationAt() == nil {
				return true
			}

			return false
		}
		if alarms[j].GetLastNotificationAt() == nil {
			return true
		}

		return alarms[i].GetLastNotificationAt().Before(*alarms[j].GetLastNotificationAt())
	})

	limit = min(limit, len(alarms))

	for i := 0; i < limit; i++ {
		results = append(results, alarms[i])
	}

	return results, nil
}

// DeleteTrainAlarm returns an error if alarm not found
func (m *MemoryDatabase) DeleteTrainAlarm(ctx context.Context, id string) error {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	if _, ok := m.TrainAlarms[id]; !ok {
		return errors.New("alarm not found")
	}

	delete(m.TrainAlarms, id)

	return nil
}

// DeleteOldTrainAlarms returns NO error if nothing to do
func (m *MemoryDatabase) DeleteOldTrainAlarms(ctx context.Context, threshold time.Time) error {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	for key, alarm := range m.TrainAlarms {
		if alarm.GetFinalArrivalAt().Before(threshold) {
			delete(m.TrainAlarms, key)
		}
	}

	return nil
}

// UpdateTrainAlarm create a transaction, find the model and save it after updateFn. returns a error if model not found or pipe the error from updateFn
func (m *MemoryDatabase) UpdateTrainAlarm(ctx context.Context, alarmID string, updateFn func(alarm *trainalarm.TrainAlarm) (*trainalarm.TrainAlarm, error)) error {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	if _, ok := m.TrainAlarms[alarmID]; !ok {
		return errors.New("not found")
	}

	alarm, err := updateFn(m.TrainAlarms[alarmID])
	if err != nil {
		return err
	}

	m.TrainAlarms[alarmID] = alarm
	return nil
}

// findTrainAlarm returns NO error if nothing found
func (m *MemoryDatabase) findTrainAlarm(query *trainalarm.TrainAlarm) (*trainalarm.TrainAlarm, error) {
	for _, alarm := range m.TrainAlarms {
		if alarm.Compare(query) {
			return alarm, nil
		}
	}
	return nil, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GetState returns an error if no state exists
func (m *MemoryDatabase) GetState(ctx context.Context, identifyer string, plattform string) (string, string, error) {
	for _, model := range m.CurrentState {
		if model.Identifyer == identifyer && model.Plattform == plattform {
			return model.State, model.StatePlayload, nil
		}
	}

	return "", "", errors.New("not found")
}

// UpdateState create a new if no exists
func (m *MemoryDatabase) UpdateState(ctx context.Context, identifyer string, plattform string, updateFn func(state string, payload string) (string, string, error)) error {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	// find
	var state *MemoryStateModel

	for _, model := range m.CurrentState {
		if model.Identifyer == identifyer && model.Plattform == plattform {
			state = model
			break
		}
	}

	if state == nil {
		state = &MemoryStateModel{
			Identifyer:    identifyer,
			Plattform:     plattform,
			State:         "",
			StatePlayload: "",
		}
		m.CurrentState = append(m.CurrentState, state)
	}

	newState, newPayload, err := updateFn(state.State, state.StatePlayload)
	if err != nil {
		return err
	}

	state.State = newState
	state.StatePlayload = newPayload

	return nil
}
