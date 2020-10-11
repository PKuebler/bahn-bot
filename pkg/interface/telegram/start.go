package telegram

import (
	"fmt"

	"github.com/pkuebler/bahn-bot/pkg/infrastructure/telegramconversation"
)

// Start new conversation and/or print help
func (t *TelegramService) Start(ctx telegramconversation.TContext) telegramconversation.TContext {
	log := ctx.LogFields(t.log)

	if !ctx.IsCommand("start") && !ctx.IsCommand("help") {
		// redirection to start (no user interaction)
		return ctx
	}

	log.Trace("Start()")

	txt := fmt.Sprintf(`Telegrambot für automatische Verspätungsalarme.

/help Diese Befehlsübersicht

/myalarms Bearbeite deine Alarme
/newalarm Erzeuge neuen Alarm

/webhooks Bearbeite deine Webhooks
/newwebhook Verknüpfe travelynx.de

/cancel Breche aktuelle Option ab `)

	ctx.Send(txt)

	return ctx
}
