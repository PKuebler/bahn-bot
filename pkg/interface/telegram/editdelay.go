package telegram

import (
	"context"
	"fmt"

	"github.com/pkuebler/bahn-bot/pkg/infrastructure/telegramconversation"
)

// EditDelay request train alarm threshold
func (t *TelegramService) EditDelay(ctx telegramconversation.TContext) telegramconversation.TContext {
	log := ctx.LogFields(t.log)
	log.Trace("EditDelay()")

	if !ctx.IsButtonPressed() {
		return ctx.SendWithState("Irgendetwas ist schief gelaufen. :/", "start")
	}

	ctx.DeleteMessage(ctx.MessageID())

	alarm, err := t.trainAlarmRepository.GetTrainAlarm(context.Background(), ctx.ButtonData())
	if err != nil {
		log.Error(err)
		return ctx.SendWithState("Irgendetwas ist schief gelaufen. :/", "start")
	}

	if alarm == nil {
		log.Error("not found")
		return ctx.SendWithState("Alarm nicht gefunden.", "start")
	}

	ctx.SetStatePayload(alarm.GetID())
	ctx.ChangeState("savedelay")

	txt := fmt.Sprintf("Ab wie viel Abweichung von `%s > %s` soll bescheid gesagt werden? z.B. 1h5min oder 10m", alarm.GetTrainName(), alarm.GetFinalDestinationName())

	buttons := []telegramconversation.TButton{}
	units := []string{"1m", "5m", "10m", "15m", "20m", "30m", "1h", "1h30m"}
	for _, unit := range units {
		button := telegramconversation.NewTButton(unit, unit)
		buttons = append(buttons, button)
	}

	return ctx.SendWithSuggestions(txt, buttons, 1)
}
