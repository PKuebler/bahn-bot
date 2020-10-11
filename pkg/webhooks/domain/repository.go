package domain

import (
	"context"
)

// Repository interface to datahandling between repository and application
type Repository interface {
	CreateWebhook(ctx context.Context, webhook *Webhook) (*Webhook, error)
	GetWebhook(ctx context.Context, id string) (*Webhook, error)
	GetWebhooksByIdentifyer(ctx context.Context, identifyer string, plattform string) ([]*Webhook, error)
	GetWebhookByURLHash(ctx context.Context, urlHash string) (*Webhook, error)
	DeleteWebhook(ctx context.Context, webhookID string) error
}
