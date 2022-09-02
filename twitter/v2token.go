package twitter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	twitterv2 "github.com/g8rswimmer/go-twitter/v2"
)

type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

func oauth2Token(consumerKey, consumerSecret string) (string, error) {
	reader := strings.NewReader("grant_type=client_credentials")
	req, err := http.NewRequest(http.MethodPost, "https://api.twitter.com/oauth2/token", reader)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(consumerKey, consumerSecret)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("expected: %d, got %d", http.StatusOK, resp.StatusCode)
	}

	data := struct {
		TokenType   string `json:"token_type"`
		AccessToken string `json:"access_token"`
	}{}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(raw, &data)
	if err != nil {
		return "", err
	}

	return data.AccessToken, nil
}

func (b *Bot) TwitterV2Client() *twitterv2.Client {
	return b.clientv2
}
