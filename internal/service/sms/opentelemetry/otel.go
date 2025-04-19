package opentelemetry

import (
	"context"
	"github.com/jym/mywebook/internal/service/sms"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	svc    sms.Service
	tracer trace.Tracer
}

func NewService(svc sms.Service) sms.Service {
	tp := otel.GetTracerProvider()
	t := tp.Tracer("guthub.com/jym/mywebook/internal/service/sms/opentelemetry/otel.go")
	return &Service{
		svc:    svc,
		tracer: t,
	}
}

func (s *Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {

	ctx, span := s.tracer.Start(ctx, "top-span")
	defer span.End()

	err := s.svc.Send(ctx, tpl, args, numbers...)
	if err != nil {
		span.RecordError(err)
	}
	return err
}
