package spotibot

import (
	"errors"
	"fmt"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type SlackConfig struct {
	AppToken string
	BotToken string
}

func (cc *SlackConfig) validate() error {
	if cc.AppToken == "" {
		return fmt.Errorf("AppToken empty")
	}

	if !strings.HasPrefix(cc.AppToken, "xapp-") {
		return fmt.Errorf("AppToken must have the prefix \"xapp-\"")
	}

	if cc.BotToken == "" {
		return fmt.Errorf("BotToken empty")
	}

	if !strings.HasPrefix(cc.BotToken, "xoxb-") {
		return fmt.Errorf("BotToken must have the prefix \"xoxb-\"")
	}
	return nil
}

func newSlackClient(config SlackConfig) (*socketmode.Client, error) {
	if err := config.validate(); err != nil {
		return nil, err
	}

	api := slack.New(
		config.BotToken,
		slack.OptionAppLevelToken(config.AppToken),
	)

	client := socketmode.New(api)
	return client, nil
}

func noticeMessage(client *slack.Client, msg *slackevents.MessageEvent) error {
	ref := slack.NewRefToMessage(msg.Channel, msg.EventTimeStamp)
	if ref.Timestamp == "" {
		return errors.New("could not get message ref")
	}

	err := client.AddReaction("eyes", ref)
	if err != nil {
		return fmt.Errorf("could not add reaction: %e", err)
	}
	return nil
}
