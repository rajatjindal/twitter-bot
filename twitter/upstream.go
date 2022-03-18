package twitter

import "github.com/dghubble/go-twitter/twitter"

func (b *Bot) TwitterClient() *twitter.Client {
	return b.client
}

func (b *Bot) TwitterConfig() *BotConfig {
	return b.config
}
