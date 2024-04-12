package oauth2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dmzlingyin/utils/config"
	"golang.org/x/oauth2"
)

const (
	DiscordAuthURL  = "https://discord.com/oauth2/authorize"
	DiscordTokenURL = "https://discord.com/api/oauth2/token"
	DiscordUserURL  = "https://discord.com/api/users/@me"
)

const (
	DiscordScopeUser  = "identity"
	DiscordScopeEmail = "email"
)

func NewDiscord() Provider {
	cfg := &oauth2.Config{
		ClientID:     config.GetString("oauth2.discord.client_id"),
		ClientSecret: config.GetString("oauth2.discord.client_secret"),
		Endpoint: oauth2.Endpoint{
			AuthURL:   DiscordAuthURL,
			TokenURL:  DiscordTokenURL,
			AuthStyle: oauth2.AuthStyleInParams,
		},
		RedirectURL: config.GetString("oauth2.discord.redirect_url"),
		Scopes:      []string{DiscordScopeUser, DiscordScopeEmail},
	}

	return &discord{
		cfg: cfg,
	}
}

type discord struct {
	cfg *oauth2.Config
}

func (d *discord) Authorize(ctx context.Context, args *AuthArgs) (token *oauth2.Token, user *User, err error) {
	token, err = d.cfg.Exchange(ctx, args.Code)
	if err != nil {
		return
	} else if !token.Valid() {
		err = errors.New("invalid token")
		return
	}

	res, err := d.cfg.Client(ctx, token).Get(DiscordUserURL)
	if err != nil {
		return
	}
	defer res.Body.Close()

	var u struct {
		ID     string `json:"id"`
		Name   string `json:"username"`
		Email  string `json:"email"`
		Avatar string `json:"avatar"`
	}
	if err = json.NewDecoder(res.Body).Decode(&u); err != nil {
		return nil, nil, err
	}

	var avatar string
	if u.Avatar != "" {
		avatar = fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.jpg", u.ID, u.Avatar)
	}
	return token, &User{
		ID:       u.ID,
		Username: u.Name,
		Email:    u.Email,
		Avatar:   avatar,
	}, nil
}
