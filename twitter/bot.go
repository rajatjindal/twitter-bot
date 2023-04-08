package twitter

import (
	"context"
	"net/http"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	twitterv2 "github.com/g8rswimmer/go-twitter/v2"
)

const (
	twitterAPIHost = "https://api.twitter.com"
	twitterAPIBase = twitterAPIHost + "/1.1"
)

type Bot struct {
	config    *BotConfig
	twitter   *twitter.Client
	twitterv2 *twitterv2.Client
	client    *http.Client
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
	return &Bot{
		twitter: twitter.NewClient(oauthClient),
		config:  config,
		client:  oauthClient,
	}, nil
}

func (b *Bot) MustEnableV2Client() {
	err := b.EnableV2Client()
	if err != nil {
		panic(err)
	}
}

func (b *Bot) EnableV2Client() error {
	oauth2Token, err := b.oauth2Token(b.config.Tokens.ConsumerKey, b.config.Tokens.ConsumerToken)
	if err != nil {
		return err
	}

	b.twitterv2 = &twitterv2.Client{
		Authorizer: authorize{
			Token: oauth2Token,
		},
		Host: twitterAPIHost,
	}

	return nil
}
