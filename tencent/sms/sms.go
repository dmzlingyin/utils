package sms

import (
	Error "errors"
	"github.com/dmzlingyin/utils/config"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

// Sms 详情: https://cloud.tencent.com/document/product/382/43199
type Sms struct {
	client     *sms.Client
	sdkAppID   string
	signName   string
	templateID string
}

func New() *Sms {
	secretID := config.GetString("tencent.sms.secret_id")
	secretKey := config.GetString("tencent.sms.secret_key")
	region := config.GetString("tencent.sms.region")

	credential := common.NewCredential(secretID, secretKey)
	client, err := sms.NewClient(credential, region, profile.NewClientProfile())
	if err != nil {
		panic("failed to create sms client")
	}
	return &Sms{
		client:     client,
		sdkAppID:   config.GetString("tencent.sms.sdk_app_id"),
		signName:   config.GetString("tencent.sms.sign_name"),
		templateID: config.GetString("tencent.sms.template_id"),
	}
}

func (s *Sms) Send(phone, captcha string) error {
	request := sms.NewSendSmsRequest()
	request.SmsSdkAppId = common.StringPtr(s.sdkAppID)
	request.SignName = common.StringPtr(s.signName)
	request.TemplateId = common.StringPtr(s.templateID)
	request.TemplateParamSet = common.StringPtrs([]string{captcha})
	request.PhoneNumberSet = common.StringPtrs([]string{phone})

	_, err := s.client.SendSms(request)
	if v, ok := err.(*errors.TencentCloudSDKError); ok {
		return Error.New(v.Message)
	}
	return err
}
