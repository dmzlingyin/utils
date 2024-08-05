package sms

import (
	"context"
	"github.com/dmzlingyin/utils/config"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	ses "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002"
)

type SendEmailArgs struct {
	ToEmailAddress string
	Subject        string
	TemplateData   string
}

type Ses struct {
	client           *ses.Client
	fromEmailAddress string
	templateID       uint64
}

func New() *Ses {
	secretID := config.GetString("tencent.ses.secret_id")
	secretKey := config.GetString("tencent.ses.secret_key")
	region := config.GetString("tencent.ses.region")

	credential := common.NewCredential(secretID, secretKey)
	client, err := ses.NewClient(credential, region, profile.NewClientProfile())
	if err != nil {
		panic("failed to create ses client")
	}
	return &Ses{
		client:           client,
		fromEmailAddress: config.GetString("tencent.ses.from_email_address"),
		templateID:       config.GetUint64("tencent.ses.template_id"),
	}
}

func (s *Ses) Send(ctx context.Context, args *SendEmailArgs) error {
	req := ses.NewSendEmailRequest()
	req.FromEmailAddress = &s.fromEmailAddress
	req.Destination = []*string{&args.ToEmailAddress}
	req.Subject = &args.Subject
	req.Template = &ses.Template{
		TemplateID:   &s.templateID,
		TemplateData: &args.TemplateData,
	}
	_, err := s.client.SendEmailWithContext(ctx, req)
	return err
}
