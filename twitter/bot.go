package twitter

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"sigs.k8s.io/yaml"
)

const twitterAPIBase = "https://api.twitter.com/1.1"

//Tokens is for twitter tokens
type Tokens struct {
	ConsumerKey   string `json:"consumerKey"`
	ConsumerToken string `json:"consumerToken"`
	Token         string `json:"token"`
	TokenSecret   string `json:"tokenSecret"`
	ClientId      string `json:"clientId"`
	ClientSecret  string `json:"clientSecret"`
}

//Bot is a twitter bot
type Bot struct {
	config      *BotConfig
	client      *twitter.Client
	debug       bool
	oauthClient *http.Client
}

//BotConfig is config for initializing new twitter bot
type BotConfig struct {
	Tokens               Tokens `json:"tokens"`
	Environment          string `json:"environment"`
	WebhookHost          string `json:"webhook-host"`
	WebhookPath          string `json:"webhook-path"`
	OverrideRegistration bool   `json:"override-registration"`
}

//NotFoundError not found error
type NotFoundError struct {
	Msg string
}

func (n NotFoundError) Error() string {
	return n.Msg
}

func NewBotFromFile(credsFile string) (*Bot, error) {
	data, err := ioutil.ReadFile(credsFile)
	if err != nil {
		return nil, err
	}

	config := &BotConfig{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return NewBot(config)
}

//NewBot returns new bot
func NewBot(config *BotConfig) (*Bot, error) {
	oauthConfig := oauth1.NewConfig(config.Tokens.ConsumerKey, config.Tokens.ConsumerToken)
	oauthToken := oauth1.NewToken(config.Tokens.Token, config.Tokens.TokenSecret)
	oauthClient := oauthConfig.Client(oauth1.NoContext, oauthToken)

	return &Bot{
		oauthClient: oauthClient,
		client:      twitter.NewClient(oauthClient),
		config:      config,
		debug:       true,
	}, nil
}

func (b *Bot) WebhookPath() string {
	return b.config.WebhookPath
}

func (b *Bot) currentWebhook() (*WebhookConfig, error) {
	webhookConfigs, err := b.getAllWebhooks()
	if err != nil {
		return nil, err
	}

	for _, c := range webhookConfigs {
		if strings.HasPrefix(c.URL, b.config.WebhookHost) {
			return c, nil
		}
	}

	return nil, NotFoundError{Msg: fmt.Sprintf("no webhook found for host %s", b.config.WebhookHost)}
}

//DoRegistrationAndSubscribeBusiness registers the webhook for twitter bot
func (b *Bot) DoRegistrationAndSubscribeBusiness() error {
	webhook, err := b.currentWebhook()
	if err != nil {
		_, ok := err.(NotFoundError)
		if !ok {
			return err
		}
	}

	switch {
	case webhook != nil:
		err = b.triggerCRC(webhook.ID)
		if err != nil {
			return err
		}
	case webhook == nil:
		_, err = b.registerWebhook()
		if err != nil {
			return err
		}
	}

	subscribed, err := b.isSubscribed()
	if err != nil {
		return err
	}

	if subscribed {
		return nil
	}

	return b.subscribeWebhook()
}
