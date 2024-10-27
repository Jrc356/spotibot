package main

import (
	"os"

	"github.com/Jrc356/spotibot/pkg/slack"
	"github.com/Jrc356/spotibot/pkg/spotibot"
	"github.com/zmb3/spotify/v2"
	"go.uber.org/zap"
)

func main() {
	// ctx := context.Background()
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	log := logger.Sugar()

	slack_config := slack.ClientConfig{
		AppToken: os.Getenv("SLACK_APP_TOKEN"),
		BotToken: os.Getenv("SLACK_BOT_TOKEN"),
		Logger:   log.With("spotibot.service", "slack"),
	}
	slack_client, err := slack.New(slack_config)
	if err != nil {
		panic(err)
	}

	spotify_client := spotify.Client{}
	// spotify_config := spotify.ClientConfig{
	// 	ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
	// 	ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
	// 	Logger:       log.With("spotibot.service", "spotify"),
	// }
	// spotify_client, err := spotify.New(spotify_config, ctx)
	// if err != nil {
	// 	panic(err)
	// }

	bot, err := spotibot.New(
		slack_client,
		&spotify_client,
		log.With("spotibot.service", "bot"),
	)
	if err != nil {
		panic(err)
	}

	bot.Run()
}
