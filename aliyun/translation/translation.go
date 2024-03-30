package translation

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alimt"
	"os"
)

func New() (*Translator, error) {
	region := "cn-hangzhou"
	key := os.Getenv("TRANSLATION_KEY_ID")
	secret := os.Getenv("TRANSLATION_KEY_SECRET")
	client, err := alimt.NewClientWithAccessKey(region, key, secret)
	if err != nil {
		return nil, err
	}
	return &Translator{client: client}, nil
}

type Translator struct {
	client *alimt.Client
}

// Translate 调用阿里云接口实现文本翻译
// https://help.aliyun.com/document_detail/158244.html?spm=a2c4g.197055.0.0.60797bbcfGQ1Up
func (s *Translator) Translate(text string) (res string, err error) {
	request := alimt.CreateTranslateECommerceRequest()
	request.Method = "POST"
	request.FormatType = "text"   //翻译文本的格式
	request.SourceLanguage = "zh" //源语言
	request.SourceText = text     //原文
	request.TargetLanguage = "en" //目标语言
	request.Scene = "title"       //目标语言

	// 发起请求并处理异常
	result, err := s.client.TranslateECommerce(request)
	if err != nil {
		res = text
	} else if len(result.Data.Translated) == 0 {
		res = text
	} else {
		res = result.Data.Translated
	}
	return res, err
}
