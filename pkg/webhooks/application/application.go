package application

import (
	"github.com/sirupsen/logrus"

	trainalarmDomain "github.com/pkuebler/bahn-bot/pkg/trainalarms/domain"
	webhookDomain "github.com/pkuebler/bahn-bot/pkg/webhooks/domain"
)

// Application represents all usecases
type Application struct {
	alarmRepo   trainalarmDomain.Repository
	webhookRepo webhookDomain.Repository
	log         *logrus.Entry
}

// NewApplication returns a application service object
func NewApplication(alarmRepo trainalarmDomain.Repository, webhookRepo webhookDomain.Repository, log *logrus.Entry) *Application {
	return &Application{
		alarmRepo:   alarmRepo,
		webhookRepo: webhookRepo,
		log:         log,
	}
}
