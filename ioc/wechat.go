package ioc

import "github.com/jym/mywebook/internal/service/oauth2/wechat"

func InitWechat() wechat.Service {
	appId := "123456"
	appSecret := "123456"
	return wechat.Newservice(appId, appSecret)
}
