package cmd

import (
	"context"
	"net/http"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/pkuebler/bahn-bot/pkg/application"
	"github.com/pkuebler/bahn-bot/pkg/config"
	"github.com/pkuebler/bahn-bot/pkg/infrastructure/marudor"
	"github.com/pkuebler/bahn-bot/pkg/infrastructure/repository"
	"github.com/pkuebler/bahn-bot/pkg/infrastructure/telegramconversation"
	"github.com/pkuebler/bahn-bot/pkg/interface/cron"
	"github.com/pkuebler/bahn-bot/pkg/interface/telegram"
	"github.com/pkuebler/bahn-bot/pkg/metrics"
)

// NewBotCmd create a command to start the bot
func NewBotCmd(ctx context.Context) *cobra.Command {
	command := &cobra.Command{
		Use:     "bot",
		Short:   "start bot",
		Long:    "start the bahn bot",
		Example: "bot",
		Run: func(cmd *cobra.Command, args []string) {
			BotCommand(ctx, cmd, args)
		},
	}

	return command
}

// BotCommand starts the bot
func BotCommand(ctx context.Context, cmd *cobra.Command, args []string) {
	log := logrus.NewEntry(logrus.StandardLogger())

	cfg := config.ReadConfig("config.json", log)

	switch cfg.LogLevel {
	case "info":
		log.Logger.Level = logrus.InfoLevel
	case "trace":
		log.Logger.Level = logrus.TraceLevel
	case "warn":
		log.Logger.Level = logrus.WarnLevel
	case "error":
		log.Logger.Level = logrus.ErrorLevel
	default:
		log.Logger.Level = logrus.InfoLevel
	}

	// metrics
	metricsRegistry := metrics.NewPrometheusMetric()

	// external interfaces
	api, err := marudor.NewAPIClient(cfg.APIConfig.APIEndpoint, nil, log, cfg.EnableMetrics)
	if err != nil {
		panic(err)
	}
	hafas := api.HafasService

	// storage
	// repo := repository.NewMemoryDatabase()
	repo := repository.NewSQLDatabase(cfg.Database.Dialect, cfg.Database.Path)
	defer repo.Close()

	// application
	app := application.NewApplication(hafas, repo, log)

	// interfaces
	service := telegram.NewTelegramService(log, repo, app, hafas)
	cronService := cron.NewCronJob(log, app, cfg.EnableMetrics)

	// conversationengine
	router := telegramconversation.NewConversationRouter("start")
	router.OnUnknownState(func(ctx telegramconversation.TContext) telegramconversation.TContext {
		log.Tracef("[%d] State `%s` nicht gefunden.", ctx.MessageID(), ctx.State())
		return ctx.SendWithState("Irgendwas ist schief gelaufen :/", "start")
	})

	router.OnCommand("start", "start")
	router.OnCommand("help", "start")
	router.OnState("start", service.Start)
	router.OnState("cancel", service.Cancel)

	router.OnCommand("myalarms", "listalarms")
	router.OnState("listalarms", service.ListTrainAlarms)
	router.OnState("alarm", service.AlarmMenu)
	router.OnState("deletealarm", service.DeleteAlarm)

	router.OnCommand("newalarm", "newalarm")
	router.OnState("newalarm", service.NewAlarm)
	router.OnState("newalarmselect", service.NewAlarmSelect)
	router.OnState("savealarm", service.SaveAlarm)

	router.OnState("editdelay", service.EditDelay)
	router.OnState("savedelay", service.SaveDelay)

	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Key)
	if err != nil {
		panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	go func() {
		for update := range updates {
			var tctx telegramconversation.TContext

			startTime := time.Now()
			metricsRegistry.TelegramInUpdatesTotal.Inc()

			if update.CallbackQuery != nil {
				log.Tracef("[%d][%s] Callback Query %s", update.CallbackQuery.Message.MessageID, update.CallbackQuery.From.UserName, update.CallbackQuery.Data)
				metricsRegistry.TelegramInQueriesTotal.Inc()

				chatID := strconv.FormatInt(update.CallbackQuery.Message.Chat.ID, 10)

				tctx = telegramconversation.NewTContext(chatID)
				tctx.SetMessageID(update.CallbackQuery.Message.MessageID)
				tctx.SetButtonData(update.CallbackQuery.Data)
			} else if update.Message != nil {
				log.Tracef("[%d][%s] Message %s", update.Message.MessageID, update.Message.From.UserName, update.Message.Text)
				chatID := strconv.FormatInt(update.Message.Chat.ID, 10)
				tctx = telegramconversation.NewTContext(chatID)
				tctx.SetMessageID(update.Message.MessageID)

				command := update.Message.Command()
				if command != "" {
					log.Tracef("[%d][%s] Command %s", tctx.MessageID(), update.Message.From.UserName, command)
					metricsRegistry.TelegramInCommandsTotal.Inc()
					tctx.SetCommand(command, update.Message.CommandArguments())
				} else {
					metricsRegistry.TelegramInMessagesTotal.Inc()
				}

				tctx.SetMessage(update.Message.Text)

				if state, payload, err := repo.GetState(ctx, tctx.ChatID(), "telegram"); err == nil {
					log.Tracef("[%d] Load State %s (Payload: %s)", tctx.MessageID(), state, payload)
					tctx.SetStatePayload(payload)
					tctx.ChangeState(state)
				}
			} else {
				continue
			}

			msgLog := tctx.LogFields(log)

			msgLog.Tracef("Route")
			tctx = router.Route(tctx)

			if update.CallbackQuery != nil {
				bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
			}

			if tctx.IsReply() {
				msgLog.Tracef("Reply: %s", tctx.Reply())
				reply := tgbotapi.NewMessage(tctx.ChatID64(), tctx.Reply())
				reply.ParseMode = tgbotapi.ModeMarkdown

				if tctx.IsRemoveSuggestions() {
					reply.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				} else if tctx.IsSuggestions() {
					buttons := []tgbotapi.KeyboardButton{}
					for _, button := range tctx.Suggestions() {
						buttons = append(buttons, tgbotapi.NewKeyboardButton(button.Label))
					}
					keyboard := tgbotapi.NewReplyKeyboard(buttons)
					reply.ReplyMarkup = keyboard
				}

				if tctx.IsKeyboard() {
					buttons := [][]tgbotapi.InlineKeyboardButton{}
					for _, button := range tctx.Keyboard() {
						buttons = append(
							buttons,
							tgbotapi.NewInlineKeyboardRow(
								tgbotapi.NewInlineKeyboardButtonData(button.Label, button.Data),
							),
						)
						//						buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(button.Label, button.Data))
					}
					keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)
					reply.ReplyMarkup = keyboard
				}

				metricsRegistry.TelegramOutMessagesTotal.Inc()
				bot.Send(reply)
			}

			if tctx.IsDeleteMessage() {
				msgLog.Tracef("Delete Message %d", tctx.DeletedMessageID())
				bot.DeleteMessage(tgbotapi.NewDeleteMessage(tctx.ChatID64(), tctx.DeletedMessageID()))
			}

			repo.UpdateState(ctx, tctx.ChatID(), "telegram", func(state string, payload string) (string, string, error) {
				msgLog.Tracef("Save State %s with payload `%s` (old: %s / `%s`)", tctx.State(), tctx.StatePayload(), state, payload)
				return tctx.State(), tctx.StatePayload(), nil
			})

			duration := time.Now().Sub(startTime).Seconds()
			metricsRegistry.TelegramRequestDuration.Observe(duration)
		}
	}()

	go func() {
		notificationChannel := cronService.NotificationChannel()
		for tctx := range notificationChannel {
			if tctx.IsReply() {
				log.Tracef("Reply: %s", tctx.Reply())

				reply := tgbotapi.NewMessage(tctx.ChatID64(), tctx.Reply())
				reply.ParseMode = tgbotapi.ModeMarkdown

				if tctx.IsRemoveSuggestions() {
					reply.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				} else if tctx.IsSuggestions() {
					buttons := []tgbotapi.KeyboardButton{}
					for _, button := range tctx.Suggestions() {
						buttons = append(buttons, tgbotapi.NewKeyboardButton(button.Label))
					}
					keyboard := tgbotapi.NewReplyKeyboard(buttons)
					reply.ReplyMarkup = keyboard
				}

				if tctx.IsKeyboard() {
					buttons := []tgbotapi.InlineKeyboardButton{}
					for _, button := range tctx.Keyboard() {
						buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(button.Label, button.Data))
					}
					keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons)
					reply.ReplyMarkup = keyboard
				}

				bot.Send(reply)
			}
		}
	}()

	cronService.Start(ctx)

	// metric endpoint
	if cfg.EnableMetrics {
		log.Info("start metrics endpoint...")
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Error(err)
			return
		}
	}

	<-ctx.Done()
}
