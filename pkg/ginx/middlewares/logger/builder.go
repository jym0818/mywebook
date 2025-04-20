package logger

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/jym/mywebook/pkg/logger"
	"io/ioutil"
)

type MiddlewareBuilder struct {
	allowReqBody  bool
	allowRespBody bool
	logger        logger.Logger
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{}
}

// 请求结构体
type AccessLog struct {
	Method string
	Url    string
	//如果请求体很大 要控制大小
	ReqBody  string
	RespBody string
	Status   int
}

// builder模式
func (m *MiddlewareBuilder) AllowReqBody() *MiddlewareBuilder {
	m.allowReqBody = true
	return m
}

// builder模式
func (m *MiddlewareBuilder) AllowRespBody() *MiddlewareBuilder {
	m.allowRespBody = true
	return m
}

func (m *MiddlewareBuilder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		url := c.Request.URL.String()
		if len(url) > 1024 {
			url = url[:1024]
		}
		al := &AccessLog{
			Method: c.Request.Method,
			//URL本身也可以很长 可以考虑不打印全部
			Url: url,
		}

		if m.allowReqBody && c.Request.Body != nil {
			body, _ := c.GetRawData()
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			//这其实很消耗CPU和内存的操作
			al.ReqBody = string(body)
			//读出来就没有了，需要再写回去

		}
		// 替换掉
		if m.allowRespBody {
			c.Writer = responseWriter{
				al:             al,
				ResponseWriter: c.Writer,
			}
		}
		//为什么再这里打印，而不是C.nEXT()，因为有可能panic
		defer func() {
			//读取响应

			//记录日志
			//m.logger.Debug()
		}()

		c.Next()

	}
}

// 装饰器模式
// 什么时候用组合？ 我们只想装饰部分方法的时候
type responseWriter struct {
	al *AccessLog
	gin.ResponseWriter
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.al.RespBody = string(b)
	return w.ResponseWriter.Write(b)
}

func (w responseWriter) WriteString(s string) (int, error) {
	w.al.RespBody = s
	return w.ResponseWriter.WriteString(s)
}

func (w responseWriter) WriteHeader(statusCode int) {
	w.al.Status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)

}
