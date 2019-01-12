package tgbothelper

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Arman92/go-tdlib"
)

// BotTimerFunc - обработчик таймера
type BotTimerFunc func(t time.Time)

// HandleCallback - обработчик
type HandleCallback func()

// Bot - tg бот
type Bot struct {
	client       *Client
	message      *MessageData
	command      *Command
	commandChain *CommandChain
	handleBoot   HandleCallback
	handleRoute  routeCallback
	handleTimer  BotTimerFunc
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

// SetMessage - установить информацию о полученном сообщении
func (b *Bot) SetMessage(message *MessageData) {
	b.message = message
}

// SetCommand - установить информацию о отправленной команде
func (b *Bot) SetCommand(command *Command) {
	if command == nil {
		return
	}

	b.command = command
}

// Client - tg клиент
func (b *Bot) Client() *Client {
	return b.client
}

// Start - запуск бота
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
			go b.handleRoute(botMsg)
		case t := <-ticker.C:
			go b.handleTimer(t)
		case <-ch:
			os.Exit(1)
			return
		}
	}
}

func (b *Bot) startBot(chBot chan *MessageData) {
	b.GetClient()
	b.handleBoot()
	eventFilter := func(msg *tdlib.TdMessage) bool {
		return true
	}

	receiver := b.GetClient().AddEventReceiver(&tdlib.UpdateNewMessage{}, eventFilter, 10)
	for newMsg := range receiver.Chan {
		updateMsg := newMsg.(*tdlib.UpdateNewMessage)
		chBot <- ParseMessage(updateMsg)
	}
}

// SendCommand -  отправить сообщение
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

// GetCommand - последняя отправленная комманда
func (b *Bot) GetCommand() *Command {
	return b.command
}

// GetCommandChain - последняя цепочка команд
func (b *Bot) GetCommandChain() *CommandChain {
	return b.commandChain
}

// NewCommandChain - новая цепочка команд
func (b *Bot) NewCommandChain(id string, commands ...*СhainElement) *CommandChain {
	b.commandChain = &CommandChain{id: id, commands: commands, created: time.Now()}
	return b.commandChain
}

// SetHandleRoute - установить обработчк команд
func (b *Bot) SetHandleRoute(initFunc routeCallback) {
	b.handleRoute = initFunc
}

// SetHandleTimer - установить обработчик таймера
func (b *Bot) SetHandleTimer(timerFunc BotTimerFunc) {
	b.handleTimer = timerFunc
}

// SetHandleBoot - установить обработчик
func (b *Bot) SetHandleBoot(bootFunc HandleCallback) {
	b.handleBoot = bootFunc
}

// NewBot - новый бот
func NewBot(APIID, APIHash, accountName string) *Bot {
	bot := Bot{}
	bot.init(APIID, APIHash, accountName)

	return &bot
}
