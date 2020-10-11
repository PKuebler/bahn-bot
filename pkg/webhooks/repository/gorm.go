package repository

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/pkuebler/bahn-bot/pkg/webhooks/domain"
)

// SQLDatabase to persistence
type SQLDatabase struct {
	db *gorm.DB
}

// NewSQLDatabase by gorm
func NewSQLDatabase(dialect string, path string) *SQLDatabase {
	db, err := gorm.Open(dialect, path)
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&SQLWebhookModel{})

	return &SQLDatabase{
		db: db,
	}
}

// SQLWebhookModel to save webhooks
type SQLWebhookModel struct {
	ID         string `gorm:"primary_key"`
	Identifyer string
	Plattform  string
	URLHash    string
	Token      string
	Protocol   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// NewSQLWebhookModel converts domain to model
func NewSQLWebhookModel(hook *domain.Webhook) *SQLWebhookModel {
	return &SQLWebhookModel{
		ID:         hook.GetID(),
		Identifyer: hook.GetIdentifyer(),
		Plattform:  hook.GetPlattform(),
		URLHash:    hook.GetURLHash(),
		Token:      hook.GetToken(),
		Protocol:   string(hook.GetProtocol()),
	}
}

// Webhook convert to Webhook domain model
func (s *SQLWebhookModel) Webhook() *domain.Webhook {
	protocol, _ := domain.NewWebhookProtocol(s.Protocol)

	webhook, _ := domain.NewWebhookFromRepository(
		s.ID,
		s.Identifyer,
		s.Plattform,
		s.URLHash,
		s.Token,
		protocol,
	)

	return webhook
}

// Close database
func (s *SQLDatabase) Close() {
	s.db.Close()
}

// CreateWebhook at database
func (s *SQLDatabase) CreateWebhook(ctx context.Context, hook *domain.Webhook) (*domain.Webhook, error) {
	model := NewSQLWebhookModel(hook)
	if err := s.db.Create(model).Error; err != nil {
		return nil, err
	}
	return model.Webhook(), nil
}

// GetWebhook by id
func (s *SQLDatabase) GetWebhook(ctx context.Context, id string) (*domain.Webhook, error) {
	var hook SQLWebhookModel
	if res := s.db.Where("id = ?", id).Take(&hook); res.Error != nil {
		if res.RecordNotFound() {
			return nil, nil
		}
		return nil, res.Error
	}

	return hook.Webhook(), nil
}

// GetWebhooksByIdentifyer returns a empty array if nothing found
func (s *SQLDatabase) GetWebhooksByIdentifyer(ctx context.Context, identifyer string, plattform string) ([]*domain.Webhook, error) {
	var results []SQLWebhookModel
	if res := s.db.Where("identifyer = ? AND plattform = ?", identifyer, plattform).Find(&results); res.Error != nil {
		if res.RecordNotFound() {
			return nil, nil
		}
		return nil, res.Error
	}

	// convert
	hooks := []*domain.Webhook{}
	for _, result := range results {
		hooks = append(hooks, result.Webhook())
	}

	return hooks, nil
}

// GetWebhookByURLHash returns nothing if not found
func (s *SQLDatabase) GetWebhookByURLHash(ctx context.Context, urlHash string) (*domain.Webhook, error) {
	var hook SQLWebhookModel
	if res := s.db.Where("urlHash = ?", urlHash).Take(&hook); res.Error != nil {
		if res.RecordNotFound() {
			return nil, nil
		}
		return nil, res.Error
	}

	return hook.Webhook(), nil
}

// DeleteWebhook returns a error if webhook not found
func (s *SQLDatabase) DeleteWebhook(ctx context.Context, webhookID string) error {
	err := s.db.Where("id = ?", webhookID).Delete(SQLWebhookModel{}).Error
	return err
}
