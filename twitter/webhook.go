package twitter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/sirupsen/logrus"
)

func listWebhooksEndpoint(environment string) string {
	listWebhooksPath := "account_activity/all/%s/webhooks.json"
	return fmt.Sprintf("%s/%s", twitterAPIBase, fmt.Sprintf(listWebhooksPath, environment))
}

//WebhookConfig from twitter
type WebhookConfig struct {
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

//registerWebhook registers our webhook
func (b *Bot) registerWebhook() (*WebhookConfig, error) {
	req, err := http.NewRequest("POST", registerWebhookEndpoint(b.config.Environment, b.config.WebhookHost, b.config.WebhookPath), nil)
	if err != nil {
		return nil, err
	}

	if b.debug {
		d, err := httputil.DumpRequest(req, true)
		if err != nil {
			return nil, err
		}

		logrus.Debug(string(d))
	}

	resp, err := b.oauthClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if b.debug {
		bodyString := string(bodyBytes)
		logrus.Debug("response is: ", bodyString)
	}

	webhookConfig := &WebhookConfig{}
	err = json.Unmarshal(bodyBytes, webhookConfig)
	if err != nil {
		return nil, err
	}

	return webhookConfig, nil
}

//getAllWebhooks lists our webhook
func (b *Bot) getAllWebhooks() ([]*WebhookConfig, error) {
	req, err := http.NewRequest("GET", listWebhooksEndpoint(b.config.Environment), nil)
	if err != nil {
		return nil, err
	}

	if b.debug {
		d, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			return nil, err
		}

		logrus.Debug(string(d))
	}

	resp, err := b.oauthClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if b.debug {
		bodyString := string(bodyBytes)
		logrus.Debug("response is: ", bodyString)
	}

	webhookConfigs := []*WebhookConfig{}
	err = json.Unmarshal(bodyBytes, &webhookConfigs)
	if err != nil {
		return nil, err
	}

	return webhookConfigs, nil
}

func deleteWebhookEndpoint(webhookID, environment string) string {
	return fmt.Sprintf("%s/account_activity/all/%s/webhooks/%s.json", twitterAPIBase, environment, webhookID)
}

//deleteWebhook deletes the webhook
func (b *Bot) deleteWebhook(webhookID string) error {
	req, err := http.NewRequest("DELETE", deleteWebhookEndpoint(webhookID, b.config.Environment), nil)
	if err != nil {
		return err
	}

	if b.debug {
		d, err := httputil.DumpRequest(req, true)
		if err != nil {
			return err
		}

		logrus.Debug(string(d))
	}

	resp, err := b.oauthClient.Do(req)
	if err != nil {
		return err
	}

	//TODO(rajatjindal): add status code check
	if b.debug {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		bodyString := string(bodyBytes)
		logrus.Debug("response is: ", bodyString)
	}

	return nil
}
