package tgbothelper

import (
	"fmt"
	"strings"

	tdlib "github.com/Arman92/go-tdlib"
)

var lastMessage *MessageData

// ParseMessage парсинг содержимого сообщения
func ParseMessage(message *tdlib.UpdateNewMessage) {
	lastMessage = &MessageData{
		ChatID:              message.Message.ChatID,
		MessageID:           message.Message.ID,
		Message:             getText(message),
		Buttons:             *NewButtonList(message.Message.ReplyMarkup, message.Message.ChatID, message.Message.ID),
		ButtonsShowKeyboard: getButtonsShowKeyboard(message.Message.ReplyMarkup),
	}
}

// GetLastMessage последнее сообщение из указанного чата
func GetLastMessage(client *tdlib.Client, chatID int64) *MessageData {
	if lastMessage == nil {
		chat, err := client.GetChat(chatID)
		if err != nil {
			return nil
		}

		lastMessage = &MessageData{
			ChatID:              chatID,
			MessageID:           chat.LastMessage.ID,
			Message:             getTextContent(chat.LastMessage.Content),
			Buttons:             *NewButtonList(chat.LastMessage.ReplyMarkup, chatID, chat.LastMessage.ID),
			ButtonsShowKeyboard: getButtonsShowKeyboard(chat.LastMessage.ReplyMarkup),
		}

		fmt.Println("LAST MESSAGE:", lastMessage)
	}

	return lastMessage
}

// MessageData данные сообщения
type MessageData struct {
	ChatID              int64
	MessageID           int64
	Message             string
	Buttons             ButtonList
	ButtonsShowKeyboard []string
}

func (md *MessageData) ContainText(text string) bool {
	return strings.Contains(md.Message, text)
}

func (md *MessageData) SendMessage(text string) {

}

func (md *MessageData) ReplyMessage(text string) {

}

func getText(message *tdlib.UpdateNewMessage) string {
	if message.Message.Content.GetMessageContentEnum() == tdlib.MessageTextType {
		msg := message.Message.Content.(*tdlib.MessageText)
		return msg.Text.Text
	}

	if message.Message.Content.GetMessageContentEnum() == tdlib.MessagePhotoType {
		msg := message.Message.Content.(*tdlib.MessagePhoto)
		return msg.Caption.Text
	}

	return ""
}

func getButtonsShowKeyboard(reply tdlib.ReplyMarkup) []string {
	if reply == nil || reply.GetReplyMarkupEnum() != tdlib.ReplyMarkupShowKeyboardType {
		return []string{}
	}

	buttonsText := []string{}
	rShowKeyboard := reply.(*tdlib.ReplyMarkupShowKeyboard)
	for _, row := range rShowKeyboard.Rows {
		for _, button := range row {
			buttonsText = append(buttonsText, button.Text)
		}
	}

	if len(buttonsText) == 0 && lastMessage != nil {
		return lastMessage.ButtonsShowKeyboard
	}

	return buttonsText
}

func getButtonsInlineKeyboard(reply tdlib.ReplyMarkup) []string {
	return []string{}
}

func getTextContent(content tdlib.MessageContent) string {
	if content.GetMessageContentEnum() == tdlib.MessageTextType {
		msg := content.(*tdlib.MessageText)
		return msg.Text.Text
	}

	if content.GetMessageContentEnum() == tdlib.MessagePhotoType {
		msg := content.(*tdlib.MessagePhoto)
		return msg.Caption.Text
	}

	return ""
}
