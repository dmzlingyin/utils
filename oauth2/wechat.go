package oauth2

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dmzlingyin/utils/config"
	"golang.org/x/oauth2"
	"net/http"
	"time"
)

const (
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers = "0123456789"
)

const DefaultAvatar = "https://avatar-04.oss-cn-beijing.aliyuncs.com/assets/avatar/1.jpg"

type WechatLoginResp struct {
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

type Phone struct {
	ErrCode   int       `json:"errcode"`
	ErrMsg    string    `json:"errmsg"`
	PhoneInfo phoneInfo `json:"phone_info"`
}

type phoneInfo struct {
	PhoneNumber     string `json:"phoneNumber"`
	PurePhoneNumber string `json:"purePhoneNumber"`
	CountryCode     string `json:"countryCode"`
}

func NewWechat() *Wechat {
	appid := config.GetString("oauth2.wechat.app_id")
	secret := config.GetString("oauth2.wechat.app_secret")
	if appid == "" || secret == "" {
		panic("the appid or secret of wechat get failed")
	}
	return &Wechat{
		appid:  appid,
		secret: secret,
	}
}

type Wechat struct {
	appid  string
	secret string
}

func (w *Wechat) Authorize(ctx context.Context, code, pcode string) (*oauth2.Token, *User, error) {
	url := "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
	url = fmt.Sprintf(url, w.appid, w.secret, code)

	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	var wResp WechatLoginResp
	if err = json.NewDecoder(resp.Body).Decode(&wResp); err != nil {
		return nil, nil, err
	}
	// 判断微信接口返回的是否是一个异常情况
	if wResp.ErrCode != 0 {
		return nil, nil, errors.New(wResp.ErrMsg)
	}

	phone, err := w.getPhoneNumber(pcode)
	if err != nil {
		return nil, nil, err
	}
	user := &User{
		ID:     wResp.OpenId,
		Avatar: DefaultAvatar,
		Phone:  phone,
	}
	token := &oauth2.Token{
		Expiry: time.Now().Add(time.Hour * 24),
	}

	return token, user, nil
}

func (w *Wechat) getPhoneNumber(code string) (string, error) {
	ac, err := w.getAccessToken()
	if err != nil {
		return "", err
	}

	url := "https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=" + ac
	scode := struct {
		Code string `json:"code"`
	}{code}

	buffer, err := json.Marshal(scode)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewReader(buffer))
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var phone Phone
	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&phone); err != nil {
		return "", err
	}
	return phone.PhoneInfo.PhoneNumber, nil
}

func (w *Wechat) getAccessToken() (string, error) {
	// access token
	type AT struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int32  `json:"expires_int"`
	}

	url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential"
	url = fmt.Sprintf("%s&appid=%s&secret=%s", url, w.appid, w.secret)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var at AT
	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&at); err != nil {
		return "", err
	}
	return at.AccessToken, nil
}
