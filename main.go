package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"syscall"
	"time"

	twitter "github.com/rajatjindal/twitter-bot/twitter"
	"github.com/sirupsen/logrus"
)

const (
	webhookHost     = "https://your-webhook-host"
	webhookPath     = "/webhook/twitter"
	environmentName = "prod" //has to be same as provided when getting tokens from twitter developers console
)

type webhookHandler struct{}

func (wh webhookHandler) handler(w http.ResponseWriter, r *http.Request) {
	x, _ := httputil.DumpRequest(r, true)
	logrus.Info(string(x))
}

func main() {
	bot, err := twitter.NewBot(
		&twitter.BotConfig{
			Tokens: twitter.Tokens{
				ConsumerKey:   "<consumer-key>",
				ConsumerToken: "<consumer-token>",
				Token:         "<token>",
				TokenSecret:   "<token-secret>",
			},
			Environment:          environmentName,
			WebhookHost:          webhookHost,
			WebhookPath:          webhookPath,
			OverrideRegistration: true,
		},
	)
	if err != nil {
		logrus.Fatal(err)
	}

	wh := webhookHandler{}
	http.HandleFunc(webhookPath, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			bot.HandleCRCResponse(w, r)
		case http.MethodPost:
			wh.handler(w, r)
		}
	})

	go func() {
		fmt.Println(http.ListenAndServe(":8080", nil))
	}()

	// give time for http server to start and be ready
	// server needs to be up before we do registration and subscription
	// so that we can respond to CRC request
	time.Sleep(3)

	err = bot.DoRegistrationAndSubscribeBusiness()
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("started listening the events successfully")

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM)
	<-signals
	logrus.Info("Received SIGTERM. Terminating...")
}
