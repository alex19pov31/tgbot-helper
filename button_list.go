package tgbothelper

import (
	"fmt"
	"strings"

	tdlib "github.com/Arman92/go-tdlib"
)

// ButtonList - список кнопок
type ButtonList struct {
	buttons []Button
}

func (bl *ButtonList) GetList() []Button {
	return bl.buttons
}

// Add - добавить кнопку  список
func (bl *ButtonList) Add(button Button) {
	if bl.buttons == nil {
		bl.buttons = []Button{}
	}

	bl.buttons = append(bl.buttons, button)
}

func (bl *ButtonList) getButtonsType(buttonType string) *ButtonList {
	newBL := &ButtonList{}
	for _, button := range bl.buttons {
		if button.GetType() == buttonType {
			newBL.Add(button)
		}
	}

	return newBL
}

// GetShowKeybordButtons - список кнопок клавиатуры
func (bl *ButtonList) GetShowKeybordButtons() *ButtonList {
	return bl.getButtonsType(ShowKeyboardButtonType)
}

// GetInlineButtons - список inline кнопок
func (bl *ButtonList) GetInlineButtons() *ButtonList {
	return bl.getButtonsType(InlineButtonType)
}

// GetCallbackButtons - список callback кнопок
func (bl *ButtonList) GetCallbackButtons() *ButtonList {
	return bl.getButtonsType(CallbackButtonType)
}

// GetButtonByText - возвращает кнопку по тексту
func (bl *ButtonList) GetButtonByText(text string) *Button {
	for _, button := range bl.buttons {
		if button.GetText() == text {
			return &button
		}
	}

	return nil
}

// GetButtonByContainText - возвращает кнопку по вхождению текста
func (bl *ButtonList) GetButtonByContainText(text string) *Button {
	for _, button := range bl.buttons {
		if strings.Contains(button.GetText(), text) {
			return &button
		}
	}

	return nil
}

// GetButtonByData - возвращает кноку по данным для callback вызова
func (bl *ButtonList) GetButtonByData(data string) *Button {
	for _, button := range bl.buttons {
		if button.GetData() == data {
			return &button
		}
	}

	return nil
}

func (bl *ButtonList) hasButtonType(buttonType string) bool {
	for _, button := range bl.buttons {
		if button.GetType() == buttonType {
			return true
		}
	}

	return false
}

// HasCallbackButton - наличие callback кнопок
func (bl *ButtonList) HasCallbackButton() bool {
	return bl.hasButtonType(CallbackButtonType)
}

// HasInlineButton - наличие inline кнопок
func (bl *ButtonList) HasInlineButton() bool {
	return bl.hasButtonType(InlineButtonType)
}

// HasShowKeyboardButton - наличие кнопок кливиатуры
func (bl *ButtonList) HasShowKeyboardButton() bool {
	return bl.hasButtonType(ShowKeyboardButtonType)
}

// NewButtonList - новый список кнопок из сообщения
func NewButtonList(reply tdlib.ReplyMarkup, chatID, messageID int64) *ButtonList {
	bl := &ButtonList{}
	if reply == nil {
		return bl
	}

	fmt.Println("Set chatID: ", chatID)

	if reply.GetReplyMarkupEnum() == tdlib.ReplyMarkupInlineKeyboardType {
		replyKeyboard := reply.(*tdlib.ReplyMarkupInlineKeyboard)
		for _, row := range replyKeyboard.Rows {
			for _, button := range row {
				if button.Type.GetInlineKeyboardButtonTypeEnum() == tdlib.InlineKeyboardButtonTypeCallbackType {
					btCallback := button.Type.(*tdlib.InlineKeyboardButtonTypeCallback)
					bl.Add(
						newCallbackButton(
							button.Text,
							string(btCallback.Data),
							chatID,
							messageID,
						),
					)
					continue
				}

				bl.Add(newInlineButton(button.Text, "", chatID, messageID))
				continue
			}

		}
	}

	if reply.GetReplyMarkupEnum() == tdlib.ReplyMarkupShowKeyboardType {
		replyKeyboard := reply.(*tdlib.ReplyMarkupShowKeyboard)
		for _, row := range replyKeyboard.Rows {
			for _, button := range row {
				bl.Add(newShowKeyboardButton(button.Text, chatID, messageID))
			}
		}
	}

	return bl
}
