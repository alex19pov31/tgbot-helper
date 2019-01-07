package tgbothelper

import (
	"time"

	"github.com/Arman92/go-tdlib"
)

var lastCommand = &command{}

type command struct {
	request   string
	response  string
	timeSend  time.Time
	lockUntil time.Time
}

func (c *command) isLock() bool {
	return c.lockUntil.Unix() > time.Now().Unix()
}

// GetLastCommandTime возвращает время отправки последней команды
func GetLastCommandTime() time.Time {
	return lastCommand.timeSend
}

// GetLastCommandMessage возвращает последнюю отправленную команду
func GetLastCommandMessage() string {
	return lastCommand.request
}

// SetCommandLock утсановить блокировку на отправку команд
func SetCommandLock(lockUntil time.Time) {
	lastCommand.lockUntil = lockUntil
}

// UnlockCommand разблокировать команды
func UnlockCommand() {
	lastCommand.lockUntil = time.Time{}
}

// ForseSendCommand отправить команду минуя блокировки и не сохраняя информацию о ней
func ForseSendCommand(client *tdlib.Client, text string, chatID int64) error {
	inputMsgTxt := tdlib.NewInputMessageText(tdlib.NewFormattedText(text, nil), false, false)
	_, err := client.SendMessage(chatID, 0, false, true, nil, inputMsgTxt)

	return err
}

// SendCommand отправить команду боту
func SendCommand(client *tdlib.Client, text string, chatID int64) {
	if lastCommand.isLock() {
		return
	}

	lastCommand.request = text
	lastCommand.response = ""
	lastCommand.timeSend = time.Now()
	go ForseSendCommand(client, text, chatID)
}

// SendMessage отправляет сообщение в указанный чат
func SendMessage(client *tdlib.Client, text string, chatID, replyMessageID int64) {
	go func() {
		inputMsgTxt := tdlib.NewInputMessageText(tdlib.NewFormattedText(text, nil), true, true)
		client.SendMessage(chatID, replyMessageID, false, true, nil, inputMsgTxt)
	}()
}