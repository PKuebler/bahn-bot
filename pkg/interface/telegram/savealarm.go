package telegram

import (
	"context"
	"fmt"

	"github.com/pkuebler/bahn-bot/pkg/infrastructure/telegramconversation"
	"github.com/pkuebler/bahn-bot/pkg/trainalarms/application"
)

// SaveAlarm to database
func (t *TelegramService) SaveAlarm(ctx telegramconversation.TContext) telegramconversation.TContext {
	log := ctx.LogFields(t.log)
	log.Trace("SaveAlarm()")

	if !ctx.IsButtonPressed() {
		ctx.ChangeState("start")
		return ctx
	}

	ctx.DeleteMessage(ctx.MessageID())

	trainName, stationEVA, stationDate, err := ParseButtonQuery(ctx.ButtonData())
	if err != nil {
		log.Error(err)
		return ctx.SendWithState("Irgendetwas ist schief gelaufen. :/", "start")
	}

	log.Tracef("Train %s, EVA %d, Date %d", trainName, stationEVA, stationDate)

	cmd := application.AddTrainAlarmCmd{
		Identifyer:  ctx.ChatID(),
		Plattform:   "telegram",
		TrainName:   trainName,
		StationEVA:  stationEVA,
		StationDate: stationDate,
	}

	err = t.trainalarmApp.AddTrainAlarm(context.Background(), cmd)
	if err != nil {
		log.Error(err)
		return ctx.SendWithState("Irgendetwas ist schief gelaufen. :/", "start")
	}

	txt := fmt.Sprintf("Neuer Alarm `%s` hinzugefügt. Über /myalarms kann die erlaubte Abweichung geändert werden.", trainName)
	return ctx.SendWithState(txt, "start")
}
