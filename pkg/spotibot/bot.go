package spotibot

import (
	"errors"
	"fmt"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"github.com/zmb3/spotify/v2"
	"go.uber.org/zap"
)

type spotibot struct {
	SlackClient   *socketmode.Client
	SpotifyClient *spotify.Client
	Logger        *zap.SugaredLogger
}

func New(slackClient *socketmode.Client, spotifyClient *spotify.Client, logger *zap.SugaredLogger) (spotibot, error) {
	if slackClient == nil {
		return spotibot{}, errors.New("slackSocketmodeHandler is nil")
	}
	if spotifyClient == nil {
		return spotibot{}, errors.New("spotifyClient is nil")
	}

	return spotibot{
		SlackClient:   slackClient,
		SpotifyClient: spotifyClient,
		Logger:        logger,
	}, nil
}

func (bot *spotibot) Run() error {
	go bot.run()
	return bot.SlackClient.Run()
}

func (bot *spotibot) run() {
	for evt := range bot.SlackClient.Events {
		switch evt.Type {
		case socketmode.EventTypeConnecting:
			bot.Logger.Info("Connecting to Slack with Socket Mode...")
		case socketmode.EventTypeConnectionError:
			bot.Logger.Info("Connection failed. Retrying later...")
		case socketmode.EventTypeConnected:
			bot.Logger.Info("Connected to Slack with Socket Mode.")
		case socketmode.EventTypeHello:
			continue
		case socketmode.EventTypeEventsAPI:
			event := evt.Data.(slackevents.EventsAPIEvent)
			err := bot.handleEvent(event)
			if err != nil {
				bot.Logger.Error(err)
				continue
			}
			bot.SlackClient.Ack(*evt.Request)
		default:
			bot.Logger.Warnf("Unexpected event type received: %s\n", evt.Type)
		}
	}
}

func (bot *spotibot) handleEvent(event slackevents.EventsAPIEvent) error {
	switch ev := event.InnerEvent.Data.(type) {
	case *slackevents.MessageEvent:
		bot.Logger.Infof("heard message event: %+v", ev)
		err := bot.handleMessageEvent(ev)
		if err != nil {
			return err
		}
	default:
		bot.Logger.Warnf("Unexpected inner event: %+v", ev)
	}
	return nil
}

func (bot *spotibot) handleMessageEvent(msgEvent *slackevents.MessageEvent) error {
	if msgEvent.Channel == "" || msgEvent.EventTimeStamp == "" {
		return fmt.Errorf("invalid message event: %+v", msgEvent)
	}
	ref := slack.NewRefToMessage(msgEvent.Channel, msgEvent.EventTimeStamp)
	if ref.Timestamp == "" {
		return errors.New("could not get message ref")
	}
	println(ref.Comment)
	err := bot.SlackClient.AddReaction("+1", ref)
	if err != nil {
		return err
	}
	return nil
}
