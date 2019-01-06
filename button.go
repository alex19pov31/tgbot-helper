package tgbothelper

import (
	"fmt"
)

const ShowKeyboardButtonType = "ShowKeyboardButton"
const InlineButtonType = "InlineButton"
const CallbackButtonType = "CallbackButton"

type Button interface {
	GetData() string
	GetType() string
	Click()
}

type baseButton struct {
	text       string
	data       string
	typeButton string
}

func (b *baseButton) Init(text, data, typeButton string) {
	b.data = data
	b.text = text
	b.typeButton = typeButton
}

func (b *baseButton) GetData() string {
	return b.data
}

func (b *baseButton) GetType() string {
	return b.typeButton
}

func (b *baseButton) Click() {
	fmt.Println("Clicked on " + b.GetData() + " button")
}

type ShowKeyboardButton struct {
	baseButton
}

type InlineButton struct {
	baseButton
}

type CallbackButton struct {
	baseButton
}

func NewShowKeyboardButton(text string) *ShowKeyboardButton {
	button := &ShowKeyboardButton{}
	button.Init(text, "", ShowKeyboardButtonType)

	return button
}

func NewInlineButton(text, data string) *InlineButton {
	button := &InlineButton{}
	button.Init(text, data, InlineButtonType)

	return button
}

func NewCallbackButton(text, data string) *CallbackButton {
	button := &CallbackButton{}
	button.Init(text, data, CallbackButtonType)

	return button
}
