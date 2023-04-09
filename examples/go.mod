module github.com/rajatjindal/twitter-bot/examples

go 1.20

replace github.com/rajatjindal/twitter-bot/v2 => ../

require (
	github.com/g8rswimmer/go-twitter/v2 v2.1.2
	github.com/gorilla/mux v1.8.0
	github.com/rajatjindal/twitter-bot/v2 v2.0.0-00010101000000-000000000000
)

require github.com/dghubble/oauth1 v0.7.1 // indirect
