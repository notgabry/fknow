package listeners

import (
	"fknow/utils"
	"fmt"
	"regexp"
	"time"

	"github.com/charmbracelet/log"
	tele "gopkg.in/telebot.v4"
)

func OnText(c tele.Context, i18n utils.Translator) error {
	regex := `([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})`
	re := regexp.MustCompile(regex)
	match := re.FindStringSubmatch(c.Message().Text)

	isGroup := c.Message().FromGroup()

	// if this is not an url and the bot is in a group do nothing
	if match == nil && isGroup {
		return nil

	} else if match == nil {
		return c.Send(i18n.GT("invalidURL"), tele.NoPreview, tele.ModeMarkdownV2)
	}

	url, desc := utils.GetPDF(match[0])
	if url == "" {
		return c.Send(i18n.GT("invalidURL"), tele.NoPreview, tele.ModeMarkdownV2)
	}

	log.Info("New File Requested", "id", match[0], "user", c.Message().Sender.Username)

	a := &tele.ReplyMarkup{}
	btnURL := a.URL(i18n.TOnly("openInApp"), fmt.Sprintf("https://knowunity.it/knows/%s", match[0]))
	btnQuery := a.QueryChat(i18n.TOnly("querySearch"), desc)

	a.Inline(
		a.Row(btnURL),
		a.Row(btnQuery),
	)

	err := c.Send(
		&tele.Document{
			File: tele.FromURL(url),
			FileName: fmt.Sprintf("appunti_%d",
				time.Now().UnixMilli()),
		}, a)

	if err != nil {
		c.Send(i18n.GT("invalidPerms"), tele.ModeMarkdownV2)
	}

	return nil
}
