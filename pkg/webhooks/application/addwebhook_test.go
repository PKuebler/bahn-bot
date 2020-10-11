package application

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	webhookDomain "github.com/pkuebler/bahn-bot/pkg/webhooks/domain"
)

func TestAddWebhook(t *testing.T) {
	app, alarmRepo, webhookRepo := createTestCase(true)
	ctx := context.Background()

	cmd := AddWebhookCmd{
		Identifyer: uuid.New().String(),
		Plattform:  "telegram",
		Protocol:   string(webhookDomain.TravelynxProtocol),
	}

	_, err := app.AddWebhook(ctx, cmd)
	assert.Nil(t, err)
	assert.Len(t, webhookRepo.Webhooks, 1)
}
