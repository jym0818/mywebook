package service

import (
	"context"
	"fmt"
	"github.com/jym/mywebook/internal/repository"
	"github.com/jym/mywebook/internal/service/sms"
	"math/rand"
)

const codeTplId = "1877556"

type CodeService interface {
	Send(ctx context.Context, biz string, phone string) error
	Verify(ctx context.Context, biz string, phone string, inputPhone string) (bool, error)
}

type codeService struct {
	repo repository.CodeRepository
	sms  sms.Service
}

func (c *codeService) Send(ctx context.Context, biz string, phone string) error {
	//生成随机验证码
	code := c.generateCode()
	//存起来
	err := c.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	err = c.sms.Send(ctx, codeTplId, []string{code}, phone)
	return err
}
func (c *codeService) generateCode() string {
	// 六位数，num 在 0, 999999 之间，包含 0 和 999999
	num := rand.Intn(1000000)
	// 不够六位的，加上前导 0
	// 000001
	return fmt.Sprintf("%06d", num)
}

func (c *codeService) Verify(ctx context.Context, biz string, phone string, inputPhone string) (bool, error) {
	return c.repo.Verfiy(ctx, biz, phone, inputPhone)
}

func NewcodeService(repo repository.CodeRepository, sms sms.Service) CodeService {
	return &codeService{
		repo: repo,
		sms:  sms,
	}
}
