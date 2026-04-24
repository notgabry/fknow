package listeners

import (
	"fknow/utils"

	"github.com/charmbracelet/log"
	tele "gopkg.in/telebot.v4"
)

func OnAddedToGroup(c tele.Context, i18n utils.Translator) error {
	log.Info("Added to group", "group", c.Chat().Title, "by", c.Sender().Username)
	return c.Send(i18n.GT("addToGroup"), tele.ModeMarkdownV2)
}
