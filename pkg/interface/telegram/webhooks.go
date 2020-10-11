package telegram

import (
	"context"
	"fmt"

	"github.com/pkuebler/bahn-bot/pkg/infrastructure/telegramconversation"
)

// Webhooks managed webhook
func (t *TelegramService) Webhooks(ctx telegramconversation.TContext) telegramconversation.TContext {
	log := ctx.LogFields(t.log)
	log.Trace("Webhooks()")

	if ctx.IsButtonPressed() {
		ctx.DeleteMessage(ctx.MessageID())
	}

	webhooks, err := t.webhookRepository.GetWebhooksByIdentifyer(context.Background(), ctx.ChatID(), "telegram")
	if err != nil {
		log.Error(err)
		return ctx.SendWithState("Irgendetwas ist schief gelaufen. :/", "start")
	}

	if len(webhooks) == 0 {
		return ctx.SendWithState(`Du hast noch keine Webhooks angelegt.
		
Über den Webhook können automatisiert Alarme verwaltet werden, wenn z.B. über die Plattform travelynx.de eine aktive Zugfahrt getrackt wird.
		
/help`, "start")
	}

	txt := "Welcher Webhook soll bearbeitet werden?"
	buttons := []telegramconversation.TButton{}
	for _, hook := range webhooks {
		hookName := fmt.Sprintf("%s: %s", hook.GetProtocol(), hook.GetURLHash())
		button := telegramconversation.NewTButton(hookName, fmt.Sprintf("webhook|%s", hook.GetID()))
		buttons = append(buttons, button)
	}
	buttons = append(buttons, telegramconversation.NewTButton("Abbruch", "cancel"))

	return ctx.SendWithKeyboard(txt, buttons, 2)
}
