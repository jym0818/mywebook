package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jym/mywebook/internal/domain"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/http"
	"net/url"
)

var redirectURI = url.PathEscape("https://meoying.com/oauth2/wechat/callback")

type Service interface {
	Authurl(ctx context.Context) (string, error)
	VerifyCode(ctx context.Context, code string, state string) (domain.WechatInfo, error)
}
type service struct {
	appID     string
	appSecret string

	client *http.Client
}

func (s *service) VerifyCode(ctx context.Context, code string, state string) (domain.WechatInfo, error) {
	const targetPattern = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	target := fmt.Sprintf(targetPattern, s.appID, s.appSecret, code)

	req, err := http.NewRequest(http.MethodGet, target, nil)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	var result Result
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return domain.WechatInfo{}, err
	}

	if result.ErrCode != 0 {
		return domain.WechatInfo{},
			fmt.Errorf("微信返回错误响应，错误码：%d，错误信息：%s", result.ErrCode, result.ErrMsg)
	}
	return domain.WechatInfo{
		UnionID: result.UnionID,
		OpenID:  result.OpenID,
	}, nil
}

func Newservice(appID, appSecret string) Service {
	return &service{
		appID:     appID,
		appSecret: appSecret,
	}
}
func (s *service) Authurl(ctx context.Context) (string, error) {
	const urlPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect"
	state := uuid.New()
	return fmt.Sprintf(urlPattern, s.appID, redirectURI, state), nil
}

type Result struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`

	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`

	OpenID  string `json:"openid"`
	Scope   string `json:"scope"`
	UnionID string `json:"unionid"`
}
