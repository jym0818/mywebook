package ioc

import (
	"github.com/jym/mywebook/internal/service/sms"
	"github.com/jym/mywebook/internal/service/sms/memory"
)

func InitSMS() sms.Service {
	return memory.NewService()
}
