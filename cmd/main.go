package main

import (
	"context"
	"os"

	"github.com/Jrc356/spotibot/pkg/spotibot"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	log := logger.Sugar()

	bot, err := spotibot.New(
		spotibot.SlackConfig{
			AppToken: os.Getenv("SLACK_APP_TOKEN"),
			BotToken: os.Getenv("SLACK_BOT_TOKEN"),
		},
		spotibot.SpotifyConfig{
			ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
			ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		},
		log.With("spotibot.service", "bot"),
		ctx,
	)
	if err != nil {
		log.Panic(err)
	}

	spotibot.Run(bot, ctx)
}
