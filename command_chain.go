package tgbothelper

import (
	"strings"
	"time"

	tdlib "github.com/Arman92/go-tdlib"
)

var lastCommandChain = &CommandChain{}

type chainElement struct {
	command       string
	containText   string
	containButton string
	countSend     int
	countSkeep    int
	timeSend      time.Time
}

func (ce *chainElement) isValid(text string, buttons ButtonList) bool {
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

func (ce *chainElement) run(client *tdlib.Client, text string, message *MessageData) bool {
	if !ce.isValid(text, message.Buttons) {
		ce.countSkeep++
		return false
	}

	SendCommand(client, text, message.ChatID)
	ce.countSend++
	ce.timeSend = time.Now()

	return true
}

func (ce *chainElement) forseRun(client *tdlib.Client, chatID int64) {
	SendCommand(client, ce.command, chatID)
	ce.countSend++
	ce.timeSend = time.Now()
}

type CommandChain struct {
	id       string
	commands []*chainElement
	finished bool
	created  time.Time
}

func (ch *CommandChain) Run(client *tdlib.Client, text string, message *MessageData) bool {
	if ch.finished {
		return false
	}

	for _, command := range ch.commands {
		if command.run(client, text, message) {
			return true
		} else {
		}
	}

	ch.finished = true
	return false
}

func (ch *CommandChain) ForseRun(client *tdlib.Client, chatID int64) {
	if ch.finished {
		return
	}

	for _, command := range ch.commands {
		command.forseRun(client, chatID)
		return
	}
}

func GetCommandChain() *CommandChain {
	return lastCommandChain
}

func NewCommandChain(id string, commands ...*chainElement) *CommandChain {
	lastCommandChain = &CommandChain{id: id, commands: commands, created: time.Now()}

	return lastCommandChain
}

func NewCommandButton(button string) *chainElement {
	return &chainElement{command: button, containButton: button}
}

func NewCommandMessage(message string) *chainElement {
	return &chainElement{command: message, containText: message}
}

/*func ForwardMessage(chatID int64, messgesID []int64) {
	go GetClient().ForwardMessages(chatID, BotID, messgesID, false, false, false)
}*/
