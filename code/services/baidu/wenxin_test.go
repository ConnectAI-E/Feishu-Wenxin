package baidu

import (
	ai_customv1 "github.com/ConnectAI-E/go-wenxin/gen/go/baidubce/ai_custom/v1"
	"github.com/k0kubun/pp/v3"
	"start-feishubot/initialization"
	"testing"
)

func TestWenXin_Completions(t *testing.T) {
	config := initialization.LoadConfig("../../config.yaml")

	wenxin := LoadWenXin(*config)

	completions, err := wenxin.Completions([]*ai_customv1.Message{
		{
			Role:    "user",
			Content: "hello",
		},
	})
	if err != nil {
		return
	}
	pp.Println(completions)
}
