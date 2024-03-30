package aduit

import (
	"errors"
	"os"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	imageaudit20191230 "github.com/alibabacloud-go/imageaudit-20191230/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

// 指定文本检测的应用场景 https://help.aliyun.com/document_detail/155010.htm?spm=a2c4g.477836.0.0.630a661eLBlscB

var LABELS = []string{"spam", "politics", "abuse", "terrorism", "porn", "contraband", "ad"}

func NewAuditor() (*Auditor, error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(os.Getenv("AUDIT_KEY_ID")),
		AccessKeySecret: tea.String(os.Getenv("AUDIT_KEY_SECRET")),
		Endpoint:        tea.String("imageaudit.cn-shanghai.aliyuncs.com"),
	}
	client, err := imageaudit20191230.NewClient(config)
	return &Auditor{client: client}, err
}

type Auditor struct {
	client *imageaudit20191230.Client
}

func (a *Auditor) Audit(text string) error {
	task := &imageaudit20191230.ScanTextRequestTasks{
		Content: tea.String(text),
	}
	labels := make([]*imageaudit20191230.ScanTextRequestLabels, 0, len(LABELS))
	for _, label := range LABELS {
		strl := &imageaudit20191230.ScanTextRequestLabels{Label: tea.String(label)}
		labels = append(labels, strl)
	}

	scanTextRequest := &imageaudit20191230.ScanTextRequest{
		Tasks:  []*imageaudit20191230.ScanTextRequestTasks{task},
		Labels: labels,
	}
	runtime := &util.RuntimeOptions{}
	scanTextResponse, err := a.client.ScanTextWithOptions(scanTextRequest, runtime)
	if err != nil {
		return err
	}

	elements := scanTextResponse.Body.Data.Elements
	if len(elements) > 0 {
		res := elements[0].Results
		if len(res) > 0 && *res[0].Suggestion != "pass" {
			return errors.New("not pass")
		}
	}
	return nil
}
