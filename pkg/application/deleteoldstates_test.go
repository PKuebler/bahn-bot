package application

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func (a *Application) TestDeleteOldStates(t *testing.T) error {
	app, _ := createTestCase(true)
	ctx := context.Background()

	err := app.DeleteOldStates(ctx)
	assert.Nil(t, err)

	return nil
}
