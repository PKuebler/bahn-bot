package application

import (
	"github.com/sirupsen/logrus"

	trainalarmRepository "github.com/pkuebler/bahn-bot/pkg/trainalarms/repository"
	webhookRepository "github.com/pkuebler/bahn-bot/pkg/webhooks/repository"
)

func createTestCase(trainExists bool) (*Application, *trainalarmRepository.MemoryDatabase, *webhookRepository.MemoryDatabase) {
	log := logrus.NewEntry(logrus.StandardLogger())

	webhookRepo := webhookRepository.NewMemoryDatabase()
	trainalarmRepo := trainalarmRepository.NewMemoryDatabase()
	app := NewApplication(trainalarmRepo, webhookRepo, log)

	return app, trainalarmRepo, webhookRepo
}
