package spotibot

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	spotifyLink = "open.spotify.com"
)

type SpotifyConfig struct {
	ClientID     string
	ClientSecret string
}

func (config SpotifyConfig) validate() error {
	if config.ClientID == "" {
		return errors.New("could not create spotify client: ClientID empty")
	}

	if config.ClientSecret == "" {
		return errors.New("could not create spotify client: ClientSecret empty")
	}
	return nil
}

func newSpotifyClient(config SpotifyConfig, ctx context.Context) (*spotify.Client, error) {
	if err := config.validate(); err != nil {
		return nil, err
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

func containsSpotifyLink(text string) bool {
	return strings.Contains(text, spotifyLink)
}

func trackFromText(client *spotify.Client, text string, ctx context.Context) (*spotify.FullTrack, error) {
	url, err := linkFromText(text)
	if err != nil {
		return &spotify.FullTrack{}, err
	}
	return trackFromLink(client, url, ctx)

}

func linkFromText(text string) (*url.URL, error) {
	matches, err := regexMatch(fmt.Sprintf(`(?P<URL>https://%s/.*)`, spotifyLink), text)
	if err != nil {
		return &url.URL{}, err
	}
	return url.Parse(matches["URL"])
}

func trackFromLink(client *spotify.Client, link *url.URL, ctx context.Context) (*spotify.FullTrack, error) {
	s := strings.Split(link.Path, "/")
	id := s[len(s)-1]
	track, err := client.GetTrack(ctx, spotify.ID(id))
	if err != nil {
		return &spotify.FullTrack{}, err
	}
	return track, nil
}
