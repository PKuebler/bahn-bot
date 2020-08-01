package telegram

import (
	"context"
	"fmt"

	"github.com/pkuebler/bahn-bot/pkg/infrastructure/telegramconversation"
)

// ListTrainAlarms command show current listalarmsings
func (t *TelegramService) ListTrainAlarms(ctx telegramconversation.TContext) telegramconversation.TContext {
	log := ctx.LogFields(t.log)
	log.Trace("ListTrainAlarms()")

	if ctx.IsButtonPressed() {
		ctx.DeleteMessage(ctx.MessageID())
	}

	// list trains
	alarms, err := t.trainAlarmRepository.GetTrainAlarms(context.Background(), ctx.ChatID(), "telegram")
	if err != nil {
		log.Error(err)
		return ctx.SendWithState("Irgendetwas ist schief gelaufen. :/", "start")
	}

	if len(alarms) == 0 {
		if ctx.IsButtonPressed() {
			ctx.DeleteMessage(ctx.MessageID())
		}
		return ctx.SendWithState("Du beobachtest noch keine Züge. /help", "start")
	}

	txt := "Welcher Alarm soll bearbeitet werden?"
	buttons := []telegramconversation.TButton{}
	for _, alarm := range alarms {
		button := telegramconversation.NewTButton(alarm.GetTrainName(), fmt.Sprintf("alarm|%s", alarm.GetID()))
		buttons = append(buttons, button)
	}
	buttons = append(buttons, telegramconversation.NewTButton("Abbruch", "cancel"))

	return ctx.SendWithKeyboard(txt, buttons, 2)
}
