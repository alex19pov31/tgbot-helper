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
	command      *command
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
	b.command = &command{}
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

func (b *Bot) GetClient() *tdlib.Client {
	return b.client.GetClient()
}

func (b *Bot) GetMessage() *MessageData {
	return b.message
}

func (b *Bot) GetCommand() *command {
	return b.command
}

func (b *Bot) GetCommandChain() *CommandChain {
	return b.commandChain
}

func NewBot(APIID, APIHash, accountName string, initRoute routeCallback, initTime BotTimerFunc) *Bot {
	bot := Bot{}
	bot.init(APIID, APIHash, accountName)
	bot.initRoute = initRoute
	bot.initTime = initTime

	return &bot
}
