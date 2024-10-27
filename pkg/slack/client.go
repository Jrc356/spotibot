package slack

import (
	"fmt"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"go.uber.org/zap"
)

type ClientConfig struct {
	Logger   *zap.SugaredLogger
	AppToken string
	BotToken string
}

func (cc *ClientConfig) validate() error {
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

func New(config ClientConfig) (*socketmode.Client, error) {
	err := config.validate()
	if err != nil {
		return nil, err
	}

	api := slack.New(
		config.BotToken,
		slack.OptionAppLevelToken(config.AppToken),
	)

	client := socketmode.New(api)
	return client, nil
}
