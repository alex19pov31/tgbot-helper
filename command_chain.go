package tgbothelper

import (
	"strings"
	"time"

	tdlib "github.com/Arman92/go-tdlib"
)

var lastCommandChain = &CommandChain{}

// СhainElement - элемент цепочки команд
type СhainElement struct {
	command           string
	containText       string
	containButton     string
	containButtonText string
	callbackData      string
	countSend         int
	countSkeep        int
	timeSend          time.Time
	button            *Button
}

func (ce *СhainElement) isCallback() bool {
	return ce.callbackData != ""
}

func (ce *СhainElement) isValid(text string, buttons ButtonList) bool {
	if ce.countSend > 0 || ce.countSkeep > 0 {
		return false
	}

	if ce.containText != "" && !strings.Contains(text, ce.containText) {
		return false
	}

	if ce.containButton != "" {
		ce.button = buttons.GetButtonByText(ce.containButton)
		if ce.button == nil {
			return false
		}
	}

	if ce.callbackData != "" {
		ce.button = buttons.GetButtonByText(ce.callbackData)
		if ce.button == nil {
			return false
		}
	}

	if ce.containButtonText != "" {
		ce.button = buttons.GetButtonByContainText(ce.containButtonText)
		if ce.button == nil {
			return false
		}
	}

	return true
}

func (ce *СhainElement) run(client *tdlib.Client, text string, message *MessageData) *Command {
	if !ce.isValid(text, message.Buttons) {
		ce.countSkeep++
		return nil
	}

	if ce.button != nil {
		return (*ce.button).Click(client)
	} else {
		return SendMessage(client, ce.command, message.ChatID, 0)
	}

	ce.countSend++
	ce.timeSend = time.Now()

	return nil
}

func (ce *СhainElement) forseRun(client *tdlib.Client, chatID int64) {
	SendMessage(client, ce.command, chatID, 0)
	ce.countSend++
	ce.timeSend = time.Now()
}

// CommandChain - цепочка команд
type CommandChain struct {
	id       string
	commands []*СhainElement
	finished bool
	created  time.Time
}

// Run - запуск цепочки команд
func (ch *CommandChain) Run(client *tdlib.Client, message *MessageData, initFunc routeCallback) *Command {
	text := message.Message
	if ch.finished {
		return nil
	}

	for _, command := range ch.commands {
		cmd := command.run(client, text, message)
		if cmd != nil {
			if command.button != nil && (*command.button).GetType() == CallbackButtonType {
				newMessage := GetMessageByID(client, message.ChatID, message.MessageID)
				initFunc(newMessage)
			}

			return cmd
		}
	}

	ch.finished = true
	return nil
}

// ForseRun - принудительный запус команд
func (ch *CommandChain) ForseRun(client *tdlib.Client, chatID int64) {
	if ch.finished {
		return
	}

	for _, command := range ch.commands {
		command.forseRun(client, chatID)
		return
	}
}

// GetCommandChain - текущая цепочка команд
func GetCommandChain() *CommandChain {
	return lastCommandChain
}

// NewCommandChain - новая цепочка команд
func NewCommandChain(id string, commands ...*СhainElement) *CommandChain {
	lastCommandChain = &CommandChain{id: id, commands: commands, created: time.Now()}

	return lastCommandChain
}

// NewCommandButton - новый элемент цепочки команд (нажатие на кнопку)
func NewCommandButton(button string) *СhainElement {
	return &СhainElement{command: button, containButton: button}
}

func NewContainTextButton(button string) *СhainElement {
	return &СhainElement{containButtonText: button}
}

// NewCallbackButton - новый элемент цепочки команд (нажатие на кнопку)
func NewCallbackButton(button string) *СhainElement {
	return &СhainElement{callbackData: button}
}

// NewCommandMessage - новый элемент цепочки команд (отправка команды)
func NewCommandMessage(message string) *СhainElement {
	return &СhainElement{command: message, containText: message}
}

/*func ForwardMessage(chatID int64, messgesID []int64) {
	go GetClient().ForwardMessages(chatID, BotID, messgesID, false, false, false)
}*/
