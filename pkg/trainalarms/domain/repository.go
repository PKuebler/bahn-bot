package domain

import (
	"context"
	"time"
)

// Repository interface to datahandling between repository and application
type Repository interface {
	GetOrCreateTrainAlarm(ctx context.Context, alarm *TrainAlarm) (*TrainAlarm, error)
	GetTrainAlarm(ctx context.Context, alarmID string) (*TrainAlarm, error)
	GetTrainAlarms(ctx context.Context, identifyer string, plattform string) ([]*TrainAlarm, error)
	GetTrainAlarmsSortByLastNotificationAt(ctx context.Context, limit int) ([]*TrainAlarm, error)
	DeleteTrainAlarm(ctx context.Context, alarmID string) error
	DeleteOldTrainAlarms(ctx context.Context, threshold time.Time) error
	DeleteOldStates(ctx context.Context, threshold time.Time) error
	UpdateTrainAlarm(ctx context.Context, alarmID string, updateFn func(alarm *TrainAlarm) (*TrainAlarm, error)) error
}
