package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

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
