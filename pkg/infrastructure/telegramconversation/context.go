package telegramconversation

import (
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

// TButton for keyboard
type TButton struct {
	Label string
	Data  string
}

// NewTButton creates a telegram button
func NewTButton(label string, data string) TButton {
	return TButton{label, data}
}

// TContext for telegram requests
type TContext struct {
	state        string
	statePayload string

	message   string
	messageID int
	chatID    string

	isButtonPressed bool
	buttonData      string
	command         string
	commandQuery    string

	reply             *string
	removeSuggestions bool
	buttons           *[]TButton
	suggestionButtons *[]TButton
	deleteMessageID   *int
	columns           int
}

// NewTContext with chat context
func NewTContext(chatID string) TContext {
	return TContext{
		chatID:            chatID,
		removeSuggestions: false,
		columns:           2,
		isButtonPressed:   false,
	}
}

// State return the current state
func (c *TContext) State() string {
	return c.state
}

// IsState compare the state
func (c *TContext) IsState(state string) bool {
	return c.state == state
}

// SetStatePayload to save payload on state
func (c *TContext) SetStatePayload(payload string) {
	c.statePayload = payload
}

// StatePayload to get payload
func (c *TContext) StatePayload() string {
	return c.statePayload
}

// ChangeState of the conversation
func (c *TContext) ChangeState(state string) {
	c.state = state
}

// IsButtonPressed on telegram keyboard
func (c *TContext) IsButtonPressed() bool {
	return c.isButtonPressed
}

// ButtonData from telegram button. empty string if not pressed.
func (c *TContext) ButtonData() string {
	return c.buttonData
}

// SetButtonData to context
func (c *TContext) SetButtonData(data string) {
	parts := strings.SplitN(data, "|", 2)
	if len(parts) == 0 {
		return
	}

	c.state = parts[0]

	if len(parts) == 2 {
		c.buttonData = parts[1]
	}

	c.isButtonPressed = true
}

// IsCommand compare current command with string. if no command is the result false.
func (c *TContext) IsCommand(cmd string) bool {
	return c.command == cmd
}

// Command returns the command. empty string if no command
func (c *TContext) Command() string {
	return c.command
}

// SetCommand to context
func (c *TContext) SetCommand(cmd string, query string) {
	c.command = cmd
	c.commandQuery = query
}

// CommandQuery returns the query. empty string if no command or query
func (c *TContext) CommandQuery() string {
	return c.commandQuery
}

// MessageID returns the current message id
func (c *TContext) MessageID() int {
	return c.messageID
}

// SetMessageID to context
func (c *TContext) SetMessageID(id int) {
	c.messageID = id
}

// Message returns the current message
func (c *TContext) Message() string {
	return c.message
}

// SetMessage to context
func (c *TContext) SetMessage(msg string) {
	c.message = msg
}

// ChatID returns the current chat id
func (c *TContext) ChatID() string {
	return c.chatID
}

// ChatID64 return the current chat id as int64
func (c *TContext) ChatID64() int64 {
	n, _ := strconv.ParseInt(c.chatID, 10, 64)
	return n
}

// Send a message
func (c *TContext) Send(txt string) TContext {
	return c.SendAndRemoveSuggestions(txt)
}

// SendWithState a message and set the state
func (c *TContext) SendWithState(txt string, state string) TContext {
	c.ChangeState(state)
	return c.SendAndRemoveSuggestions(txt)
}

// SendWithKeyboard send a message and show buttons
func (c *TContext) SendWithKeyboard(txt string, buttons []TButton, columns int) TContext {
	c.reply = &txt
	c.buttons = &buttons
	c.columns = columns
	return *c
}

// IsKeyboard set
func (c *TContext) IsKeyboard() bool {
	return c.buttons != nil
}

// Keyboard buttons
func (c *TContext) Keyboard() []TButton {
	if c.buttons == nil {
		return []TButton{}
	}

	return *c.buttons
}

// SendWithSuggestions send a message and show buttons by the textfield
func (c *TContext) SendWithSuggestions(txt string, buttons []TButton, columns int) TContext {
	c.reply = &txt
	c.suggestionButtons = &buttons
	c.columns = columns
	return *c
}

// SendAndRemoveSuggestions send a message and remove the keyboard if exists
func (c *TContext) SendAndRemoveSuggestions(txt string) TContext {
	c.reply = &txt
	c.removeSuggestions = true
	return *c
}

// IsSuggestions set
func (c *TContext) IsSuggestions() bool {
	return c.suggestionButtons != nil
}

// Suggestions buttons
func (c *TContext) Suggestions() []TButton {
	if c.suggestionButtons == nil {
		return []TButton{}
	}

	return *c.suggestionButtons
}

// IsRemoveSuggestions active
func (c *TContext) IsRemoveSuggestions() bool {
	return c.removeSuggestions
}

// IsReply message
func (c *TContext) IsReply() bool {
	return c.reply != nil
}

// Reply message string. returns empty string if no reply set
func (c *TContext) Reply() string {
	if !c.IsReply() {
		return ""
	}
	return *c.reply
}

// DeleteMessage from context chat
func (c *TContext) DeleteMessage(msgID int) {
	c.deleteMessageID = &msgID
}

// IsDeleteMessage set
func (c *TContext) IsDeleteMessage() bool {
	return c.deleteMessageID != nil
}

// DeletedMessageID to delete. returns -1 if no message set
func (c *TContext) DeletedMessageID() int {
	if !c.IsDeleteMessage() {
		return -1
	}

	return *c.deleteMessageID
}

// LogFields adds log fields to logrus
func (c *TContext) LogFields(log *logrus.Entry) *logrus.Entry {
	return log.WithFields(logrus.Fields{
		"msgID":         c.MessageID(),
		"chatID":        c.ChatID(),
		"cmd":           c.Command(),
		"cmdQuery":      c.CommandQuery(),
		"state":         c.State(),
		"statePayload":  c.StatePayload(),
		"buttonPressed": c.IsButtonPressed(),
		"button":        c.ButtonData(),
	})
}
