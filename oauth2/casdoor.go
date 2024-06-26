package oauth2

import (
	"context"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/dmzlingyin/utils/config"
	"golang.org/x/oauth2"
	"os"
)

func NewCasdoor() Provider {
	endpoint := config.GetString("oauth2.casdoor.endpoint")
	clientID := config.GetString("oauth2.casdoor.client_id")
	clientSecret := config.GetString("oauth2.casdoor.client_secret")
	organization := config.GetString("oauth2.casdoor.organization")
	application := config.GetString("oauth2.casdoor.application")
	file, err := os.ReadFile("config/cert")
	if err != nil {
		panic(err)
	}
	casdoorsdk.InitConfig(endpoint, clientID, clientSecret, string(file), organization, application)
	return &casdoor{}
}

type casdoor struct{}

func (c *casdoor) Authorize(_ context.Context, args *AuthArgs) (token *oauth2.Token, user *User, err error) {
	state := "marmot"
	token, err = casdoorsdk.GetOAuthToken(args.Code, state)
	if err != nil {
		return
	}
	claims, err := casdoorsdk.ParseJwtToken(token.AccessToken)
	if err != nil {
		return
	}
	user = &User{
		ID:       claims.Subject,
		Username: claims.Name,
		Avatar:   claims.Avatar,
		Email:    claims.Email,
		Phone:    claims.Phone,
	}
	return
}
