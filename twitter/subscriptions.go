package twitter

import (
	"fmt"
	"io"
	"net/http"
)

func getSubsribeEndpoint(environment string) string {
	subscribeEndpointPath := "account_activity/all/%s/subscriptions.json"
	return fmt.Sprintf("%s/%s", twitterAPIBase, fmt.Sprintf(subscribeEndpointPath, environment))
}

// SubscribeWebhook subscribes
func (b *Bot) subscribeWebhook() error {
	path := getSubsribeEndpoint(b.config.Environment)
	resp, err := b.oauthClient.PostForm(path, nil)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//If response code is 204 it was successful
	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	return fmt.Errorf("could not subscribe the webhook. response below: %s", body)
}

// isSubscribed check if subscription was already done
func (b *Bot) isSubscribed() (bool, error) {
	path := getSubsribeEndpoint(b.config.Environment)
	resp, err := b.oauthClient.Get(path)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	//If response code is 204 it was successful
	if resp.StatusCode == 204 {
		return true, nil
	}

	return false, nil
}
