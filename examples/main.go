package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"syscall"
	"time"

	twitterv2 "github.com/g8rswimmer/go-twitter/v2"
	"github.com/gorilla/mux"
	twitter "github.com/rajatjindal/twitter-bot/v2/twitter"
)

type webhookHandler struct {
	bot *twitter.Bot
}

func (wh webhookHandler) handler(w http.ResponseWriter, r *http.Request) {
	x, _ := httputil.DumpRequest(r, true)
	fmt.Println("webhook received is ", string(x))

	_, err := wh.bot.AsOwnerOfApp().CreateTweet(context.TODO(), twitterv2.CreateTweetRequest{
		Text: "hello back @rajatjindal",
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	client := &http.Client{}
	botConf := &twitter.BotConfig{
		Tokens: twitter.Tokens{
			ConsumerKey:   "consumer-key",
			ConsumerToken: "consumer-token",
			Token:         "token",
			TokenSecret:   "token-secret",
			ClientId:      "oauth-client-id",
			ClientSecret:  "oauth-client-secret",
		},
		WebhookConfig: twitter.WebhookConfig{
			Environment:      "development",
			URL:              "https://your-webhook-host",
			Path:             "/webhook/twitter",
			OverWriteOnLimit: true,
			MaxAllowed:       1,
		},
	}

	bot, err := twitter.NewBotWithClient(client, botConf)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}

	router := mux.NewRouter().StrictSlash(true)
	router.Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// this helps with CRC validation that Twitter do periodically
	router.Methods(http.MethodGet).Path(bot.WebhookPath()).HandlerFunc(bot.HandleCRCResponse)

	// this is your webhook handler
	wh := &webhookHandler{}
	router.Methods(http.MethodPost).Path(bot.WebhookPath()).HandlerFunc(wh.handler)

	go func() {
		server := &http.Server{
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  10 * time.Second,
			Addr:         ":8080",
			Handler:      router,
		}

		fmt.Printf("Starting HTTP server on %s\n", server.Addr)
		err := server.ListenAndServe()
		if err != nil {
			fmt.Printf("ERROR: server.ListendAndServe() failed with %v\n", err)
			os.Exit(1)
		}
	}()

	// give time for http server to start and be ready
	// server needs to be up before we do registration and subscription
	// so that we can respond to CRC request
	time.Sleep(3 * time.Second)

	err = bot.EnsureWebhookIsActive()
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("started listening the events successfully")

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM)
	<-signals
	fmt.Println("Received SIGTERM. Terminating...")
}
