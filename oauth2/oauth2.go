package oauth2

import (
	"context"
	"github.com/cuigh/auxo/app/ioc"
	"github.com/cuigh/auxo/errors"
	"golang.org/x/oauth2"
	"time"
)

const (
	KindGoogle   = "google"
	KindFacebook = "facebook"
	KindApple    = "apple"
	KindMobile   = "mobile"
	KindWechat   = "wechat"
	KindAuth0    = "auth0"
	KindCasdoor  = "casdoor"
	KindDiscord  = "discord"
	KindTwitter  = "twitter"
)

type User struct {
	ID       string
	Username string
	Avatar   string
	Email    string
	Phone    string
}

type AuthArgs struct {
	AppName    string
	Ctx        context.Context
	Type       string
	Code       string
	PCode      string
	MobileNum  string
	MobileCode string
}

type Service interface {
	Authorize(args *AuthArgs) (*oauth2.Token, *User, error)
}

func NewService() Service {
	return &service{
		google:   NewGoogle(),
		apple:    NewApple(),
		facebook: NewFacebook(),
		discord:  NewDiscord(),
		twitter:  NewTwitter(),
		casdoor:  NewCasdoor(),
		wechat:   NewWechat(),
	}
}

type service struct {
	google   *Google
	facebook *Facebook
	apple    *Apple
	wechat   *Wechat
	casdoor  *Casdoor
	discord  *Discord
	twitter  *Twitter
}

func (s *service) Authorize(args *AuthArgs) (*oauth2.Token, *User, error) {
	switch args.Type {
	case KindGoogle:
		return s.google.Authorize(args.Ctx, args.Code, args.AppName)
	case KindFacebook:
		return s.facebook.Authorize(args.Ctx, args.Code)
	case KindApple:
		return s.apple.Authorize(args.Ctx, args.Code)
	case KindWechat:
		return s.wechat.Authorize(args.Ctx, args.Code, args.PCode)
	case KindCasdoor:
		return s.casdoor.Authorize(args.Code)
	case KindDiscord:
		return s.discord.Authorize(args.Ctx, args.Code)
	case KindTwitter:
		return s.twitter.Authorize(args.Ctx, args.Code)
	default:
		return nil, nil, errors.New("not supported auth type")
	}
}

func init() {
	ioc.Put(NewService, ioc.Name("oauth2.service"))
}

// createExpiry 指定默认的登录时效: 2周
func createExpiry() time.Time {
	return time.Now().Add(time.Hour * 24 * 14)
}
