package listeners

import (
	"fknow/utils"
	"fmt"
	"regexp"
	"time"

	"github.com/charmbracelet/log"
	tele "gopkg.in/telebot.v4"
)

var uuidRegex = regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)

func OnText(c tele.Context, i18n utils.Translator) error {
	match := uuidRegex.FindString(c.Message().Text)

	if match == "" {
		if c.Message().FromGroup() {
			return nil
		}
		return c.Send(i18n.GT("invalidURL"), tele.NoPreview, tele.ModeMarkdownV2)
	}

	pdfURL, desc := utils.GetPDF(match)
	if pdfURL == "" {
		return c.Send(i18n.GT("invalidURL"), tele.NoPreview, tele.ModeMarkdownV2)
	}

	log.Info("File requested", "id", match, "user", c.Message().Sender.Username)

	markup := &tele.ReplyMarkup{}
	markup.Inline(
		markup.Row(markup.URL(i18n.TOnly("openInApp"), fmt.Sprintf("https://knowunity.it/knows/%s", match))),
		markup.Row(markup.QueryChat(i18n.TOnly("querySearch"), desc)),
	)

	if err := c.Send(&tele.Document{
		File:     tele.FromURL(pdfURL),
		FileName: fmt.Sprintf("appunti_%d", time.Now().UnixMilli()),
	}, markup); err != nil {
		c.Send(i18n.GT("invalidPerms"), tele.ModeMarkdownV2)
	}

	return nil
}
