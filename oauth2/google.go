package oauth2

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/dmzlingyin/utils/config"
	"github.com/dmzlingyin/utils/oauth2/jwt"
	"golang.org/x/oauth2"
)

const (
	GoogleAuthURL  = "https://accounts.google.com/o/oauth2/v2/auth"
	GoogleTokenURL = "https://oauth2.googleapis.com/token"
	GoogleUserURL  = "https://www.googleapis.com/oauth2/v3/userinfo"
	GoogleKeyURL   = "https://www.googleapis.com/oauth2/v3/certs"
)

const (
	GoogleScopeProfile = "https://www.googleapis.com/auth/userinfo.profile"
	GoogleScopeEmail   = "https://www.googleapis.com/auth/userinfo.email"
)

func NewGoogle() Provider {
	cfg := &oauth2.Config{
		ClientID:     config.GetString("oauth2.google.client_id"),
		ClientSecret: config.GetString("oauth2.google.client_secret"),
		Endpoint: oauth2.Endpoint{
			AuthURL:   GoogleAuthURL,
			TokenURL:  GoogleTokenURL,
			AuthStyle: oauth2.AuthStyleInParams,
		},
		RedirectURL: config.GetString("oauth2.google.redirect_url"),
		Scopes:      []string{GoogleScopeProfile, GoogleScopeEmail},
	}

	return &google{
		cfg:     cfg,
		decoder: jwt.NewDecoder(GoogleKeyURL),
	}
}

type google struct {
	cfg     *oauth2.Config
	decoder *jwt.Decoder
}

func (g *google) Authorize(ctx context.Context, args *AuthArgs) (token *oauth2.Token, user *User, err error) {
	token, err = g.cfg.Exchange(ctx, args.Code)
	if err != nil {
		return
	} else if !token.Valid() {
		err = errors.New("invalid token")
		return
	}

	res, err := g.cfg.Client(ctx, token).Get(GoogleUserURL)
	if err != nil {
		return
	}
	defer res.Body.Close()

	var u struct {
		Sub     string `json:"sub"`
		Name    string `json:"name"`
		Email   string `json:"email"`
		Picture string `json:"picture"`
	}
	if err = json.NewDecoder(res.Body).Decode(&u); err != nil {
		return nil, nil, err
	}
	return token, &User{
		ID:       u.Sub,
		Username: u.Name,
		Email:    u.Email,
		Avatar:   u.Picture,
	}, nil
}
