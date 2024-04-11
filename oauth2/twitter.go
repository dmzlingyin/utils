package oauth2

import (
	"context"
	"encoding/json"
	"github.com/cuigh/auxo/config"
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/auxo/log"
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

func NewTwitter() *Twitter {
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

	return &Twitter{
		cfg:    cfg,
		logger: log.Get("twitter"),
	}
}

type Twitter struct {
	cfg    *oauth2.Config
	logger log.Logger
}

func (d *Twitter) Authorize(ctx context.Context, code string) (token *oauth2.Token, user *User, err error) {
	opt := oauth2.VerifierOption("challenge")
	token, err = d.cfg.Exchange(ctx, code, opt)
	if err != nil {
		return
	} else if !token.Valid() {
		err = errors.New("invalid token")
		return
	}
	token.Expiry = createExpiry()

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

	if u.Data.Avatar == "" {
		u.Data.Avatar = "https://file.aitubo.ai/images/avatars/aituboer.png"
	}

	return token, &User{
		ID:       u.Data.ID,
		Username: u.Data.Name,
		Email:    u.Data.Email,
		Avatar:   u.Data.Avatar,
	}, nil
}
