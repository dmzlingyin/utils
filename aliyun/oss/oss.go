package oss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/dmzlingyin/utils/backoff"
	"github.com/google/uuid"
	"strconv"

	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

type ObjectMeta struct {
	Size int
}

func NewClient(opts ...Option) (*Client, error) {
	endpoint := os.Getenv("OSS_ENDPOINT")
	keyID := os.Getenv("OSS_KEY_ID")
	keySecret := os.Getenv("OSS_KEY_SECRET")
	bucketName := os.Getenv("OSS_BUCKET")

	client, err := oss.New(endpoint, keyID, keySecret)
	if err != nil {
		return nil, err
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return nil, err
	}

	oc := &Client{
		client:     client,
		bucket:     bucket,
		retryCount: 3,
		tmpDir:     "tmp",
	}
	for _, opt := range opts {
		opt.apply(oc)
	}
	return oc, nil
}

type Client struct {
	// oss实例
	client *oss.Client
	// bucket实例
	bucket *oss.Bucket
	// 失败重试次数
	retryCount int
	// 网络资源上传至oss的默认路径
	tmpDir string
}

// Copy 将src复制到dst
func (c *Client) Copy(src, dst string) error {
	options := []oss.Option{
		oss.MetadataDirective(oss.MetaReplace),
		// 指定复制源Object的对象标签到目标 Object。
		oss.TaggingDirective(oss.TaggingCopy),
		// 指定复制源Object的元数据到目标Object。
		//oss.MetadataDirective(oss.MetaCopy),
		// 指定CopyObject操作时是否覆盖同名目标Object。此处设置为true，表示禁止覆盖同名Object。
		oss.ForbidOverWrite(false),
		// 指定Object的存储类型。此处设置为Standard，表示标准存储类型。
		oss.StorageClass("Standard"),
		oss.ContentType("application/octet-stream"),
	}

	if isURL(src) {
		objectName, err := c.PutURLObject(src)
		if err != nil {
			return err
		}
		src = objectName
	}

	return backoff.Retry(c.retryCount, func() error {
		_, err := c.bucket.CopyObject(src, dst, options...)
		return err
	})
}

// Delete 删除单个文件
func (c *Client) Delete(objectName string) error {
	err := backoff.Retry(c.retryCount, func() error {
		return c.bucket.DeleteObject(objectName)
	})
	return err
}

// DeleteMulti 批量删除文件
func (c *Client) DeleteMulti(objectNames []string) error {
	err := backoff.Retry(c.retryCount, func() error {
		_, err := c.bucket.DeleteObjects(objectNames)
		return err
	})
	return err
}

// isNoSuchKeyErr 返回是否为找不到oss资源错误
func (c *Client) isNoSuchKeyErr(err error) bool {
	if v, ok := err.(oss.ServiceError); ok {
		if v.StatusCode == 404 {
			return true
		}
	}
	return false
}

// PutObject 将资源上传到oss
func (c *Client) PutObject(name string, reader io.Reader) error {
	return c.bucket.PutObject(name, reader)
}

// PutURLObject 将网络资源上传到oss
func (c *Client) PutURLObject(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	name := uuid.NewString() + path.Ext(url)
	objectName := path.Join(c.tmpDir, name)

	return objectName, c.bucket.PutObject(objectName, io.Reader(res.Body))
}

// GetObject 获取资源
func (c *Client) GetObject(objectName string) (io.ReadCloser, error) {
	return c.bucket.GetObject(objectName)
}

// GetObjectMeta 获取资源元信息
func (c *Client) GetObjectMeta(objectName string) (*ObjectMeta, error) {
	res := &ObjectMeta{}
	meta, err := c.bucket.GetObjectMeta(objectName)
	if err != nil {
		return res, err
	}
	res.Size, err = strconv.Atoi(meta.Get("Content-Length"))
	return res, err
}

// isURL 判断是否为网络资源
func isURL(objectName string) bool {
	return strings.HasPrefix(objectName, "http://") || strings.HasPrefix(objectName, "https://")
}
