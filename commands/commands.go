package commands

import (
	"fknow/utils"

	tele "gopkg.in/telebot.v4"
)

func Start(c tele.Context, i18n utils.Translator) error {
	return c.Send(i18n.GT("start"), tele.ModeMarkdownV2)
}

func Help(c tele.Context, i18n utils.Translator) error {
	return c.Send(&tele.Photo{
		File:    tele.FromDisk("./assets/fknow.png"),
		Caption: i18n.GT("help"),
	}, tele.ModeMarkdownV2)
}
