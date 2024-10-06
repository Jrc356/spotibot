package spotibot

import (
	"errors"

	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"github.com/zmb3/spotify/v2"
	"go.uber.org/zap"
)

// TODO: convert these to interfaces
type spotibot struct {
	SlackSocketHandler *socketmode.SocketmodeHandler
	SpotifyClient      *spotify.Client
	Logger             *zap.SugaredLogger
}

func New(slackSocketmodeHandler *socketmode.SocketmodeHandler, spotifyClient *spotify.Client, logger *zap.SugaredLogger) (spotibot, error) {
	if slackSocketmodeHandler == nil {
		return spotibot{}, errors.New("slackSocketmodeHandler is nil")
	}
	if spotifyClient == nil {
		return spotibot{}, errors.New("spotifyClient is nil")
	}

	return spotibot{
		SlackSocketHandler: slackSocketmodeHandler,
		SpotifyClient:      spotifyClient,
		Logger:             logger,
	}, nil
}

func (bot *spotibot) Run() error {
	bot.SlackSocketHandler.HandleEvents(slackevents.Message, bot.handleMessage)
	bot.SlackSocketHandler.RunEventLoop()
	return nil
}

func (bot *spotibot) handleMessage(evt *socketmode.Event, client *socketmode.Client) {
	msgEvent := evt.Data.(slackevents.MessageEvent)
	// TODO:
	if msgEvent.ChannelType != "channel" {
		return
	}
}
