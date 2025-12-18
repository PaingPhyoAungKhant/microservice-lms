package zoom

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/paingphyoaungkhant/asto-microservice/shared/config"
)

type ZoomAuth struct {
	config *config.ZoomConfig
	token  *AccessToken
}

type AccessToken struct {
	Token     string
	ExpiresAt time.Time
}

func NewZoomAuth(cfg *config.ZoomConfig) *ZoomAuth {
	return &ZoomAuth{
		config: cfg,
	}
}

func (za *ZoomAuth) GetAccessToken() (string, error) {
	if za.token != nil && time.Now().Before(za.token.ExpiresAt.Add(-5*time.Minute)) {
		return za.token.Token, nil
	}

	token, expiresIn, err := za.generateAccessToken()
	if err != nil {
		return "", err
	}

	za.token = &AccessToken{
		Token:     token,
		ExpiresAt: time.Now().Add(time.Duration(expiresIn) * time.Second),
	}

	return token, nil
}

func (za *ZoomAuth) generateAccessToken() (string, int, error) {
	tokenURL := "https://zoom.us/oauth/token"
	
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", za.config.ClientID, za.config.ClientSecret)))
	
	data := url.Values{}
	data.Set("grant_type", "account_credentials")
	data.Set("account_id", za.config.AccountID)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", auth))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", 0, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var response struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", 0, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.AccessToken, response.ExpiresIn, nil
}

