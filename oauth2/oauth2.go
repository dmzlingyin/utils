package oauth2

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
)

type OauthType string

const (
	TypeGoogle   OauthType = "google"
	TypeApple    OauthType = "apple"
	TypeFacebook OauthType = "facebook"
	TypeDiscord  OauthType = "discord"
	TypeTwitter  OauthType = "twitter"
	TypeWechat   OauthType = "wechat"
	TypeCasdoor  OauthType = "casdoor"
)

type User struct {
	ID       string // 第三方用户ID
	Username string // 用户名
	Avatar   string // 头像
	Email    string // 邮箱
	Phone    string // 手机
}

type AuthArgs struct {
	// 登录类型
	Type OauthType
	// 授权码
	Code string
	// 随机值, 防止 XSRF 攻击
	State string
	// 客户端可以直接传递token, 省略了code换取token的步骤
	Token string
	// 如果微信登录要获取手机号,需多传一个code
	PCode string
}

type Builder func() Provider

type Provider interface {
	Authorize(ctx context.Context, args *AuthArgs) (*oauth2.Token, *User, error)
}

func New(ots ...OauthType) Provider {
	c := &Client{
		providers: make(map[OauthType]Provider),
		builders: map[OauthType]Builder{
			TypeGoogle:   NewGoogle,
			TypeApple:    NewApple,
			TypeFacebook: NewFacebook,
			TypeDiscord:  NewDiscord,
			TypeTwitter:  NewTwitter,
			TypeWechat:   NewWechat,
			TypeCasdoor:  NewCasdoor,
		},
	}
	for _, ot := range ots {
		if builder, ok := c.builders[ot]; ok {
			c.register(ot, builder())
		}
	}
	return c
}

type Client struct {
	providers map[OauthType]Provider
	builders  map[OauthType]Builder
}

func (c *Client) register(ot OauthType, p Provider) {
	if _, ok := c.providers[ot]; ok {
		panic("duplicate processor: " + ot)
	}
	c.providers[ot] = p
}

func (c *Client) Authorize(ctx context.Context, args *AuthArgs) (*oauth2.Token, *User, error) {
	if p, ok := c.providers[args.Type]; ok {
		return p.Authorize(ctx, args)
	}
	return nil, nil, errors.New(fmt.Sprintf("not supported oauth type: %s", args.Type))
}
