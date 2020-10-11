package domain

import (
	"errors"

	"github.com/google/uuid"
	"github.com/lithammer/shortuuid/v3"
)

// WebhookProtocol to support multiple services
type WebhookProtocol string

const (
	// TravelynxProtocol to support travelynx webhooks
	TravelynxProtocol WebhookProtocol = "travelynx"
)

// NewWebhookProtocol converts string to WebhookProtocol
func NewWebhookProtocol(protocol string) (WebhookProtocol, error) {
	switch protocol {
	case "travelynx":
		return TravelynxProtocol, nil
	}

	return TravelynxProtocol, errors.New("unknown protocol")
}

// Webhook represent a webhook to use it e.g. travelynx
type Webhook struct {
	id         string
	identifyer string
	plattform  string

	urlHash  string
	token    string
	protocol WebhookProtocol
}

// NewWebhook returns a new webhook
func NewWebhook(
	identifyer string,
	plattform string,
	token string,
	protocol WebhookProtocol,
) (*Webhook, error) {
	return &Webhook{
		id:         uuid.New().String(),
		identifyer: identifyer,
		plattform:  plattform,
		urlHash:    shortuuid.New(),
		token:      token,
		protocol:   protocol,
	}, nil
}

// NewWebhookFromRepository returns a new webhook with all protected fields
func NewWebhookFromRepository(
	id string,
	identifyer string,
	plattform string,
	urlHash string,
	token string,
	protocol WebhookProtocol,
) (*Webhook, error) {
	return &Webhook{
		id:         id,
		identifyer: identifyer,
		plattform:  plattform,
		urlHash:    urlHash,
		token:      token,
		protocol:   protocol,
	}, nil
}

// GetID from domain
func (w *Webhook) GetID() string {
	return w.id
}

// GetIdentifyer from webhook
func (w *Webhook) GetIdentifyer() string {
	return w.identifyer
}

// GetPlattform from webhook
func (w *Webhook) GetPlattform() string {
	return w.plattform
}

// GetURLHash from webhook
func (w *Webhook) GetURLHash() string {
	return w.urlHash
}

// GetToken from webhook
func (w *Webhook) GetToken() string {
	return w.token
}

// GetProtocol from webhook
func (w *Webhook) GetProtocol() WebhookProtocol {
	return w.protocol
}
