package sms

import (
	Error "errors"
	"fmt"
	"github.com/dmzlingyin/utils/config"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
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

func New(keys ...string) *Sms {
	key := "sms"
	if len(keys) > 0 {
		key = keys[0]
	}
	secretID := config.GetString(fmt.Sprintf("tencent.%s.secret_id", key))
	secretKey := config.GetString(fmt.Sprintf("tencent.%s.secret_key", key))
	region := config.GetString(fmt.Sprintf("tencent.%s.region", key))

	credential := common.NewCredential(secretID, secretKey)
	client, err := sms.NewClient(credential, region, profile.NewClientProfile())
	if err != nil {
		panic("failed to create sms client")
	}
	return &Sms{
		client:     client,
		sdkAppID:   config.GetString(fmt.Sprintf("tencent.%s.sdk_app_id", key)),
		signName:   config.GetString(fmt.Sprintf("tencent.%s.sign_name", key)),
		templateID: config.GetString(fmt.Sprintf("tencent.%s.template_id", key)),
	}
}

func (s *Sms) Send(phone, captcha string) error {
	request := sms.NewSendSmsRequest()
	request.SmsSdkAppId = common.StringPtr(s.sdkAppID)
	request.SignName = common.StringPtr(s.signName)
	request.TemplateId = common.StringPtr(s.templateID)
	request.TemplateParamSet = common.StringPtrs([]string{captcha})
	request.PhoneNumberSet = common.StringPtrs([]string{phone})

	res, err := s.client.SendSms(request)
	if err != nil {
		return err
	}
	if res == nil {
		return Error.New("短信发送失败")
	}
	if len(res.Response.SendStatusSet) <= 0 {
		return Error.New("短信发送失败")
	}

	if *res.Response.SendStatusSet[0].Code != "Ok" {
		return Error.New(*res.Response.SendStatusSet[0].Message)
	}
	return nil
}
