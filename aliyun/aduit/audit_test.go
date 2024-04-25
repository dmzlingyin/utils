package aduit

import (
	"testing"
	"time"
)

func TestAudit(t *testing.T) {
	a, err := NewAuditor()
	if err != nil {
		t.Fatal(err)
	}
	texts := []string{
		"本校小额贷款，安全、快捷、方便、无抵押，随机随贷，当天放款，上门服务。联系weixin 123456",
		"硬长直——简称——污,A 日韩 云盘 超轻午马 魏鑫",
		"王天刚去饭店吃饭后发现自己的车子被刮了，破口大骂是哪个傻逼干的?",
	}
	for _, text := range texts {
		time.Sleep(time.Second)
		if err = a.Audit(text); err == nil {
			t.Fatal("audit error")
		}
	}
	t.Log("audit success")
}
