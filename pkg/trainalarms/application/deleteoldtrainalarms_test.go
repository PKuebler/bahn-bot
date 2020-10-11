package application

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteOldTrainAlarms(t *testing.T) {
	app, _ := createTestCase(true)
	ctx := context.Background()

	err := app.DeleteOldTrainAlarms(ctx)
	assert.Nil(t, err)
}
