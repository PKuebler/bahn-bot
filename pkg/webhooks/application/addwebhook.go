package application

import (
	"context"
	"errors"

	"github.com/lithammer/shortuuid"
	"github.com/sirupsen/logrus"

	webhookDomain "github.com/pkuebler/bahn-bot/pkg/webhooks/domain"
)

// AddWebhookCmd for AddWebhook
type AddWebhookCmd struct {
	Identifyer string
	Plattform  string
	Protocol   string
}

// AddWebhook to database
func (a *Application) AddWebhook(ctx context.Context, cmd AddWebhookCmd) (*webhookDomain.Webhook, error) {
	log := a.log.WithFields(logrus.Fields{
		"identifyer": cmd.Identifyer,
		"plattform":  cmd.Plattform,
		"protocol":   cmd.Protocol,
	})

	protocol, err := webhookDomain.NewWebhookProtocol(cmd.Protocol)
	if err != nil {
		log.Error(err)
		return nil, errors.New("unknown protocol")
	}

	token := shortuuid.New()

	hook, err := webhookDomain.NewWebhook(
		cmd.Identifyer,
		cmd.Plattform,
		token,
		protocol,
	)

	_, err = a.webhookRepo.CreateWebhook(ctx, hook)
	if err != nil {
		log.Error(err)
		return nil, errors.New("internal server error")
	}

	log.Trace("webhook created")
	return hook, nil
}
