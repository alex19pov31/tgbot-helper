package tgbothelper

import (
	"strings"

	tdlib "github.com/Arman92/go-tdlib"
)

var lastMessage *MessageData

// ParseMessage парсинг содержимого сообщения
func ParseMessage(message *tdlib.UpdateNewMessage) *MessageData {
	return &MessageData{
		ChatID:              message.Message.ChatID,
		MessageID:           message.Message.ID,
		SenderID:            message.Message.SenderUserID,
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
			SenderID:            chat.LastMessage.SenderUserID,
			Message:             getTextContent(chat.LastMessage.Content),
			Buttons:             *NewButtonList(chat.LastMessage.ReplyMarkup, chatID, chat.LastMessage.ID),
			ButtonsShowKeyboard: getButtonsShowKeyboard(chat.LastMessage.ReplyMarkup),
		}
	}

	return lastMessage
}

// GetMessageByID сообщение по идентификатору
func GetMessageByID(client *tdlib.Client, chatID, messageID int64) *MessageData {
	message, err := client.GetMessage(chatID, messageID)
	if err != nil {
		return nil
	}

	return &MessageData{
		ChatID:              chatID,
		MessageID:           messageID,
		SenderID:            message.SenderUserID,
		Message:             getTextContent(message.Content),
		Buttons:             *NewButtonList(message.ReplyMarkup, chatID, messageID),
		ButtonsShowKeyboard: getButtonsShowKeyboard(message.ReplyMarkup),
	}
}

// MessageData данные сообщения
type MessageData struct {
	ChatID              int64
	SenderID            int32
	MessageID           int64
	Message             string
	Buttons             ButtonList
	ButtonsShowKeyboard []string
}

// ContainText - наличие указанного текста в сообщении
func (md *MessageData) ContainText(text string) bool {
	return strings.Contains(md.Message, text)
}

// SendMessage - отправка сообщения
func (md *MessageData) SendMessage(client *tdlib.Client, text string) {
	SendMessage(client, text, md.ChatID, 0)
}

// ReplyMessage - ответ на сообщение
func (md *MessageData) ReplyMessage(client *tdlib.Client, text string) {
	SendMessage(client, text, md.ChatID, md.MessageID)
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
