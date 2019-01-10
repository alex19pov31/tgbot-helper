package tgbothelper

import (
	"github.com/Arman92/go-tdlib"
)

// ShowKeyboardButtonType - клавиатурные кнопки
const ShowKeyboardButtonType = "ShowKeyboardButton"

// InlineButtonType - инлайновые кнопки
const InlineButtonType = "InlineButton"

// CallbackButtonType - кнопки с обработкой callback запросов
const CallbackButtonType = "CallbackButton"

// Button - интерфейс кнопки
type Button interface {
	GetData() string
	GetText() string
	GetType() string
	GetChatID() int64
	GetMessageID() int64
	Click(client *tdlib.Client)
}

type baseButton struct {
	chatID     int64
	messageID  int64
	text       string
	data       string
	typeButton string
}

func (b *baseButton) Init(text, data, typeButton string, chatID, messageID int64) {
	b.data = data
	b.text = text
	b.chatID = chatID
	b.messageID = messageID
	b.typeButton = typeButton
}

func (b *baseButton) GetText() string {
	return b.text
}

func (b *baseButton) GetData() string {
	return b.data
}

func (b *baseButton) GetType() string {
	return b.typeButton
}

func (b *baseButton) GetChatID() int64 {
	return b.chatID
}

func (b *baseButton) GetMessageID() int64 {
	return b.messageID
}

func (b *baseButton) Click(client *tdlib.Client) {
	SendMessage(client, b.GetText(), b.GetChatID(), 0)
}

// ShowKeyboardButton - клавиатурная кнопка
type ShowKeyboardButton struct {
	baseButton
}

// InlineButton - инлайновая кнопка
type InlineButton struct {
	baseButton
}

// CallbackButton - кнопка с обработкой callback запроса
type CallbackButton struct {
	baseButton
}

// Click - клик по кнопке
func (cb *CallbackButton) Click(client *tdlib.Client) {
	client.GetCallbackQueryAnswer(
		cb.GetChatID(),
		cb.GetMessageID(),
		tdlib.NewCallbackQueryPayloadData([]byte(cb.GetData())),
	)
}

func newShowKeyboardButton(text string, chatID, messageID int64) *ShowKeyboardButton {
	button := &ShowKeyboardButton{}
	button.Init(text, "", ShowKeyboardButtonType, chatID, messageID)

	return button
}

func newInlineButton(text, data string, chatID, messageID int64) *InlineButton {
	button := &InlineButton{}
	button.Init(text, data, InlineButtonType, chatID, messageID)

	return button
}

func newCallbackButton(text, data string, chatID, messageID int64) *CallbackButton {
	button := &CallbackButton{}
	button.Init(text, data, CallbackButtonType, chatID, messageID)

	return button
}
