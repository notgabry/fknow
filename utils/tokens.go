package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/log"
)

const authURL = "https://auth.knowunity.com"

type tokenManager struct {
	mu           sync.Mutex
	accessToken  string
	refreshToken string
	expiresAt    time.Time
}

var tokens = &tokenManager{}

// InitTokens loads the refresh token from .env, immediately fetches a fresh
// access token, and is the only startup call needed.
func InitTokens(refreshToken string) error {
	tokens.mu.Lock()
	defer tokens.mu.Unlock()
	tokens.refreshToken = refreshToken
	return tokens.doRefresh()
}

func (tm *tokenManager) getAccessToken() (string, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	if time.Now().Add(60 * time.Second).After(tm.expiresAt) {
		if err := tm.doRefresh(); err != nil {
			return "", err
		}
	}
	return tm.accessToken, nil
}

func (tm *tokenManager) forceRefresh() error {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	return tm.doRefresh()
}

// doRefresh must be called with tm.mu held.
func (tm *tokenManager) doRefresh() error {
	body, _ := json.Marshal(map[string]string{"refreshToken": tm.refreshToken})
	req, err := http.NewRequest("POST", authURL+"/oauth/token", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("X-Platform", "web")

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("token refresh failed %d: %s", resp.StatusCode, string(b))
	}

	var rr struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&rr); err != nil {
		return err
	}

	tm.accessToken = rr.AccessToken
	if rr.RefreshToken != "" {
		tm.refreshToken = rr.RefreshToken
	}
	tm.expiresAt = jwtExpiry(rr.AccessToken)
	log.Info("Token refreshed", "expires_at", tm.expiresAt.Format(time.RFC3339))
	return nil
}

func jwtExpiry(jwt string) time.Time {
	parts := strings.Split(jwt, ".")
	if len(parts) != 3 {
		return time.Now().Add(30 * time.Minute)
	}
	payload := parts[1]
	switch len(payload) % 4 {
	case 2:
		payload += "=="
	case 3:
		payload += "="
	}
	payload = strings.NewReplacer("-", "+", "_", "/").Replace(payload)
	decoded, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return time.Now().Add(30 * time.Minute)
	}
	var claims struct {
		Exp int64 `json:"exp"`
	}
	if err := json.Unmarshal(decoded, &claims); err != nil || claims.Exp == 0 {
		return time.Now().Add(30 * time.Minute)
	}
	return time.Unix(claims.Exp, 0)
}
