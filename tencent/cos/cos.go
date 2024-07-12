package cos

import (
	"context"
	"errors"
	"github.com/dmzlingyin/utils/config"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"net/http"
	"net/url"
)

const HashTypeSHA256 = "sha256"

// NewCosClient 返回cos实例 详情: https://cloud.tencent.com/document/product/436/31215
func NewCosClient() (*Client, error) {
	secretID := config.GetString("tencent.cos.secret_id")
	secretKey := config.GetString("tencent.cos.secret_key")
	bucketURL := config.GetString("tencent.cos.bucket_url")
	if secretID == "" || secretKey == "" || bucketURL == "" {
		return nil, errors.New("invalid secret_id  or secret_key")
	}
	u, err := url.Parse(bucketURL)
	if err != nil {
		return nil, errors.New("invalid bucket url")
	}

	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretID,
			SecretKey: secretKey,
		},
	})
	return &Client{client: c}, nil
}

type Client struct {
	client *cos.Client
}

func (c *Client) PutObject(name string, reader io.Reader) error {
	_, err := c.client.Object.Put(context.Background(), name, reader, nil)
	return err
}

func (c *Client) PutFromFile(src, dst string) error {
	_, err := c.client.Object.PutFromFile(context.Background(), dst, src, nil)
	return err
}

func (c *Client) GetFileHash(name, hashType string) (string, error) {
	opt := &cos.GetFileHashOptions{
		CIProcess:   "filehash", // 固定写法
		Type:        hashType,
		AddToHeader: false,
	}
	res, _, err := c.client.CI.GetFileHash(context.Background(), name, opt)
	if err != nil {
		return "", err
	}
	if res != nil && res.FileHashCodeResult != nil {
		return res.FileHashCodeResult.SHA256, nil
	}
	return "", errors.New("empty result")
}
