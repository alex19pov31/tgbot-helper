package tgbothelper

import (
	"strings"
	"time"

	tdlib "github.com/Arman92/go-tdlib"
)

var lastCommandChain = &CommandChain{}

// СhainElement - элемент цепочки команд
type СhainElement struct {
	command       string
	containText   string
	containButton string
	countSend     int
	countSkeep    int
	timeSend      time.Time
}

func (ce *СhainElement) isValid(text string, buttons ButtonList) bool {
	if ce.countSend > 0 || ce.countSkeep > 0 {
		return false
	}

	if ce.containText != "" && !strings.Contains(text, ce.containText) {
		return false
	}

	if ce.containButton != "" && buttons.GetButtonByText(ce.containButton) == nil {
		return false
	}

	return true
}

func (ce *СhainElement) run(client *tdlib.Client, text string, message *MessageData) bool {
	if !ce.isValid(text, message.Buttons) {
		ce.countSkeep++
		return false
	}

	SendCommand(client, text, message.ChatID)
	ce.countSend++
	ce.timeSend = time.Now()

	return true
}

func (ce *СhainElement) forseRun(client *tdlib.Client, chatID int64) {
	SendCommand(client, ce.command, chatID)
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
func (ch *CommandChain) Run(client *tdlib.Client, text string, message *MessageData) bool {
	if ch.finished {
		return false
	}

	for _, command := range ch.commands {
		if command.run(client, text, message) {
			return true
		}
	}

	ch.finished = true
	return false
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

// NewCommandMessage - новый элемент цепочки команд (отправка команды)
func NewCommandMessage(message string) *СhainElement {
	return &СhainElement{command: message, containText: message}
}

/*func ForwardMessage(chatID int64, messgesID []int64) {
	go GetClient().ForwardMessages(chatID, BotID, messgesID, false, false, false)
}*/
