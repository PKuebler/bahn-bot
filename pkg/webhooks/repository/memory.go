package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/pkuebler/bahn-bot/pkg/webhooks/domain"
)

// MemoryDatabase to test without persistence
type MemoryDatabase struct {
	Webhooks map[string]*domain.Webhook
	Mutex    *sync.Mutex
}

// NewMemoryDatabase with mutex
func NewMemoryDatabase() *MemoryDatabase {
	return &MemoryDatabase{
		Webhooks: map[string]*domain.Webhook{},
		Mutex:    &sync.Mutex{},
	}
}

// CreateWebhook at database
func (m *MemoryDatabase) CreateWebhook(ctx context.Context, webhook *domain.Webhook) (*domain.Webhook, error) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	m.Webhooks[webhook.GetID()] = webhook
	return webhook, nil
}

// GetWebhook by id
func (m *MemoryDatabase) GetWebhook(ctx context.Context, id string) (*domain.Webhook, error) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	if _, ok := m.Webhooks[id]; !ok {
		return nil, nil
	}

	return m.Webhooks[id], nil
}

// GetWebhooksByIdentifyer returns a empty array if nothing found
func (m *MemoryDatabase) GetWebhooksByIdentifyer(ctx context.Context, identifyer string, plattform string) ([]*domain.Webhook, error) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	webhooks := []*domain.Webhook{}

	for _, hook := range m.Webhooks {
		if identifyer == hook.GetIdentifyer() && plattform == hook.GetPlattform() {
			webhooks = append(webhooks, hook)
		}
	}

	return webhooks, nil
}

// GetWebhookByURLHash returns nothing if not found
func (m *MemoryDatabase) GetWebhookByURLHash(ctx context.Context, urlHash string) (*domain.Webhook, error) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	for _, hook := range m.Webhooks {
		if urlHash == hook.GetURLHash() {
			return hook, nil
		}
	}

	return nil, nil
}

// DeleteWebhook returns a error if webhook not found
func (m *MemoryDatabase) DeleteWebhook(ctx context.Context, webhookID string) error {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	if _, ok := m.Webhooks[webhookID]; !ok {
		return errors.New("not found")
	}

	delete(m.Webhooks, webhookID)

	return nil
}
