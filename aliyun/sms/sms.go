package sms

import (
	"encoding/json"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/cuigh/auxo/errors"
	"github.com/dmzlingyin/utils/config"
	"github.com/dmzlingyin/utils/misc"
	"strings"
)

type Sms struct {
	client       *dysmsapi.Client
	signName     string
	templateCode string
}

func New() *Sms {
	client, err := dysmsapi.NewClient(&openapi.Config{
		AccessKeyId:     tea.String(config.GetString("sms.access_key_id")),
		AccessKeySecret: tea.String(config.GetString("sms.access_key_secret")),
	})
	if err != nil {
		panic("failed to create sms client")
	}
	return &Sms{
		client:       client,
		signName:     config.GetString("sms.sign_name"),
		templateCode: config.GetString("sms.template_code"),
	}
}

func (s *Sms) Send(phone, captcha string) error {
	b, err := json.Marshal(map[string]any{"code": captcha})
	if err != nil {
		return err
	}

	req := &dysmsapi.SendSmsRequest{
		PhoneNumbers:  tea.String(strings.Join([]string{phone}, ",")),
		SignName:      tea.String(s.signName),
		TemplateCode:  tea.String(s.templateCode),
		TemplateParam: tea.String(string(b)),
	}

	res, err := s.client.SendSms(req)
	if err != nil {
		return err
	}
	if tea.StringValue(res.Body.Code) == "OK" {
		return nil
	}
	return errors.New(tea.StringValue(res.Body.Message))
}
