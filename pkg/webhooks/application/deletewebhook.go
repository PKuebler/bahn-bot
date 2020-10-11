package application

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"

	webhookDomain "github.com/pkuebler/bahn-bot/pkg/webhooks/domain"
)

// DeleteWebhookCmd for DeleteWebhook
type DeleteWebhookCmd struct {
	WebhookID string
}

// DeleteWebhook at database
func (a *Application) DeleteWebhook(ctx context.Context, cmd DeleteWebhookCmd) (*webhookDomain.Webhook, error) {
	log := a.log.WithFields(logrus.Fields{
		"webhookID": cmd.WebhookID,
	})

	// search webhook at repository
	webhook, err := a.webhookRepo.GetWebhook(ctx, cmd.WebhookID)
	if err != nil {
		log.Error(err)
		return nil, errors.New("internal server error")
	}
	if webhook == nil {
		log.Trace("not found")
		return nil, errors.New("not found")
	}

	// delete webhook at repository
	if err := a.webhookRepo.DeleteWebhook(ctx, webhook.GetID()); err != nil {
		log.Error(err)
		return nil, errors.New("internal server error")
	}

	log.Trace("webhook deleted")
	return webhook, nil
}
