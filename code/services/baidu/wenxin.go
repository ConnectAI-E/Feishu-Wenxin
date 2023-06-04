package baidu

import (
	"context"
	"github.com/ConnectAI-E/go-wenxin/baidubce"
	ai_customv1 "github.com/ConnectAI-E/go-wenxin/gen/go/baidubce/ai_custom/v1"
	baidubcev1 "github.com/ConnectAI-E/go-wenxin/gen/go/baidubce/v1"
	"start-feishubot/initialization"
)

type WenXin struct {
	Client *baidubce.Client
}

func LoadWenXin(config initialization.Config) *WenXin {
	var opts []baidubce.Option
	opts = append(opts, baidubce.WithTokenRequest(&baidubcev1.TokenRequest{
		GrantType:    "client_credentials",
		ClientId:     config.WenXinClientId,
		ClientSecret: config.WenXinClientSecret,
	}))

	client, err := baidubce.New(opts...)
	if err != nil {
		panic(err)
	}
	return &WenXin{
		Client: client,
	}
}

func (wenxin *WenXin) Completions(msg []*ai_customv1.
	Message) (
	resp *ai_customv1.Message, err error) {
	ctx := context.Background()
	req := &ai_customv1.ChatCompletionsRequest{
		User:     "feishu-user",
		Messages: msg,
	}
	response, err := wenxin.Client.ChatCompletions(ctx, req)

	if err != nil {
		return nil, err
	}

	return &ai_customv1.Message{
		Role:    "assistant",
		Content: response.Result,
	}, nil
}
