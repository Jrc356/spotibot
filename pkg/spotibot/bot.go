package spotibot

import (
	"context"
	"fmt"

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

func New(
	slackConfig SlackConfig,
	spotifyConfig SpotifyConfig,
	logger *zap.SugaredLogger,
	ctx context.Context,
) (*spotibot, error) {
	if err := slackConfig.validate(); err != nil {
		return &spotibot{}, fmt.Errorf("invalid slack config: %e", err)
	}
	if err := spotifyConfig.validate(); err != nil {
		return &spotibot{}, fmt.Errorf("invalid spotify config: %e", err)
	}

	slackClient, err := newSlackClient(slackConfig)
	if err != nil {
		return &spotibot{}, err
	}
	spotifyClient, err := newSpotifyClient(spotifyConfig, ctx)
	if err != nil {
		return &spotibot{}, err
	}

	return &spotibot{
		SlackClient:   slackClient,
		SpotifyClient: spotifyClient,
		Logger:        logger,
	}, nil
}

func Run(bot *spotibot, ctx context.Context) error {
	go bot.runHandlers(ctx)
	return bot.SlackClient.RunContext(ctx)
}

func (bot *spotibot) runHandlers(ctx context.Context) {
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
			err := bot.handleEvent(event, ctx)
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

func (bot *spotibot) handleEvent(event slackevents.EventsAPIEvent, ctx context.Context) error {
	switch ev := event.InnerEvent.Data.(type) {
	case *slackevents.MessageEvent:
		bot.Logger.Infof("heard message event: %+v", ev.Text)
		err := bot.handleMessageEvent(ev, ctx)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unexpected inner event: %+v", ev)
	}
	return nil
}

func (bot *spotibot) handleMessageEvent(msgEvent *slackevents.MessageEvent, ctx context.Context) error {
	if msgEvent.SubType != "" {
		bot.Logger.Infof("skipping subtype: %s", msgEvent.SubType)
	}

	if !containsSpotifyLink(msgEvent.Text) {
		bot.Logger.Infof("skipping. text does not contain spotify link")
		return nil
	}

	if msgEvent.Channel == "" || msgEvent.EventTimeStamp == "" {
		return fmt.Errorf("invalid message event: %+v", msgEvent)
	}

	err := noticeMessage(&bot.SlackClient.Client, msgEvent)
	if err != nil {
		return err
	}
	track, err := trackFromText(bot.SpotifyClient, msgEvent.Text, ctx)
	if err != nil {
		return err
	}
	fmt.Println(track)
	// bot.updatePlaylist(track)

	return nil
}
