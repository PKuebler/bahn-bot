package telegram

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/pkuebler/bahn-bot/pkg/application"
	"github.com/pkuebler/bahn-bot/pkg/domain/trainalarm"
	"github.com/pkuebler/bahn-bot/pkg/infrastructure/marudor"
	"github.com/sirupsen/logrus"
)

// Application with business logic
type Application interface {
	DeleteTrainAlarm(ctx context.Context, cmd application.DeleteTrainAlarmCmd) (*trainalarm.TrainAlarm, error)
	AddTrainAlarm(ctx context.Context, cmd application.AddTrainAlarmCmd) error
	UpdateTrainAlarmThreshold(ctx context.Context, cmd application.UpdateTrainAlarmThresholdCmd) error
}

// HafasService to request train informations
type HafasService interface {
	FindTrain(ctx context.Context, trainName string, date time.Time) (*[]marudor.HafasTrainResult, error)
}

// TelegramService to handle requests
type TelegramService struct {
	log                  *logrus.Entry
	trainAlarmRepository trainalarm.Repository
	application          Application
	hafas                HafasService
}

// NewTelegramService to create a new service
func NewTelegramService(log *logrus.Entry, repository trainalarm.Repository, application Application, hafas HafasService) *TelegramService {
	return &TelegramService{
		log:                  log.WithField("service", "telegram"),
		trainAlarmRepository: repository,
		application:          application,
		hafas:                hafas,
	}
}

// ParseButtonQuery to trainName, stationEVA and hafasDate
func ParseButtonQuery(query string) (string, int, int64, error) {
	// parse query to parameters
	parts := strings.Split(query, "|")

	if len(parts) != 3 {
		return "", 0, 0, errors.New("wrong query")
	}

	trainName := parts[0]
	stationEVA, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", 0, 0, err
	}
	stationDate, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return "", 0, 0, err
	}

	return trainName, stationEVA, stationDate, nil
}
