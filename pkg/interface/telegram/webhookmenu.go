package telegram

import (
	"context"
	"fmt"

	"github.com/pkuebler/bahn-bot/pkg/infrastructure/telegramconversation"
)

// WebhookMenu to select webhook options
func (t *TelegramService) WebhookMenu(ctx telegramconversation.TContext) telegramconversation.TContext {
	log := ctx.LogFields(t.log)
	log.Trace("WebhookMenu()")

	if !ctx.IsButtonPressed() {
		ctx.ChangeState("start")
		return ctx
	}

	ctx.DeleteMessage(ctx.MessageID())

	hook, err := t.webhookRepository.GetWebhook(context.Background(), ctx.ButtonData())
	if err != nil {
		log.Error(err)
		return ctx.SendWithState("Webhook nicht gefunden.", "start")
	}

	txt := fmt.Sprintf("Was möchtest du für `%s: %s` ändern?", hook.GetProtocol(), hook.GetURLHash())
	buttons := []telegramconversation.TButton{
		telegramconversation.NewTButton("Löschen", fmt.Sprintf("deletewebhook|%s", hook.GetID())),
		telegramconversation.NewTButton("Zurück zur Liste", "listwebhooks"),
	}

	return ctx.SendWithKeyboard(txt, buttons, 2)
}
