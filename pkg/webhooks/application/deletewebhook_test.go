package application

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	webhookDomain "github.com/pkuebler/bahn-bot/pkg/webhooks/domain"
)

func TestDeleteWebhook(t *testing.T) {
	app, alarmRepo, webhookRepo := createTestCase(true)
	ctx := context.Background()

	hook, err := webhookDomain.NewWebhook(uuid.New().String(), "telegram", "12345", "12345", webhookDomain.TravelynxProtocol)
	assert.Nil(t, err)
	assert.NotNil(t, hook)

	cmd := DeleteWebhookCmd{
		WebhookID: hook.GetID(),
	}

	// entry not found
	a, err := app.DeleteWebhook(ctx, cmd)
	assert.NotNil(t, err)
	assert.Nil(t, a)

	// add entry
	webhookRepo.Webhooks[hook.GetID()] = hook

	// entry deleted
	a, err = app.DeleteWebhook(ctx, cmd)
	assert.Nil(t, err)
	assert.NotNil(t, a)
	assert.Len(t, webhookRepo.Webhooks, 0)
}
