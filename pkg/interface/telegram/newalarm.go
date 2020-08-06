package telegram

import "github.com/pkuebler/bahn-bot/pkg/infrastructure/telegramconversation"

// NewAlarm request train name
func (t *TelegramService) NewAlarm(ctx telegramconversation.TContext) telegramconversation.TContext {
	log := ctx.LogFields(t.log)
	log.Trace("NewAlarm()")

	if ctx.CommandQuery() != "" {
		ctx.SetMessage(ctx.CommandQuery())
		return t.NewAlarmSelect(ctx)
	}

	return ctx.SendWithState("Welcher Zug soll überwacht werden?", "newalarmselect")
}
