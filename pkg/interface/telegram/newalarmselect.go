package telegram

import (
	"context"
	"fmt"
	"time"

	"github.com/pkuebler/bahn-bot/pkg/infrastructure/telegramconversation"
)

// NewAlarmSelect parse user input and show buttons to select a train
func (t *TelegramService) NewAlarmSelect(ctx telegramconversation.TContext) telegramconversation.TContext {
	log := ctx.LogFields(t.log)
	log.Trace("NewAlarmSelect()")

	// search train at marudor
	results, err := t.hafas.FindTrain(context.Background(), ctx.Message(), time.Now())
	if err != nil {
		txt := "Zug nicht gefunden :/."
		t.log.Error(err)
		return ctx.SendWithState(txt, "start")
	}

	buttons := []telegramconversation.TButton{}
	for _, train := range *results {
		trainName := fmt.Sprintf("%s > %s", train.Train.Name, train.LastStop.Station.Title)
		button := telegramconversation.NewTButton(trainName, fmt.Sprintf("savealarm|%s %s|%s|%d", train.Train.Type, train.Train.Number, train.FirstStop.Station.ID, train.FirstStop.Departure.ScheduledTime))
		buttons = append(buttons, button)
	}
	buttons = append(buttons, telegramconversation.NewTButton("Abbruch", "cancel"))

	txt := "Wähle den passenden Zug:"
	ctx.ChangeState("start")
	return ctx.SendWithKeyboard(txt, buttons, 2)
}
