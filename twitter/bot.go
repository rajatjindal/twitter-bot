package twitter

import (
	"context"
	"net/http"

	"github.com/dghubble/oauth1"
	"github.com/g8rswimmer/go-twitter/v2"
)

const (
	twitterAPIHost = "https://api.twitter.com"
	twitterAPIBase = twitterAPIHost + "/1.1"
)

type Bot struct {
	config *BotConfig
	client *http.Client

	asOwnerOfApp *twitter.Client
	asAppItself  *twitter.Client
}

type BotConfig struct {
	Tokens        Tokens        `json:"tokens"`
	WebhookConfig WebhookConfig `json:"webhookConfig"`
}

type Tokens struct {
	ConsumerKey   string `json:"consumerKey"`
	ConsumerToken string `json:"consumerToken"`
	Token         string `json:"token"`
	TokenSecret   string `json:"tokenSecret"`
	ClientId      string `json:"clientId"`
	ClientSecret  string `json:"clientSecret"`
}

type WebhookConfig struct {
	Path             string `json:"path"`
	URL              string `json:"url"`
	Environment      string `json:"environment"`
	MaxAllowed       int    `json:"maxAllowed"`
	OverWriteOnLimit bool   `json:"overwriteOnLimit"`
}

func NewBotWithClient(client *http.Client, config *BotConfig) (*Bot, error) {
	oauthConfig := oauth1.NewConfig(config.Tokens.ConsumerKey, config.Tokens.ConsumerToken)
	oauthToken := oauth1.NewToken(config.Tokens.Token, config.Tokens.TokenSecret)

	ctx := oauth1.NoContext
	if client != nil {
		ctx = context.WithValue(ctx, oauth1.HTTPClient, client)
	}

	oauthClient := oauthConfig.Client(ctx, oauthToken)

	// setup bot
	bot := &Bot{
		config: config,
		client: client,
	}

	// add client to use when making api call as owner of app
	bot.asOwnerOfApp = &twitter.Client{
		Host:       twitterAPIHost,
		Client:     oauthClient,
		Authorizer: &noop{},
	}

	oauth2Token, err := bot.oauth2Token(config.Tokens.ConsumerKey, config.Tokens.ConsumerToken)
	if err != nil {
		return nil, err
	}

	// add client to use when making api call as app itself
	// https://developer.twitter.com/en/docs/authentication/oauth-2-0
	bot.asAppItself = &twitter.Client{
		Authorizer: appAuth{
			Token: oauth2Token,
		},
		Host:   twitterAPIHost,
		Client: client,
	}

	return bot, nil
}

func (b *Bot) AsAppItself() *twitter.Client {
	return b.asAppItself
}

func (b *Bot) AsOwnerOfApp() *twitter.Client {
	return b.asOwnerOfApp
}
