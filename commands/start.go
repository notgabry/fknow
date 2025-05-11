package commands

import (
	"fknow/utils"

	tele "gopkg.in/telebot.v4"
)

func Start(c tele.Context, i18n utils.Translator) error {
	return c.Send(i18n.GT("start"), tele.ModeMarkdownV2)
}
