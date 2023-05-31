package baidu

import (
	"context"
	"encoding/json"
	"fmt"

	"start-feishubot/initialization"

	api "github.com/zjy282/ernie-api"
)

type WenXin struct {
	client *api.Client
}

func LoadWenXin(config initialization.Config) *WenXin {
	ctx := context.Background()
	req := &api.OAuthTokenRequest{
		ClientID:     config.WenXinClientId,
		ClientSecret: config.WenXinClientSecret,
	}

	// token有效期30天
	response, err := api.CreateBCEOAuthToken(ctx, req)
	if err != nil {
		panic(err)
	}
	client := api.NewClientWithConfig(api.DefaultBCEConfig(response.AccessToken))

	return &WenXin{
		client: client,
	}
}

func (wenxin *WenXin) Completions(msg []api.ChatRequestMessage) (resp api.ChatRequestMessage, err error) {

	ctx := context.Background()
	req := &api.ChatRequest{
		User:     "test",
		Messages: msg,
	}
	response, err := wenxin.client.CreateChat(ctx, req)

	if err != nil {
		return api.ChatRequestMessage{}, err
	}

	if data, err := json.Marshal(response); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(string(data))
	}

	return api.ChatRequestMessage{}, nil
}
