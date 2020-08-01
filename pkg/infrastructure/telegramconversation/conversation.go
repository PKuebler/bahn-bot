package telegramconversation

// Action represents the action to be executed in a given state.
type Action func(ctx TContext) TContext

// ConversationRouter route all conversation queries
type ConversationRouter struct {
	DefaultState       string
	EntryCommands      map[string]string // command / state
	States             map[string]Action
	UnknownStateAction Action
}

// NewConversationRouter return a new engine
func NewConversationRouter(defaultState string) *ConversationRouter {
	return &ConversationRouter{
		DefaultState:  defaultState,
		EntryCommands: map[string]string{},
		States:        map[string]Action{},
	}
}

// OnCommand set a command to a state
func (c *ConversationRouter) OnCommand(command string, state string) {
	c.EntryCommands[command] = state
}

// OnState set a action to a state
func (c *ConversationRouter) OnState(state string, action Action) {
	c.States[state] = action
}

// OnUnknownState after migrations or failed configurations
func (c *ConversationRouter) OnUnknownState(action Action) {
	c.UnknownStateAction = action
}

// Route the context
func (c *ConversationRouter) Route(ctx TContext) TContext {
	// if no state, use default
	if ctx.State() == "" {
		ctx.ChangeState(c.DefaultState)
	}

	// change state if command defined
	cmd := ctx.Command()
	if cmd != "" {
		if _, ok := c.EntryCommands[cmd]; ok {
			ctx.ChangeState(c.EntryCommands[cmd])
		}
	}

	if _, ok := c.States[ctx.State()]; !ok {
		ctx = c.UnknownStateAction(ctx)
	} else {
		ctx = c.States[ctx.State()](ctx)
	}

	// send to telegram
	return ctx
}
