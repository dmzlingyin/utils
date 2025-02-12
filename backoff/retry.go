package backoff

import (
	"errors"
	"time"
)

func Retry(count int, fn func() error) (err error) {
	if count <= 0 {
		return errors.New("retry count must be positive")
	}
	for i := 0; i < count; i++ {
		if i > 0 {
			// 线性退避
			time.Sleep(time.Second * time.Duration(i))
		}
		if err = fn(); err == nil {
			break
		}
	}
	return
}
