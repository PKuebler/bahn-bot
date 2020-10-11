package repository

import (
	"context"
	"os"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"

	"github.com/pkuebler/bahn-bot/pkg/webhooks/domain"
)

func TestNewSQLDatabase(t *testing.T) {
	db := NewSQLDatabase(os.Getenv("DB_DIALECT"), os.Getenv("DB_PATH"))
	assert.NotNil(t, db)

	db.Close()
}

func TestConvertWebhookModel(t *testing.T) {
	hook, err := domain.NewWebhook("mychatID", "telegram", "1234567", "123456", domain.TravelynxProtocol)
	assert.NotNil(t, hook)
	assert.Nil(t, err)

	// convert
	model := NewSQLWebhookModel(hook)
	assert.NotNil(t, model)

	assert.Equal(t, hook.GetID(), model.ID)
	assert.Equal(t, hook.GetIdentifyer(), model.Identifyer)
	assert.Equal(t, hook.GetPlattform(), model.Plattform)
	assert.Equal(t, hook.GetURLHash(), model.Path)
	assert.Equal(t, hook.GetToken(), model.Token)
	assert.Equal(t, hook.GetProtocol(), model.Protocol)

	// convert back
	back := model.Webhook()
	assert.NotNil(t, back)

	assert.Equal(t, model.ID, back.GetID())
	assert.Equal(t, model.Identifyer, back.GetIdentifyer())
	assert.Equal(t, model.Plattform, back.GetPlattform())
	assert.Equal(t, model.Path, back.GetURLHash())
	assert.Equal(t, model.Token, back.GetToken())
	assert.Equal(t, model.Protocol, back.GetProtocol())
}

func TestSQLCreateWebhook(t *testing.T) {
	ctx := context.Background()
	db := NewSQLDatabase(os.Getenv("DB_DIALECT"), os.Getenv("DB_PATH"))
	assert.NotNil(t, db)
	defer db.Close()
	err := db.db.Delete(SQLWebhookModel{}).Error
	assert.Nil(t, err)

	query := createTestWebhook()

	hook, err := db.CreateWebhook(ctx, query)
	assert.Nil(t, err)
	assert.NotNil(t, hook)
	var dbHook SQLWebhookModel
	err = db.db.Where("id = ?", hook.GetID()).First(&dbHook).Error
	assert.Nil(t, err)
	assert.Equal(t, hook.GetID(), dbHook.ID)
}

func TestSQLGetWebhook(t *testing.T) {
	ctx := context.Background()
	db := NewSQLDatabase(os.Getenv("DB_DIALECT"), os.Getenv("DB_PATH"))
	assert.NotNil(t, db)
	defer db.Close()
	err := db.db.Delete(SQLWebhookModel{}).Error
	assert.Nil(t, err)

	query := createTestWebhook()

	// not found
	hook, err := db.GetWebhook(ctx, query.GetID())
	assert.Nil(t, err)
	assert.Nil(t, hook)

	// founded
	err = db.db.Create(NewSQLWebhookModel(query)).Error
	assert.Nil(t, err)

	hook, err = db.GetWebhook(ctx, query.GetID())
	assert.Nil(t, err)
	assert.NotNil(t, hook)
}

func TestSQLGetWebhooksByIdentifyer(t *testing.T) {
	ctx := context.Background()
	db := NewSQLDatabase(os.Getenv("DB_DIALECT"), os.Getenv("DB_PATH"))
	assert.NotNil(t, db)
	defer db.Close()
	err := db.db.Delete(SQLWebhookModel{}).Error
	assert.Nil(t, err)

	// empty database
	hooks, err := db.GetWebhooksByIdentifyer(ctx, "1234", "telegram")
	assert.Nil(t, err)
	assert.Len(t, hooks, 0)

	for i := 0; i < 4; i++ {
		hook, err := domain.NewWebhook("1234", "telegram", "hook", "123456789", domain.TravelynxProtocol)
		assert.Nil(t, err)
		err = db.db.Create(NewSQLWebhookModel(hook)).Error
		assert.Nil(t, err)
	}

	hooks, err = db.GetWebhooksByIdentifyer(ctx, "1234", "telegram")
	assert.Nil(t, err)
	assert.Len(t, hooks, 4)
}

func TestSQLGetWebhookByURLHash(t *testing.T) {
	ctx := context.Background()
	db := NewSQLDatabase(os.Getenv("DB_DIALECT"), os.Getenv("DB_PATH"))
	assert.NotNil(t, db)
	defer db.Close()
	err := db.db.Delete(SQLWebhookModel{}).Error
	assert.Nil(t, err)

	query := createTestWebhook()

	// not found
	hook, err := db.GetWebhookByURLHash(ctx, query.GetURLHash())
	assert.Nil(t, err)
	assert.Nil(t, hook)

	// founded
	err = db.db.Create(NewSQLWebhookModel(query)).Error
	assert.Nil(t, err)

	hook, err = db.GetWebhookByURLHash(ctx, query.GetURLHash())
	assert.Nil(t, err)
	assert.NotNil(t, hook)
}

func TestSQLDeleteWebhook(t *testing.T) {
	ctx := context.Background()
	db := NewSQLDatabase(os.Getenv("DB_DIALECT"), os.Getenv("DB_PATH"))
	assert.NotNil(t, db)
	defer db.Close()
	err := db.db.Delete(SQLWebhookModel{}).Error
	assert.Nil(t, err)

	hook := createTestWebhook()

	// not found
	err = db.DeleteWebhook(ctx, hook.GetID())
	assert.Nil(t, err)

	// delete
	err = db.db.Create(NewSQLWebhookModel(hook)).Error
	assert.Nil(t, err)

	err = db.DeleteWebhook(ctx, hook.GetID())
	assert.Nil(t, err)

	var dbHook SQLWebhookModel
	err = db.db.Where("id = ?", hook.GetID()).First(&dbHook).Error
	assert.NotNil(t, err)
}
