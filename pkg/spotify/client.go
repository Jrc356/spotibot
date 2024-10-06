package spotify

import (
	"context"
	"errors"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"go.uber.org/zap"
	"golang.org/x/oauth2/clientcredentials"
)

type ClientConfig struct {
	Logger       *zap.SugaredLogger
	ClientID     string
	ClientSecret string
}

func New(config ClientConfig, ctx context.Context) (*spotify.Client, error) {
	if config.ClientID == "" {
		return nil, errors.New("could not create spotify client: ClientID empty")
	}

	if config.ClientSecret == "" {
		return nil, errors.New("could not create spotify client: ClientSecret empty")
	}

	creds := &clientcredentials.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		TokenURL:     spotifyauth.TokenURL,
	}
	token, err := creds.Token(ctx)
	if err != nil {
		return nil, err
	}

	httpClient := spotifyauth.New().Client(ctx, token)
	return spotify.New(httpClient), nil
}
