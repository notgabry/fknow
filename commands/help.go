package commands

import (
	"fknow/utils"

	tele "gopkg.in/telebot.v4"
)

func Help(c tele.Context, i18n utils.Translator) error {
	return c.Send(&tele.Photo{File: tele.FromDisk("./assets/know.png"), Caption: i18n.GT("help")}, tele.ModeMarkdownV2)
}
