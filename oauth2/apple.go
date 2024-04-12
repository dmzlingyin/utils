package oauth2

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/dmzlingyin/utils/config"
	mjwt "github.com/dmzlingyin/utils/oauth2/jwt"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type AppleConfig struct {
	secret      []byte
	keyId       string
	teamId      string
	clientId    string
	redirectUrl string
}

const (
	AppleTokenURL = "https://appleid.apple.com/auth/token"
	AppleKeyURL   = "https://appleid.apple.com/auth/keys"
)

type JwtKeys struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

func NewApple() Provider {
	cfg := &AppleConfig{
		keyId:       config.GetString("oauth2.apple.key_id"),
		teamId:      config.GetString("oauth2.apple.team_id"),
		clientId:    config.GetString("oauth2.apple.client_id"),
		redirectUrl: config.GetString("oauth2.apple.redirect_url"),
	}
	file, err := os.ReadFile("config/key.pem")
	if err != nil {
		panic(err)
	}
	cfg.secret = file
	return &apple{
		cfg:     cfg,
		decoder: mjwt.NewDecoder(AppleKeyURL),
	}
}

type apple struct {
	cfg     *AppleConfig
	decoder *mjwt.Decoder
}

func (a *apple) Authorize(_ context.Context, args *AuthArgs) (token *oauth2.Token, user *User, err error) {
	var idToken string
	token, idToken, err = a.getToken(args.Code)
	if err != nil {
		return
	}

	claims, err := a.decoder.Decode(idToken)
	if err != nil {
		return
	}
	user = &User{
		ID:       claims.Subject,
		Username: claims.Name,
		Avatar:   claims.Picture,
		Email:    claims.Email,
	}
	return
}

func (a *apple) getToken(code string) (token *oauth2.Token, IDToken string, err error) {
	data, err := a.httpRequest("POST", AppleTokenURL, map[string]string{
		"client_id":     a.cfg.clientId,
		"client_secret": a.getAppleSecret(),
		"code":          code,
		"grant_type":    "authorization_code",
		"redirect_uri":  a.cfg.redirectUrl,
	})
	if err != nil {
		return
	}

	var res struct {
		AccessToken      string `json:"access_token"`
		TokenType        string `json:"token_type"`
		ExpiresIn        int    `json:"expires_in"`
		RefreshToken     string `json:"refresh_token"`
		IDToken          string `json:"id_token"`
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
	}
	if err = json.Unmarshal(data, &res); err != nil {
		return
	}
	if res.Error != "" {
		err = errors.New(res.ErrorDescription)
		return
	}

	token = &oauth2.Token{
		AccessToken:  res.AccessToken,
		TokenType:    res.TokenType,
		RefreshToken: res.RefreshToken,
		Expiry:       time.Now().Add(time.Second * time.Duration(res.ExpiresIn)),
	}
	IDToken = res.IDToken
	return
}

func (a *apple) httpRequest(method, addr string, params map[string]string) ([]byte, error) {
	form := url.Values{}
	for k, v := range params {
		form.Set(k, v)
	}

	var request *http.Request
	var err error
	if request, err = http.NewRequest(method, addr, strings.NewReader(form.Encode())); err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var response *http.Response
	if response, err = http.DefaultClient.Do(request); nil != err {
		return nil, err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (a *apple) getAppleSecret() string {
	token := &jwt.Token{
		Header: map[string]interface{}{
			"alg": "ES256",
			"kid": a.cfg.keyId,
		},
		Claims: jwt.MapClaims{
			"iss": a.cfg.teamId,
			"iat": time.Now().Unix(),
			// constraint: exp - iat <= 180 days
			"exp": time.Now().Add(24 * time.Hour).Unix(),
			"aud": "https://appleid.apple.com",
			"sub": a.cfg.clientId,
		},
		Method: jwt.SigningMethodES256,
	}

	ecdsaKey, _ := a.authKeyFromBytes(a.cfg.secret)
	ss, _ := token.SignedString(ecdsaKey)
	return ss
}

func (a *apple) authKeyFromBytes(key []byte) (*ecdsa.PrivateKey, error) {
	var err error

	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, errors.New("token: AuthKey must be a valid .p8 PEM file")
	}

	// Parse the key
	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
		return nil, err
	}

	var pkey *ecdsa.PrivateKey
	var ok bool
	if pkey, ok = parsedKey.(*ecdsa.PrivateKey); !ok {
		return nil, errors.New("token: AuthKey must be of type ecdsa.PrivateKey")
	}

	return pkey, nil
}
