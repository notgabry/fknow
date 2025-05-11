package listeners

import (
	"fknow/utils"

	"github.com/charmbracelet/log"
	tele "gopkg.in/telebot.v4"
)

func OnAddedToGroup(c tele.Context, i18n utils.Translator) error {
	log.Info("New user added to group", "group", c.Chat().Title, "user", c.Sender().Username)

	return c.Send(i18n.GT("addToGroup"), tele.ModeMarkdownV2)
}
