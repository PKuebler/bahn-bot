package telegram

import (
	"context"
	"fmt"
	"net/url"
	"path"

	"github.com/pkuebler/bahn-bot/pkg/infrastructure/telegramconversation"
	"github.com/pkuebler/bahn-bot/pkg/webhooks/application"
	"github.com/pkuebler/bahn-bot/pkg/webhooks/domain"
)

// NewWebhook generates a webhook that Travelynx uses to automatically set alarms for current trains.
func (t *TelegramService) NewWebhook(ctx telegramconversation.TContext) telegramconversation.TContext {
	log := ctx.LogFields(t.log)
	log.Trace("NewWebhook()")

	cmd := application.AddWebhookCmd{
		Identifyer: ctx.ChatID(),
		Plattform:  "telegram",
		Protocol:   string(domain.TravelynxProtocol),
	}

	hook, err := t.webhookApp.AddWebhook(context.Background(), cmd)
	if err != nil {
		log.Error(err)
		return ctx.SendWithState("Irgendetwas ist schief gelaufen. :/", "start")
	}

	u, _ := url.Parse(t.config.Webhook.Endpoint)
	u.Path = path.Join(u.Path, hook.GetURLHash())
	txt := fmt.Sprintf(`Neuer Webhook angelegt.

	Um automatisch Alarme angelegt zu bekommen, müssen unter https://travelynx.de/account/hooks folgende Daten hinterlegt werden:

	URL: `+"`%s`", u.String())
	return ctx.SendWithState(txt, "start")
}
