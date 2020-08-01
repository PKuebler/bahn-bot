package application

import (
	"github.com/sirupsen/logrus"

	"github.com/pkuebler/bahn-bot/pkg/infrastructure/marudor"
	"github.com/pkuebler/bahn-bot/pkg/infrastructure/repository"
)

func createTestCase(trainExists bool) (*Application, *repository.MemoryDatabase) {
	log := logrus.NewEntry(logrus.StandardLogger())

	hafas := &marudor.HafasMock{TrainExists: trainExists}
	repo := repository.NewMemoryDatabase()
	app := NewApplication(hafas, repo, log)

	return app, repo
}
