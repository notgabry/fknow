package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
)

type PDFResponse struct {
	Description string   `json:"description"`
	Documents   []PDFDoc `json:"documents"`
}

type PDFDoc struct {
	URL string `json:"contentUrl"`
}

type SearchResponse struct {
	Content []PDFItem `json:"content"`
}

type PDFItem struct {
	Know  KnowDetails `json:"know"`
	Score float64     `json:"score"`
}

type KnowDetails struct {
	ID     string `json:"uuid"`
	Title  string `json:"title"`
	Likes  int64  `json:"likes"`
	Thumb  string `json:"thumbnailLargeUrl"`
	Knower Knower `json:"knower"`
}

type Knower struct {
	User UserDetails `json:"user"`
}

type UserDetails struct {
	Name string `json:"name"`
}

type PDF struct {
	ID       string
	Title    string
	Score    float64
	ThumbURL string
	Likes    int64
	Knower   string
}

const (
	BaseURL   = "https://apiedge-eu-central-1.knowunity.com"
	UserAgent = "KnowUnityFree Downloader/1.0"
)

func fetch(url string, target interface{}) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", UserAgent)
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	json.NewDecoder(resp.Body).Decode(target)
}

func GetPDF(id string) (string, string) {
	var res PDFResponse
	fetch(fmt.Sprintf("%s/knows/%s", BaseURL, id), &res)

	if len(res.Documents) == 0 {
		return "", ""
	}
	return res.Documents[0].URL, res.Description
}

func ListPDF(query string) []PDF {
	var res SearchResponse
	fetch(fmt.Sprintf("%s/search/knows?query=%s&contentType=KNOW&limit=20&contentLanguageCode=it", BaseURL, url.QueryEscape(query)), &res)

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
