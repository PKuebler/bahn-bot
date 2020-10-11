package webhook

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/pkuebler/bahn-bot/pkg/config"
	trainalarmDomain "github.com/pkuebler/bahn-bot/pkg/trainalarms/domain"
	"github.com/pkuebler/bahn-bot/pkg/webhooks/application"
	webhookDomain "github.com/pkuebler/bahn-bot/pkg/webhooks/domain"
)

// Webhook endpoint
type Webhook struct {
	application    *application.Application
	webhookRepo    webhookDomain.Repository
	trainalarmRepo trainalarmDomain.Repository
	log            *logrus.Entry
	config         *config.Config
}

// NewWebhook endpoint to recive updates
func NewWebhook(
	application *application.Application,
	webhookRepo webhookDomain.Repository,
	trainalarmRepo trainalarmDomain.Repository,
	log *logrus.Entry,
	cfg *config.Config
) *Webhook {
	return &Webhook{
		application:    application,
		webhookRepo:    webhookRepo,
		trainalarmRepo: trainalarmRepo,
		log:            log,
		config:         cfg,
	}
}

// Start listener
func (w *Webhook) Start(ctx context.Context) {
	w.log.Info("start webhook listener")

	server := &http.Server{Addr: w.config.Webhook.Port}

	u, _ := url.Parse(w.config.Webhook.Endpoint)

	w.log.Info("listen %d/%s", w.config.Webhook.Port, u.Path)
	http.HandleFunc(u.Path, func(resp http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		defer func() {
			r := recover()
			if r != nil {
				var err error
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("Unknown error")
				}
				fmt.Println(err.Error())
				http.Error(resp, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		parts := strings.Split(r.URL.Path, "/")

		// find hash at database
		hook, err := w.webhookRepo.GetWebhookByURLHash(r.Context(), parts[len(parts)-1])
		if err != nil {
			fmt.Println(err.Error())
			http.Error(resp, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		if hook == nil {
			http.Error(resp, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}

		// json unmarschal body
		msg := TravelynxWebhook{}
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			fmt.Println(err.Error())
			http.Error(resp, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		// get current alarms
		alarms, err := w.trainalarmRepo.GetTrainAlarms(ctx, hook.GetIdentifyer(), hook.GetPlattform())
		if err != nil {
			fmt.Println(err.Error())
			http.Error(resp, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		// train.Train.Type, train.Train.Number
		// cancel webhook alarm or search train and cancel?

		// create cmd

		// application.AddAlarm
		// application.RemoveAlarm
	})

	go server.ListenAndServe()

	<-ctx.Done()

	if err := server.Shutdown(context.Background()); err != nil {
		panic(err)
	}
}
