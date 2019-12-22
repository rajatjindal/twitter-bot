`twitter-bot` is a framework to get you started with implementing a twitter webhook based bot in a quick and easy way.

## Apply for twitter developer account

##### - Navigate to https://developer.twitter.com/, and ensure you are logged in. 
##### - Once you are logged in, click on 'Apply' at the right-top corner of the screen

<img src="docs/apply-right-corner-1.png" border="4" width="460" height="250" />

##### - Make sure you read and agree to conditions here: [Restricted usecases](https://developer.twitter.com/en/developer-terms/more-on-restricted-use-cases.html)

##### - Click on 'Apply for a developer acocunt' 
<img src="docs/apply-right-corner.png" border="4" width="460" height="250" />

##### - Select your use case. we will select bot
<img src="docs/select-bot.png" border="4" width="460" height="250" />

##### - Review and complete your details
> Note that email-address and phone-number are mandatory for registring the developer account.

<img src="docs/apply-access-to-api.png" border="4" width="460" height="250" />

##### - Tell us about bot
<img src="docs/tell-about-twitter-bot.png" border="4" width="460" height="250" />

##### - Review the information you provided
<img src="docs/review-the-information.png" border="4" width="460" height="250" />

##### - review the agreement and click submit.
<img src="docs/review-the-agreement.png" border="4" width="460" height="250" />

##### - If approved, you will get email for the same
<img src="docs/developeraccount-approval.png" border="4" width="460" height="250" />


## create a new app

##### - Create a new app: https://developer.twitter.com/en/apps/create
<img src="docs/create-app-info.png" border="4" width="460" height="250" />

##### - Review terms and conditions
<img src="docs/review-developer-terms.png" border="4" width="460" height="250" />

##### - Once created, navigate to Permissions tab and make sure you select the right permissions needed for your bot.
<img src="docs/edit-permissions.png" border="4" width="460" height="250" />

##### - Now navigate back to Keys/Tokens and create access token and access token secret. Save them in a safe place, we will need them again.
<img src="docs/copy-tokens.png" border="4" width="460" height="250" />


## setup new environment

##### - Now we need to create environment. Navigate to https://developer.twitter.com/en/account/environments. 
> For the webhook bot, you should select 'Account Activity API/Sandbox'
<img src="docs/setup-dev-environment.png" border="4" width="460" height="250" />

##### - Now enter 'Dev environment label'. e.g. 'dev' or 'prod' and select the app we just created from the dropdown
<img src="docs/setup-environment-2.png" border="4" width="460" height="250" />

##### - Click complete setup and you are all set to start writing code
<img src="docs/you-are-all-set.png" border="4" width="460" height="250" />



## OK, finally writing the code now

#### create new bot object

```golang
const (
	webhookHost     = "https://your-webhook-domain"
	webhookPath     = "/webhook/twitter"
	//has to be same as provided when getting tokens from twitter developers console
	environmentName = "prod" 
)

// use the tokens saved when following above steps
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

```

#### create a handler for the events you will receive from twitter
```golang

type webhookHandler struct{}

// This handler will handle the events from twitter. You can print it to console, forward it to yourself, 
// reply using twitter api. possibilities are endless.
func (wh webhookHandler) handler(w http.ResponseWriter, r *http.Request) {
	webhookBody, _ := httputil.DumpRequest(r, true)
	logrus.Info(string(x))
}

```

#### start a webserver to handle the webhook event from twitter.
> notice that for `GET` method, you need to add `bot.HandleCRCResponse` as handler. This is required for webhook to remain active

```golang
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

```

#### finally trigger the registration and subscription of the webhook. Only after this is when you will start receiving actual webhook events

```golang
	// give time for http server to start and be ready
	// server needs to be up before we do registration and subscription
	// so that we can respond to CRC request
	time.Sleep(3)

	err = bot.DoRegistrationAndSubscribeBusiness()
	if err != nil {
		logrus.Fatal(err)
	}
```

Look at [main.go](main.go) for complete example

## Contributing

Found an issue in documentation or code, please report [issue](https://github.com/rajatjindal/twitter-bot/issues/new)