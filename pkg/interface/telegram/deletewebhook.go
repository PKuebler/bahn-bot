package telegram

import (
	"context"
	"fmt"

	"github.com/pkuebler/bahn-bot/pkg/infrastructure/telegramconversation"
	"github.com/pkuebler/bahn-bot/pkg/webhooks/application"
)

// DeleteWebhook from database
func (t *TelegramService) DeleteWebhook(ctx telegramconversation.TContext) telegramconversation.TContext {
	log := ctx.LogFields(t.log)
	log.Trace("DeleteWebhook()")

	if !ctx.IsButtonPressed() {
		return ctx.SendWithState("Irgendetwas ist schief gelaufen. :/", "start")
	}

	ctx.DeleteMessage(ctx.MessageID())

	cmd := application.DeleteWebhookCmd{
		WebhookID: ctx.ButtonData(),
	}
	hook, _ := t.webhookApp.DeleteWebhook(context.Background(), cmd)

	return ctx.SendWithState(fmt.Sprintf("Webhook `%s: %s` gelöscht.", hook.GetProtocol(), hook.GetURLHash()), "start")
}
