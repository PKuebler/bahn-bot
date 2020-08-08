package telegram

import (
	"context"
	"fmt"

	"github.com/pkuebler/bahn-bot/pkg/application"
	"github.com/pkuebler/bahn-bot/pkg/infrastructure/telegramconversation"
)

// DeleteAlarm from database
func (t *TelegramService) DeleteAlarm(ctx telegramconversation.TContext) telegramconversation.TContext {
	log := ctx.LogFields(t.log)
	log.Trace("DeleteAlarm()")

	if !ctx.IsButtonPressed() {
		return ctx.SendWithState("Irgendetwas ist schief gelaufen. :/", "start")
	}

	ctx.DeleteMessage(ctx.MessageID())

	cmd := application.DeleteTrainAlarmCmd{
		AlarmID: ctx.ButtonData(),
	}
	alarm, _ := t.application.DeleteTrainAlarm(context.Background(), cmd)

	return ctx.SendWithState(fmt.Sprintf("Alarm `%s > %s` gelöscht.", alarm.GetTrainName(), alarm.GetFinalDestinationName()), "start")
}
