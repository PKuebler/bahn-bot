package repository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/pkuebler/bahn-bot/pkg/webhooks/domain"
)

func TestNewMemoryDatabase(t *testing.T) {
	db := NewMemoryDatabase()
	assert.NotNil(t, db)
	assert.NotNil(t, db.Mutex)
}

func TestCreateWebhook(t *testing.T) {
	db := NewMemoryDatabase()
	ctx := context.Background()

	query := createTestWebhook()

	hook, err := db.CreateWebhook(ctx, query)
	assert.Nil(t, err)
	assert.NotNil(t, hook)
	_, ok := db.Webhooks[hook.GetID()]
	assert.True(t, ok)
	assert.Equal(t, hook.GetID(), db.Webhooks[hook.GetID()].GetID())
}

func TestGetWebhook(t *testing.T) {
	db := NewMemoryDatabase()
	ctx := context.Background()

	query := createTestWebhook()

	// not found
	alarm, err := db.GetWebhook(ctx, query.GetID())
	assert.Nil(t, err)
	assert.Nil(t, alarm)

	// founded
	db.Webhooks[query.GetID()] = query
	alarm, err = db.GetWebhook(ctx, query.GetID())
	assert.Nil(t, err)
	assert.NotNil(t, alarm)
}

func TestGetWebhooksByIdentifyer(t *testing.T) {
	db := NewMemoryDatabase()
	ctx := context.Background()

	// empty database
	webhooks, err := db.GetWebhooksByIdentifyer(ctx, "1234", "telegram")
	assert.Nil(t, err)
	assert.Len(t, webhooks, 0)

	for i := 0; i < 4; i++ {
		hook, err := domain.NewWebhook("1234", "telegram", "hook", "123456789", domain.TravelynxProtocol)
		assert.Nil(t, err)
		db.Webhooks[hook.GetID()] = hook
	}

	webhooks, err = db.GetWebhooksByIdentifyer(ctx, "1234", "telegram")
	assert.Nil(t, err)
	assert.Len(t, webhooks, 4)
}

func TestGetWebhookByURLHash(t *testing.T) {
	db := NewMemoryDatabase()
	ctx := context.Background()

	query := createTestWebhook()

	// not found
	hook, err := db.GetWebhookByURLHash(ctx, query.GetURLHash())
	assert.Nil(t, err)
	assert.Nil(t, hook)

	// founded
	db.Webhooks[query.GetID()] = query
	hook, err = db.GetWebhookByURLHash(ctx, query.GetURLHash())
	assert.Nil(t, err)
	assert.NotNil(t, hook)
}
func TestDeleteWebhook(t *testing.T) {
	db := NewMemoryDatabase()
	ctx := context.Background()

	query := createTestWebhook()
	db.Webhooks[query.GetID()] = query

	// not found
	err := db.DeleteWebhook(ctx, uuid.New().String())
	assert.NotNil(t, err)
	assert.Len(t, db.Webhooks, 1)

	// found
	err = db.DeleteWebhook(ctx, query.GetID())
	assert.Nil(t, err)
	assert.Len(t, db.Webhooks, 0)
}

func createTestWebhook() *domain.Webhook {
	webhook, _ := domain.NewWebhook(
		"identifyer",
		"telegram",
		uuid.New().String(),
		"1234",
		domain.TravelynxProtocol,
	)

	return webhook
}
