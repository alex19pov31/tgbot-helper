package tgbothelper

import (
	"cw_bot/helpers"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Arman92/go-tdlib"
)

type BotTimerFunc func(t time.Time)

type Bot struct {
	client       *Client
	message      *MessageData
	command      *Command
	commandChain *CommandChain
	initRoute    routeCallback
	initTime     BotTimerFunc
}

func (b *Bot) init(APIID, APIHash, accountName string) {
	b.client = &Client{
		APIID:       APIID,
		APIHash:     APIHash,
		accountName: accountName,
	}

	b.message = &MessageData{}
	b.command = &Command{}
	b.commandChain = &CommandChain{}
}

func (b *Bot) Client() *Client {
	return b.client
}

func (b *Bot) Start() {
	tdlib.SetLogVerbosityLevel(0)
	//tdlib.SetFilePath("./errors.txt")
	// Handle Ctrl+C
	ch := make(chan os.Signal, 2)
	chBot := make(chan *MessageData)
	ticker := time.NewTicker(time.Minute)

	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go b.startBot(chBot)

	for {
		select {
		case botMsg := <-chBot:
			go b.initRoute(botMsg)
		case t := <-ticker.C:
			go b.initTime(t)
		case <-ch:
			b.GetClient().DestroyInstance()
			os.Exit(1)
		}
	}
}

func (b *Bot) startBot(chBot chan *MessageData) {
	b.GetClient()
	eventFilter := func(msg *tdlib.TdMessage) bool {
		return true
	}

	receiver := helpers.GetClient().AddEventReceiver(&tdlib.UpdateNewMessage{}, eventFilter, 10)
	for newMsg := range receiver.Chan {
		updateMsg := newMsg.(*tdlib.UpdateNewMessage)
		chBot <- ParseMessage(updateMsg)
	}
}

// SendMessage -  отправить сообщение
func (b *Bot) SendCommand(text string, chatID, messageID int64) {
	if b.command.isLock() {
		return
	}

	b.command.request = text
	b.command.response = ""
	b.command.timeSend = time.Now()
	b.command.chatID = chatID
	b.command.forwardMessageID = messageID
	SendMessage(b.GetClient(), text, chatID, messageID)
}

// GetClient - tg user api клиент
func (b *Bot) GetClient() *tdlib.Client {
	return b.client.GetClient()
}

// GetMessage - данные сообщения
func (b *Bot) GetMessage() *MessageData {
	return b.message
}

func (b *Bot) GetCommand() *Command {
	return b.command
}

func (b *Bot) GetCommandChain() *CommandChain {
	return b.commandChain
}

func (b *Bot) NewCommandChain(id string, commands ...*СhainElement) *CommandChain {
	b.commandChain = &CommandChain{id: id, commands: commands, created: time.Now()}
	return b.commandChain
}

// NewBot - новый бот
func NewBot(APIID, APIHash, accountName string, initRoute routeCallback, initTime BotTimerFunc) *Bot {
	bot := Bot{}
	bot.init(APIID, APIHash, accountName)
	bot.initRoute = initRoute
	bot.initTime = initTime

	return &bot
}
