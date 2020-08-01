package telegramconversation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConversationRouter(t *testing.T) {
	router := NewConversationRouter("start")
	assert.NotNil(t, router)
	assert.Equal(t, "start", router.DefaultState)
}

func TestOnCommand(t *testing.T) {
	router := NewConversationRouter("start")
	assert.NotNil(t, router)

	router.OnCommand("help", "start")
	assert.Equal(t, "start", router.EntryCommands["help"])
}

func TestOnState(t *testing.T) {
	router := NewConversationRouter("start")
	assert.NotNil(t, router)

	fn := func(ctx TContext) TContext {
		ctx.ChangeState("test")
		return ctx
	}
	router.OnState("help", fn)
	result := router.States["help"](NewTContext("123123"))
	assert.Equal(t, "test", result.State())
}

func TestOnUnknownState(t *testing.T) {
	router := NewConversationRouter("start")
	assert.NotNil(t, router)

	fn := func(ctx TContext) TContext {
		ctx.ChangeState("test")
		return ctx
	}
	router.OnUnknownState(fn)
	result := router.UnknownStateAction(NewTContext("123123"))
	assert.Equal(t, "test", result.State())
}

func TestRoute(t *testing.T) {
	defaultState := "start"
	router := NewConversationRouter(defaultState)
	assert.NotNil(t, router)

	isUnknownState := false
	isState := false
	fnUnknownState := func(ctx TContext) TContext {
		isUnknownState = true
		return ctx
	}
	fnState := func(ctx TContext) TContext {
		isState = true
		return ctx
	}
	router.OnUnknownState(fnUnknownState)
	router.OnCommand("help", "cancel")
	router.OnState("test", fnState)

	// default state set?
	// unknown state function
	ctx := NewTContext("12341234")
	ctx = router.Route(ctx)
	assert.Equal(t, defaultState, ctx.State())
	assert.True(t, isUnknownState)
	assert.False(t, isState)

	// button router
	isUnknownState = false
	isState = false
	ctx = NewTContext("12341234")
	ctx.SetButtonData("test|abc")
	ctx = router.Route(ctx)
	assert.Equal(t, "test", ctx.State())
	assert.False(t, isUnknownState)
	assert.True(t, isState)

	// route command to state?
	isUnknownState = false
	isState = false
	ctx = NewTContext("12341234")
	ctx.SetCommand("help", "")
	ctx = router.Route(ctx)
	assert.Equal(t, "cancel", ctx.State())
	assert.True(t, isUnknownState)
	assert.False(t, isState)

	// call state action
	isUnknownState = false
	isState = false
	ctx = NewTContext("12341234")
	ctx.ChangeState("test")
	ctx = router.Route(ctx)
	assert.Equal(t, "test", ctx.State())
	assert.False(t, isUnknownState)
	assert.True(t, isState)
}
