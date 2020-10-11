package telegram

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/pkuebler/bahn-bot/pkg/config"
	"github.com/pkuebler/bahn-bot/pkg/infrastructure/marudor"
	trainalarmApplication "github.com/pkuebler/bahn-bot/pkg/trainalarms/application"
	trainalarmDomain "github.com/pkuebler/bahn-bot/pkg/trainalarms/domain"
	webhookApplication "github.com/pkuebler/bahn-bot/pkg/webhooks/application"
	webhookDomain "github.com/pkuebler/bahn-bot/pkg/webhooks/domain"
)

// Application with business logic
type Application interface {
	DeleteTrainAlarm(ctx context.Context, cmd trainalarmApplication.DeleteTrainAlarmCmd) (*trainalarmDomain.TrainAlarm, error)
	AddTrainAlarm(ctx context.Context, cmd trainalarmApplication.AddTrainAlarmCmd) error
	UpdateTrainAlarmThreshold(ctx context.Context, cmd trainalarmApplication.UpdateTrainAlarmThresholdCmd) error
	AddWebhook(ctx context.Context, cmd webhookApplication.AddWebhookCmd) (*webhookDomain.Webhook, error)
	DeleteWebhook(ctx context.Context, cmd webhookApplication.DeleteWebhookCmd) (*webhookDomain.Webhook, error)
}

// HafasService to request train informations
type HafasService interface {
	FindTrain(ctx context.Context, trainName string, date time.Time) (*[]marudor.HafasTrainResult, error)
}

// TelegramService to handle requests
type TelegramService struct {
	log                  *logrus.Entry
	config               *config.Config
	trainAlarmRepository trainalarmDomain.Repository
	trainalarmApp        *trainalarmApplication.Application
	webhookRepository    webhookDomain.Repository
	webhookApp           *webhookApplication.Application
	hafas                HafasService
}

// NewTelegramService to create a new service
func NewTelegramService(
	log *logrus.Entry,
	cfg *config.Config,
	trainAlarmRepository trainalarmDomain.Repository,
	webhookRepository webhookDomain.Repository,
	webhookApp *webhookApplication.Application,
	trainalarmApp *trainalarmApplication.Application,
	hafas HafasService,
) *TelegramService {
	return &TelegramService{
		log:                  log.WithField("service", "telegram"),
		config:               cfg,
		trainAlarmRepository: trainAlarmRepository,
		trainalarmApp:        trainalarmApp,
		webhookRepository:    webhookRepository,
		webhookApp:           webhookApp,
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
