package telegramconversation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTButton(t *testing.T) {
	button := NewTButton("label", "data")
	assert.Equal(t, "label", button.Label)
	assert.Equal(t, "data", button.Data)
}

func TestNewTContext(t *testing.T) {
	ctx := NewTContext("234")
	assert.Equal(t, "234", ctx.chatID)
	assert.Equal(t, "234", ctx.ChatID())
	assert.Equal(t, int64(234), ctx.ChatID64())

	ctx.ChangeState("mystate")
	assert.Equal(t, "mystate", ctx.state)

	assert.Equal(t, "mystate", ctx.State())
	assert.True(t, ctx.IsState("mystate"))
	assert.False(t, ctx.IsState("foo"))

	ctx.SetMessageID(123)
	assert.Equal(t, 123, ctx.MessageID())

	ctx.SetMessage("huhu")
	assert.Equal(t, "huhu", ctx.Message())

	assert.Equal(t, "", ctx.StatePayload())
	ctx.SetStatePayload("abc")
	assert.Equal(t, "abc", ctx.StatePayload())
}

func TestContextButton(t *testing.T) {
	ctxWithNoButton := TContext{
		chatID:            "234",
		removeSuggestions: false,
		columns:           2,
	}

	assert.False(t, ctxWithNoButton.IsButtonPressed())
	assert.Equal(t, "", ctxWithNoButton.ButtonData())

	ctx := TContext{
		chatID:            "234",
		removeSuggestions: false,
		columns:           2,
	}
	ctx.SetButtonData("testState|testData")

	assert.True(t, ctx.IsButtonPressed())
	assert.Equal(t, "testData", ctx.ButtonData())
	assert.Equal(t, "testState", ctx.State())
}

func TestContextCommand(t *testing.T) {
	ctxWithNoCommand := TContext{
		chatID:            "234",
		removeSuggestions: false,
		columns:           2,
	}

	assert.False(t, ctxWithNoCommand.IsCommand("test"))
	assert.Equal(t, "", ctxWithNoCommand.CommandQuery())

	ctx := TContext{
		chatID:            "234",
		removeSuggestions: false,
		columns:           2,
	}

	ctx.SetCommand("test", "querydata")

	assert.True(t, ctx.IsCommand("test"))
	assert.Equal(t, "test", ctx.Command())
	assert.Equal(t, "querydata", ctx.CommandQuery())
}

func TestSend(t *testing.T) {
	ctx := TContext{
		chatID:            "234",
		removeSuggestions: false,
		columns:           2,
	}

	assert.False(t, ctx.IsReply())
	assert.Equal(t, "", ctx.Reply())

	ctx.Send("test message")
	assert.NotNil(t, ctx.reply)
	assert.True(t, ctx.removeSuggestions)
	assert.True(t, ctx.IsReply())
	assert.Equal(t, "test message", ctx.Reply())
}

func TestSendWithState(t *testing.T) {
	ctx := TContext{
		chatID:            "234",
		removeSuggestions: false,
		columns:           2,
	}

	assert.False(t, ctx.IsReply())
	assert.Equal(t, "", ctx.Reply())

	ctx.SendWithState("test message", "test")
	assert.NotNil(t, ctx.reply)
	assert.Equal(t, "test", ctx.state)
	assert.True(t, ctx.removeSuggestions)
	assert.True(t, ctx.IsReply())
	assert.Equal(t, "test message", ctx.Reply())
}

func TestSendWithKeyboard(t *testing.T) {
	ctx := TContext{
		chatID:            "234",
		removeSuggestions: false,
		columns:           2,
	}

	assert.False(t, ctx.IsKeyboard())
	assert.Len(t, ctx.Keyboard(), 0)
	assert.False(t, ctx.IsReply())
	assert.Equal(t, "", ctx.Reply())

	buttons := []TButton{
		NewTButton("label", "data"),
	}

	ctx.SendWithKeyboard("test message", buttons, 4)
	assert.NotNil(t, ctx.reply)
	assert.NotNil(t, ctx.buttons)
	assert.Equal(t, 4, ctx.columns)

	assert.True(t, ctx.IsKeyboard())
	assert.Len(t, ctx.Keyboard(), 1)
	assert.True(t, ctx.IsReply())
	assert.Equal(t, "test message", ctx.Reply())
}

func TestSendWithSuggestions(t *testing.T) {
	ctx := TContext{
		chatID:            "234",
		removeSuggestions: false,
		columns:           2,
	}

	assert.False(t, ctx.IsSuggestions())
	assert.Len(t, ctx.Suggestions(), 0)
	assert.False(t, ctx.IsReply())
	assert.Equal(t, "", ctx.Reply())

	buttons := []TButton{
		NewTButton("label", "data"),
	}

	ctx.SendWithSuggestions("test message", buttons, 4)
	assert.NotNil(t, ctx.reply)
	assert.NotNil(t, ctx.suggestionButtons)
	assert.Equal(t, 4, ctx.columns)

	assert.True(t, ctx.IsSuggestions())
	assert.Len(t, ctx.Suggestions(), 1)
	assert.True(t, ctx.IsReply())
	assert.Equal(t, "test message", ctx.Reply())
}

func TestSendAndRemoveSuggestions(t *testing.T) {
	ctx := TContext{
		chatID:            "234",
		removeSuggestions: false,
		columns:           2,
	}

	assert.False(t, ctx.IsRemoveSuggestions())
	assert.False(t, ctx.IsReply())
	assert.Equal(t, "", ctx.Reply())

	ctx.SendAndRemoveSuggestions("test message")
	assert.NotNil(t, ctx.reply)
	assert.Equal(t, "test message", *ctx.reply)
	assert.True(t, ctx.removeSuggestions)

	assert.True(t, ctx.IsRemoveSuggestions())
	assert.True(t, ctx.IsReply())
	assert.Equal(t, "test message", ctx.Reply())
}

func TestDeleteMessage(t *testing.T) {
	ctx := TContext{
		chatID:            "234",
		removeSuggestions: false,
		columns:           2,
	}

	assert.False(t, ctx.IsDeleteMessage())

	ctx.DeleteMessage(234234)
	assert.NotNil(t, ctx.deleteMessageID)
	assert.Equal(t, 234234, ctx.DeletedMessageID())

	assert.True(t, ctx.IsDeleteMessage())
}
