package jwt

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"math/big"
	"net/http"
	"strings"
)

type Keys struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type Claims struct {
	jwt.StandardClaims
	Email   string `json:"email,omitempty"`
	Name    string `json:"name,omitempty"`
	Picture string `json:"picture,omitempty"`
}

func NewDecoder(url string) *Decoder {
	return &Decoder{
		keyURL: url,
		cache:  make(map[string]Keys),
	}
}

type Decoder struct {
	keyURL string
	cache  map[string]Keys
}

func (d *Decoder) Decode(token string) (*Claims, error) {
	resp, err := http.Get(d.keyURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data struct {
		Keys []struct {
			Kty string `json:"kty"`
			Kid string `json:"kid"`
			Use string `json:"use"`
			Alg string `json:"alg"`
			N   string `json:"n"`
			E   string `json:"e"`
		} `json:"keys"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	for _, key := range data.Keys {
		d.cache[key.Kid] = key
	}
	return d.getSubFromToken(token)
}

// 获取userID
func (d *Decoder) getSubFromToken(idToken string) (*Claims, error) {
	// 数据由 头部、载荷、签名 三部分组成
	cliTokenArr := strings.Split(idToken, ".")

	// 解析token的header获取kid
	cliHeader, err := jwt.DecodeSegment(cliTokenArr[0])
	if err != nil {
		return nil, err
	}

	var jHeader struct {
		Kid string `json:"kid"`
		Alg string `json:"alg"`
	}
	err = json.Unmarshal(cliHeader, &jHeader)
	if err != nil {
		return nil, err
	}

	// 效验pubKey 及 token
	token, err := jwt.ParseWithClaims(idToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return d.getPublicKey(jHeader.Kid), nil
	})
	if err != nil {
		return nil, err
	} else if token == nil {
		return nil, errors.New("nil token")
	}
	claims, ok := token.Claims.(*Claims)
	if ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("get userID err")
}

func (d *Decoder) getPublicKey(keyId string) *rsa.PublicKey {
	// 获取验证所需的公钥
	var pubKey rsa.PublicKey
	var keys Keys
	if key, ok := d.cache[keyId]; ok {
		keys = key
		nBin, _ := base64.RawURLEncoding.DecodeString(keys.N)
		nData := new(big.Int).SetBytes(nBin)

		eBin, _ := base64.RawURLEncoding.DecodeString(keys.E)
		eData := new(big.Int).SetBytes(eBin)

		pubKey.N = nData
		pubKey.E = int(eData.Uint64())
		return &pubKey
	}
	return nil
}
