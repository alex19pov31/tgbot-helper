package tgbothelper

import tdlib "github.com/Arman92/go-tdlib"

type ButtonList struct {
	buttons []Button
}

func (bl *ButtonList) Add(button Button) {
	if bl.buttons == nil {
		bl.buttons = []Button{}
	}

	bl.buttons = append(bl.buttons, button)
}

func (bl *ButtonList) GetShowKeybordButtons() *ButtonList {
	newBL := &ButtonList{}
	for _, button := range bl.buttons {
		if button.GetType() == ShowKeyboardButtonType {
			newBL.Add(button)
		}
	}

	return newBL
}

func (bl *ButtonList) GetInlineButtons() *ButtonList {
	newBL := &ButtonList{}
	for _, button := range bl.buttons {
		if button.GetType() == InlineButtonType {
			newBL.Add(button)
		}
	}

	return newBL
}

func (bl *ButtonList) GetCallbackButtons() *ButtonList {
	newBL := &ButtonList{}
	for _, button := range bl.buttons {
		if button.GetType() == CallbackButtonType {
			newBL.Add(button)
		}
	}

	return newBL
}

func (bl *ButtonList) GetButtonByText(text string) Button {
	return &baseButton{}
}

func (bl *ButtonList) GetButtonByData(data string) Button {
	return &baseButton{}
}

func NewButtonList(reply tdlib.ReplyMarkup) *ButtonList {
	bl := &ButtonList{}
	if reply == nil {
		return bl
	}

	if reply.GetReplyMarkupEnum() == tdlib.ReplyMarkupInlineKeyboardType {
		replyKeyboard := reply.(*tdlib.ReplyMarkupInlineKeyboard)
		for _, row := range replyKeyboard.Rows {
			for _, button := range row {
				//button.Type.GetInlineKeyboardButtonTypeEnum() == tdlib.InlineKeyboardButtonTypeCallbackType{}

			}
		}
	}

	return bl
}
