package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewWebhook(t *testing.T) {
	webhook, err := NewWebhook("identifyer", "telegram", "hook", "1233245345", "1234567890")
	assert.Nil(t, err)
	assert.NotNil(t, webhook)
}

func TestWebhookGetters(t *testing.T) {
	identifyer := uuid.New().String()
	plattform := "telegram"
	urlHash := uuid.New().String()
	token := "234234234"
	protocol := TravelynxProtocol

	webhook, err := NewWebhook(identifyer, plattform, urlHash, token, protocol)
	assert.Nil(t, err)
	assert.NotNil(t, webhook)

	assert.Equal(t, identifyer, webhook.GetIdentifyer())
	assert.Equal(t, plattform, webhook.GetPlattform())
	assert.Equal(t, urlHash, webhook.GetURLHash())
	assert.Equal(t, token, webhook.GetToken())
	assert.Equal(t, protocol, webhook.GetProtocol())
}
