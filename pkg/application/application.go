package application

import (
	"context"

	"github.com/pkuebler/bahn-bot/pkg/domain/trainalarm"
	"github.com/pkuebler/bahn-bot/pkg/infrastructure/marudor"
	"github.com/sirupsen/logrus"
)

// Train by hafas
type Train struct {
}

// HafasService to request train informations
type HafasService interface {
	GetTrainByStation(ctx context.Context, trainName string, stationEVA int, stationDate int64) (*marudor.HafasTrain, error)
}

// Application represents all usecases
type Application struct {
	hafas HafasService
	repo  trainalarm.Repository
	log   *logrus.Entry
}

// NewApplication returns a application service object
func NewApplication(hafas HafasService, repo trainalarm.Repository, log *logrus.Entry) *Application {
	return &Application{
		hafas: hafas,
		repo:  repo,
		log:   log,
	}
}
