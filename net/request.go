package net

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/dmzlingyin/utils/backoff"
	"io"
	"net/http"
)

// Request 统一请求方法, 带线性退避重试机制
func Request[T any](method, url string, headers map[string]string, body []byte) (*T, error) {
	var res T
	err := backoff.Retry(3, func() error {
		// 创建请求
		req, err := http.NewRequest(method, url, bytes.NewReader(body))
		if err != nil {
			return err
		}
		// 设置请求头
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		// 发出请求
		response, e := http.DefaultClient.Do(req)
		if e != nil {
			return e
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			info, _ := io.ReadAll(response.Body)
			return errors.New(string(info))
		}
		return json.NewDecoder(response.Body).Decode(&res)
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}
