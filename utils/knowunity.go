package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"

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

// ── HTTP helpers ──────────────────────────────────────────────────────────────

func fetchHTML(rawURL string) (string, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status %d", resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	return string(b), err
}

func fetchJSON(rawURL string, target interface{}) error {
	accessToken, err := tokens.getAccessToken()
	if err != nil {
		return fmt.Errorf("getting token: %w", err)
	}

	resp, err := doAPIRequest(rawURL, accessToken)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// On 401 force a refresh and retry once
	if resp.StatusCode == http.StatusUnauthorized {
		if err := tokens.forceRefresh(); err != nil {
			return fmt.Errorf("force refresh: %w", err)
		}
		newToken, _ := tokens.getAccessToken()
		resp2, err := doAPIRequest(rawURL, newToken)
		if err != nil {
			return err
		}
		defer resp2.Body.Close()
		if resp2.StatusCode != http.StatusOK {
			b, _ := io.ReadAll(resp2.Body)
			return fmt.Errorf("status %d after refresh: %s", resp2.StatusCode, string(b))
		}
		return json.NewDecoder(resp2.Body).Decode(target)
	}

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status %d: %s", resp.StatusCode, string(b))
	}
	return json.NewDecoder(resp.Body).Decode(target)
}

func doAPIRequest(rawURL, accessToken string) (*http.Response, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return nil, err
	}
	xMedia, timeFormat := buildXMedia(req.URL.RequestURI())
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("X-Platform", "web")
	req.Header.Set("X-Media", xMedia)
	req.Header.Set("X-Time-Format", timeFormat)
	req.Header.Set("X-Interface-Language", "it")
	req.Header.Set("Accept", "application/json")
	return httpClient.Do(req)
}

// buildXMedia generates the X-Media signature header.
// Formula reverse-engineered from the Knowunity web app JS bundle.
func buildXMedia(urlPath string) (xMedia, timeFormat string) {
	if !strings.HasPrefix(urlPath, "/") {
		urlPath = "/" + urlPath
	}
	r := fmt.Sprintf("%05d", 142+len(urlPath))

	now := time.Now().UTC()
	a := fmt.Sprintf("%04d%02d%02d%02d%02d%02d",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second())

	var timeInt int64
	fmt.Sscanf(a, "%d", &timeInt)

	versionRev := reverseStr(apiVersionProd)
	randPart := rand.Intn(8889) + 1111

	xMedia = fmt.Sprintf("11%s%d%d%s", r, 3*timeInt+4321, randPart, versionRev)
	timeFormat = a
	return
}

func reverseStr(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
