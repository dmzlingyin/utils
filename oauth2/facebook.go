package oauth2

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dmzlingyin/utils/config"
	"golang.org/x/oauth2"
)

const (
	FacebookAuthURL  = "https://www.facebook.com/v18.0/dialog/oauth"
	FacebookTokenURL = "https://graph.facebook.com/oauth/access_token"
	FacebookUserURL  = "https://graph.facebook.com/me?fields=id,name,email,picture"
)

const (
	FacebookScopeProfile = "public_profile"
	FacebookScopeEmail   = "email"
	FacebookScopePicture = "user_photos"
)

func NewFacebook() *Facebook {
	return &Facebook{
		cfg: &oauth2.Config{
			ClientID:     config.GetString("oauth2.facebook.client_id"),
			ClientSecret: config.GetString("oauth2.facebook.client_secret"),
			Endpoint: oauth2.Endpoint{
				AuthURL:   FacebookAuthURL,
				TokenURL:  FacebookTokenURL,
				AuthStyle: oauth2.AuthStyleInParams,
			},
			RedirectURL: config.GetString("oauth2.facebook.redirect_url"),
			Scopes:      []string{FacebookScopeProfile, FacebookScopeEmail, FacebookScopePicture},
		},
	}
}

type Facebook struct {
	cfg *oauth2.Config
}

func (g *Facebook) Authorize(ctx context.Context, code string) (*oauth2.Token, *User, error) {
	// code -> token
	token, err := g.cfg.Exchange(ctx, code)
	if err != nil {
		return nil, nil, err
	} else if !token.Valid() {
		return nil, nil, fmt.Errorf("invalid token %w", err)
	}
	token.Expiry = createExpiry()

	res, err := g.cfg.Client(ctx, token).Get(FacebookUserURL)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	var u struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Email   string `json:"email"`
		Picture struct {
			Data struct {
				Height       int    `json:"height"`
				IsSilhouette bool   `json:"is_silhouette"`
				URL          string `json:"url"`
				Width        int    `json:"width"`
			} `json:"data"`
		} `json:"picture"`
	}
	if err = json.NewDecoder(res.Body).Decode(&u); err != nil {
		return nil, nil, err
	}
	return token, &User{
		ID:       u.ID,
		Username: u.Name,
		Email:    u.Email,
		Avatar:   u.Picture.Data.URL,
	}, nil
}
