package oauth2

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/dmzlingyin/utils/config"
	"golang.org/x/oauth2"
)

const (
	TwitterAuthURL  = "https://twitter.com/i/oauth2/authorize"
	TwitterTokenURL = "https://api.twitter.com/2/oauth2/token"
	TwitterUserURL  = "https://api.twitter.com/2/users/me?user.fields=profile_image_url"
)

const (
	TwitterScopeUser  = "users.read"
	TwitterScopeTweet = "tweet.read"
)

func NewTwitter() Provider {
	cfg := &oauth2.Config{
		ClientID:     config.GetString("oauth2.twitter.client_id"),
		ClientSecret: config.GetString("oauth2.twitter.client_secret"),
		Endpoint: oauth2.Endpoint{
			AuthURL:   TwitterAuthURL,
			TokenURL:  TwitterTokenURL,
			AuthStyle: oauth2.AuthStyleInParams,
		},
		RedirectURL: config.GetString("oauth2.twitter.redirect_url"),
		Scopes:      []string{TwitterScopeUser, TwitterScopeTweet},
	}

	return &twitter{cfg: cfg}
}

type twitter struct {
	cfg *oauth2.Config
}

func (d *twitter) Authorize(ctx context.Context, args *AuthArgs) (token *oauth2.Token, user *User, err error) {
	opt := oauth2.VerifierOption("challenge")
	token, err = d.cfg.Exchange(ctx, args.Code, opt)
	if err != nil {
		return
	} else if !token.Valid() {
		err = errors.New("invalid token")
		return
	}

	res, err := d.cfg.Client(ctx, token).Get(TwitterUserURL)
	if err != nil {
		return
	}
	defer res.Body.Close()

	type data struct {
		ID     string `json:"id"`
		Name   string `json:"username"`
		Email  string `json:"email"`
		Avatar string `json:"profile_image_url"`
	}
	var u struct {
		Data data `json:"data"`
	}
	if err = json.NewDecoder(res.Body).Decode(&u); err != nil {
		return nil, nil, err
	}

	return token, &User{
		ID:       u.Data.ID,
		Username: u.Data.Name,
		Email:    u.Data.Email,
		Avatar:   u.Data.Avatar,
	}, nil
}
