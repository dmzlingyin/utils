package oss

type Option interface {
	apply(*Client)
}

type optionFunc func(*Client)

func (f optionFunc) apply(oc *Client) {
	f(oc)
}

// RetryCount 设置retry的次数
func RetryCount(count int) Option {
	return optionFunc(func(oc *Client) {
		oc.retryCount = count
	})
}

// TmpDir 设置临时目录
func TmpDir(dir string) Option {
	return optionFunc(func(oc *Client) {
		oc.tmpDir = dir
	})
}
