package misc

import "testing"

func TestShortURL(t *testing.T) {
	url := "https://cs.console.aliyun.com/?spm=5176.12818093_-1363046575.ProductAndResource--ali--widget-product-recent.17.2f8116d0XuQbSm#/k8s/cluster/cb7ac2b804545483fa04388055fd1749a/v2/workload/deployment/detail/prod/findmatch-api/service?type=deployment&ns=prod"
	shortURL, err := ShortURL(url)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(shortURL)
}

func BenchmarkShortURL(b *testing.B) {
	url := "https://cs.console.aliyun.com"
	for i := 0; i < b.N; i++ {
		_, _ = ShortURL(url)
	}
}
