FROM golang:1.13.4-alpine3.10 as builder

WORKDIR /go/src/github.com/rajatjindal/twitter-bot
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build --ldflags "-s -w" -o twitter-bot main.go

FROM alpine:3.10.3

RUN mkdir -p /home/app

# Add non root user
RUN addgroup -S app && adduser app -S -G app
RUN chown app /home/app

WORKDIR /home/app

USER app

COPY --from=builder /go/src/github.com/rajatjindal/twitter-bot/twitter-bot /usr/local/bin/

ENTRYPOINT "/usr/local/bin/twitter-bot"