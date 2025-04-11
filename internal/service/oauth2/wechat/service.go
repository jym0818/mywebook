package wechat

import (
	"context"
	"fmt"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/url"
)

var redirectURI = url.PathEscape("https://meoying.com/oauth2/wechat/callback")

type Service interface {
	Authurl(ctx context.Context) (string, error)
}
type service struct {
	appID     string
	appSecret string
}

func Newservice(appID, appSecret string) Service {
	return &service{
		appID:     appID,
		appSecret: appSecret,
	}
}
func (s service) Authurl(ctx context.Context) (string, error) {
	const urlPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect"
	state := uuid.New()
	return fmt.Sprintf(urlPattern, s.appID, redirectURI, state), nil
}
