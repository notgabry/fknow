package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
)

type UnityResponsePDF struct {
	Title     string       `json:"title"`
	Documents []ContentPDF `json:"documents"`
}

type ContentPDF struct {
	ContentUrl string `json:"contentUrl"`
}

type UnityReponseList struct {
	Content []KnowList `json:"content"`
}

type KnowList struct {
	Know  ContentList `json:"know"`
	Score float64     `json:"score"`
}

type KnowerList struct {
	User UserList `json:"user"`
}

type UserList struct {
	Name string `json:"name"`
}

type ContentList struct {
	UUID           string     `json:"uuid"`
	Likes          int64      `json:"likes"`
	Title          string     `json:"title"`
	LargeThumbnail string     `json:"thumbnailLargeUrl"`
	Knower         KnowerList `json:"knower"`
}

type Pdfs struct {
	ID       string
	Title    string
	Score    float64
	ThumbURL string
	Likes    int64
	Knower   string
}

const BaseURL string = "https://apiedge-eu-central-1.knowunity.com"
const UserAgent string = "KnowUnityFree Downloader/1.0"

func MakeRequest(url string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}
	}
	req.Header.Set("User-Agent", UserAgent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}
	}

	return body
}

func GetPDF(id string) string {
	body := MakeRequest(fmt.Sprintf("%s/knows/%s", BaseURL, id))
	if len(body) == 0 {
		return ""
	}

	var response UnityResponsePDF
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println(err)
		return ""
	}

	if len(response.Documents) == 0 {
		return ""
	}

	return response.Documents[0].ContentUrl
}

func ListPDF(description string) []Pdfs {
	body := MakeRequest(fmt.Sprintf("%s/search/knows?query=%s&contentType=KNOW&limit=20&contentLanguageCode=it", BaseURL, url.QueryEscape(description)))
	if len(body) == 0 {
		return []Pdfs{}
	}

	var response UnityReponseList
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println(err)
		return []Pdfs{}
	}

	var pdfs []Pdfs
	for _, item := range response.Content {
		pdfs = append(pdfs, Pdfs{
			ID:       item.Know.UUID,
			Title:    item.Know.Title,
			ThumbURL: item.Know.LargeThumbnail,
			Score:    item.Score,
			Likes:    item.Know.Likes,
			Knower:   item.Know.Knower.User.Name,
		})
	}

	sort.Slice(pdfs, func(i, j int) bool {
		return pdfs[i].Score > pdfs[j].Score
	})

	return pdfs
}
