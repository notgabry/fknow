package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"

	"github.com/charmbracelet/log"
)

const (
	baseURL         = "https://apiedge-eu-central-1.knowunity.com"
	knowunityWebURL = "https://knowunity.it/knows/u"
	userAgent       = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36 (AKA KnowunityBot/1.2)"
	apiVersionProd  = "mg58wKqwgEYCsg6"
)

var (
	httpClient    = &http.Client{}
	nextDataRegex = regexp.MustCompile(`(?s)<script id="__NEXT_DATA__"[^>]*>(.*?)</script>`)
)

// GetPDF returns the PDF download URL and description for a know UUID or URL.
func GetPDF(id string) (string, string) {
	knowURL := id
	if !strings.HasPrefix(id, "http") {
		knowURL = fmt.Sprintf("%s/%s", knowunityWebURL, id)
	}

	html, err := fetchHTML(knowURL)
	if err != nil {
		log.Error("GetPDF", "err", err)
		return "", ""
	}

	matches := nextDataRegex.FindStringSubmatch(html)
	if len(matches) < 2 {
		return "", ""
	}

	var nd NextData
	if err := json.Unmarshal([]byte(matches[1]), &nd); err != nil {
		return "", ""
	}

	know := nd.Props.PageProps.Know
	if len(know.Documents) == 0 {
		return "", ""
	}
	return know.Documents[0].URL, know.Description
}

// ListPDF searches Knowunity and returns results sorted by score.
func ListPDF(query string) []PDF {
	endpoint := fmt.Sprintf(
		"%s/search/knows?query=%s&contentType=KNOW&limit=20&contentLanguageCode=it",
		baseURL, url.QueryEscape(query),
	)

	var res SearchResponse
	if err := fetchJSON(endpoint, &res); err != nil {
		log.Error("ListPDF", "query", query, "err", err)
		return nil
	}

	pdfs := make([]PDF, len(res.Content))
	for i, item := range res.Content {
		pdfs[i] = PDF{
			ID:       item.Know.ID,
			Title:    item.Know.Title,
			ThumbURL: item.Know.Thumb,
			Score:    item.Score,
			Likes:    item.Know.Likes,
			Knower:   item.Know.Knower.User.Name,
		}
	}
	sort.Slice(pdfs, func(i, j int) bool { return pdfs[i].Score > pdfs[j].Score })
	return pdfs
}
