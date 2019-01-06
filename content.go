package tgbothelper

import (
	"fmt"

	tdlib "github.com/Arman92/go-tdlib"
)

var lastMessage *MessageData

// ParseMessage парсинг содержимого сообщения
func ParseMessage(message *tdlib.UpdateNewMessage) {
	lastMessage = &MessageData{
		Message:             getText(message),
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
			Message:             getTextContent(chat.LastMessage.Content),
			ButtonsShowKeyboard: getButtonsShowKeyboard(chat.LastMessage.ReplyMarkup),
		}

		fmt.Println("LAST MESSAGE:", lastMessage)
	}

	return lastMessage
}

// MessageData данные сообщения
type MessageData struct {
	Message               string
	ButtonsShowKeyboard   []string
	InlineButtons         []string
	CallbackInlineButtons []string
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
