package twitter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func (b *Bot) isWebhookCurrent(webhooks []*Webhook) (bool, error) {
	for _, c := range webhooks {
		if !c.Valid {
			continue
		}

		urlOK := false
		pathOK := false

		if strings.HasPrefix(c.URL, b.config.WebhookConfig.URL) {
			urlOK = true
		}

		u, err := url.Parse(c.URL)
		if err != nil {
			return false, err
		}

		if u.Path == b.config.WebhookConfig.Path {
			pathOK = true
		}

		if urlOK && pathOK {
			return true, nil
		}
	}

	return false, nil
}

func (b *Bot) WebhookPath() string {
	return b.config.WebhookConfig.Path
}

func (b *Bot) EnsureWebhookIsActive() error {
	webhooks, err := b.getAllWebhooks()
	if err != nil {
		return err
	}

	currentAlready, err := b.isWebhookCurrent(webhooks)
	if err != nil {
		return err
	}

	if currentAlready {
		return nil
	}

	// TODO: decide what would be best approach here. right now, if OverrideOnLimit is true, we will just delete oldest
	if b.config.WebhookConfig.OverWriteOnLimit && b.config.WebhookConfig.MaxAllowed > 0 && len(webhooks) == b.config.WebhookConfig.MaxAllowed {
		// reached the limit, so delete first one from the existing ones
		err := b.DeleteWebhook(webhooks[0].ID)
		if err != nil {
			return err
		}
	}

	webhook, err := b.registerWebhook()
	if err != nil {
		return err
	}

	defer func() {
		b.TriggerCRC(webhook.ID)
	}()

	subscribed, err := b.isSubscribed()
	if err != nil {
		return err
	}

	if subscribed {
		return nil
	}

	return b.subscribeWebhook()
}

func listWebhooksEndpoint(environment string) string {
	listWebhooksPath := "account_activity/all/%s/webhooks.json"
	return fmt.Sprintf("%s/%s", twitterAPIBase, fmt.Sprintf(listWebhooksPath, environment))
}

// Webhook from twitter
type Webhook struct {
	ID               string `json:"id"`
	URL              string `json:"url"`
	Valid            bool   `json:"valid"`
	CreatedTimestamp string `json:"created_timestamp"`
}

func registerWebhookEndpoint(environment, webhookHost, webhookPath string) string {
	registerWebhookPath := fmt.Sprintf("/account_activity/all/%s/webhooks.json", environment)
	webhookURL := fmt.Sprintf("%s%s", webhookHost, webhookPath)

	fmt.Println(webhookURL)
	return fmt.Sprintf("%s%s?url=%s", twitterAPIBase, registerWebhookPath, url.QueryEscape(webhookURL))
}

// registerWebhook registers our webhook
func (b *Bot) registerWebhook() (*Webhook, error) {
	req, err := http.NewRequest("POST", registerWebhookEndpoint(b.config.WebhookConfig.Environment, b.config.WebhookConfig.URL, b.config.WebhookConfig.Path), nil)
	if err != nil {
		return nil, err
	}

	resp, err := b.asOwnerOfApp.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	webhook := &Webhook{}
	err = json.Unmarshal(bodyBytes, webhook)
	if err != nil {
		return nil, err
	}

	return webhook, nil
}

// getAllWebhooks lists our webhook
func (b *Bot) getAllWebhooks() ([]*Webhook, error) {
	req, err := http.NewRequest("GET", listWebhooksEndpoint(b.config.WebhookConfig.Environment), nil)
	if err != nil {
		return nil, err
	}

	resp, err := b.asOwnerOfApp.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get all webhooks. expected code: %d, actual code: %d", http.StatusOK, resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	webhooks := []*Webhook{}
	err = json.Unmarshal(bodyBytes, &webhooks)
	if err != nil {
		return nil, err
	}

	return webhooks, nil
}

func deleteWebhookEndpoint(webhookID, environment string) string {
	return fmt.Sprintf("%s/account_activity/all/%s/webhooks/%s.json", twitterAPIBase, environment, webhookID)
}

// DeleteWebhook deletes the webhook
func (b *Bot) DeleteWebhook(webhookID string) error {
	req, err := http.NewRequest("DELETE", deleteWebhookEndpoint(webhookID, b.config.WebhookConfig.Environment), nil)
	if err != nil {
		return err
	}

	resp, err := b.asOwnerOfApp.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete webhook. expected code: %d, actual code: %d", http.StatusNoContent, resp.StatusCode)
	}

	return nil
}
