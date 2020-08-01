package telegram

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/pkuebler/bahn-bot/pkg/application"
	"github.com/pkuebler/bahn-bot/pkg/infrastructure/telegramconversation"
)

var (
	delayRegex = regexp.MustCompile(`^(([0-9]+)h)?(([0-9]+)m)?$`)
)

// SaveDelay to database
func (t *TelegramService) SaveDelay(ctx telegramconversation.TContext) telegramconversation.TContext {
	log := ctx.LogFields(t.log).WithField("alarmID", ctx.StatePayload())
	log.Trace("SaveDelay()")

	// state alarm id missing
	if ctx.StatePayload() == "" {
		log.Error("no alarm id")
		return ctx.SendWithState("Irgendetwas ist schief gelaufen. :/", "start")
	}

	if !delayRegex.MatchString(ctx.Message()) {
		log.Trace("wrong time")
		return ctx.SendWithState("Das Format wurde nicht verstanden. Abgebrochen.", "start")
	}

	parts := delayRegex.FindStringSubmatch(ctx.Message())
	hours, err := strconv.Atoi(parts[2])
	if parts[2] != "" && err != nil {
		log.Error(err)
		return ctx.SendWithState("Irgendetwas ist schief gelaufen. :/", "start")
	}
	minutes, err := strconv.Atoi(parts[4])
	if parts[4] != "" && err != nil {
		log.Error(err)
		return ctx.SendWithState("Irgendetwas ist schief gelaufen. :/", "start")
	}

	thresholdMinutes := time.Duration(minutes)*time.Minute + time.Duration(hours)*time.Hour

	log.Trace("Set delay to ", thresholdMinutes)

	cmd := application.UpdateTrainAlarmThresholdCmd{
		AlarmID:          ctx.StatePayload(),
		ThresholdMinutes: int(thresholdMinutes.Minutes()),
	}

	err = t.application.UpdateTrainAlarmThreshold(context.Background(), cmd)
	if err != nil {
		log.Error(err)
		return ctx.SendWithState("Irgendetwas ist schief gelaufen. :/", "start")
	}

	txt := fmt.Sprintf("Erlaube Abweichung liegt nun bei %d Minuten.", cmd.ThresholdMinutes)
	return ctx.SendWithState(txt, "start")
}
