package telegram

import (
	"context"
	"fmt"

	"github.com/pkuebler/bahn-bot/pkg/infrastructure/telegramconversation"
)

// AlarmMenu to select alarm options
func (t *TelegramService) AlarmMenu(ctx telegramconversation.TContext) telegramconversation.TContext {
	log := ctx.LogFields(t.log)
	log.Trace("AlarmMenu()")

	if !ctx.IsButtonPressed() {
		ctx.ChangeState("start")
		return ctx
	}

	ctx.DeleteMessage(ctx.MessageID())

	alarm, err := t.trainAlarmRepository.GetTrainAlarm(context.Background(), ctx.ButtonData())
	if err != nil {
		log.Error(err)
		return ctx.SendWithState("Alarm nicht gefunden.", "start")
	}

	txt := fmt.Sprintf("Was möchtest du für %s (Alarm ab %dm) ändern?", alarm.GetTrainName(), alarm.GetDelayThresholdMinutes())
	buttons := []telegramconversation.TButton{
		telegramconversation.NewTButton("Alarm ab ...", fmt.Sprintf("editdelay|%s", alarm.GetID())),
		telegramconversation.NewTButton("Löschen", fmt.Sprintf("deletealarm|%s", alarm.GetID())),
		telegramconversation.NewTButton("Zurück zur Liste", "listalarms"),
	}

	return ctx.SendWithKeyboard(txt, buttons, 2)
}
