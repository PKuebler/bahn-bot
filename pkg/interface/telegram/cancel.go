package telegram

import (
	"github.com/pkuebler/bahn-bot/pkg/infrastructure/telegramconversation"
)

// Cancel conversation state
func (t *TelegramService) Cancel(ctx telegramconversation.TContext) telegramconversation.TContext {
	log := ctx.LogFields(t.log)
	log.Trace("Cancel()")

	if ctx.IsButtonPressed() {
		ctx.DeleteMessage(ctx.MessageID())
	}

	if !ctx.IsCommand("cancel") {
		// redirection to cancel (no user interaction)
		return ctx
	}

	txt := "Abgebrochen. /help"
	return ctx.SendWithState(txt, "start")
}
