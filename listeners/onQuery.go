package listeners

import (
	"fknow/utils"
	"fmt"

	tele "gopkg.in/telebot.v4"
)

func OnQuery(c tele.Context, i18n utils.Translator) error {
	data := utils.ListPDF(c.Query().Text)
	results := make(tele.Results, len(data))

	for i, know := range data {
		result := &tele.ArticleResult{
			Title:       know.Title,
			Description: fmt.Sprintf("👥 %s - ⭐ %.1f - ❤️ %d", know.Knower, know.Score, know.Likes),
			ThumbURL:    know.ThumbURL,
			ResultBase: tele.ResultBase{
				Content: &tele.InputTextMessageContent{
					Text: fmt.Sprintf("https://knowunity.it/knows/%s", know.ID),
				},
			},
		}
		result.SetResultID(fmt.Sprint(i))
		results[i] = result
	}

	return c.Answer(&tele.QueryResponse{
		Results:   results,
		CacheTime: 10,
	})
}
