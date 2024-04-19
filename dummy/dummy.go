package dummy

import (
	"context"
	"fmt"
	"log/slog"
	"servicequeue/service"

	"golang.org/x/sync/errgroup"
)

type IDummyService interface {
	service.IProxy
	VoidMethod(msg string)
	NonVoidMethod(msg string) string
}

type dummyServiceProxy struct {
	service.JobProxy
	impl IDummyService
}

func NewDummyServiceProxy(impl IDummyService, queue service.IJobQueue) IDummyService {
	return &dummyServiceProxy{
		impl:     impl,
		JobProxy: service.NewJobProxy(queue),
	}
}

func (p *dummyServiceProxy) VoidMethod(msg string) {
	p.ExecuteAsync(service.Job{
		Execute: func() {
			p.impl.VoidMethod(msg)
		},
	})
}

func (p *dummyServiceProxy) NonVoidMethod(msg string) string {
	resChan := make(chan string)

	defer close(resChan)
	p.ExecuteAsync(service.Job{Execute: func() {
		result := p.impl.NonVoidMethod(msg)
		resChan <- result
	}})
	stored := <-resChan
	return stored
}

type dummyService struct {
	service.Service
}

func NewDummyService(eg *errgroup.Group, ctx context.Context) *dummyService {
	return &dummyService{
		Service: service.NewService("dummy", eg, ctx),
	}
}

func (s *dummyService) Start() {
	s.Service.Start()
	serviceProvider := service.GetServiceProvider()
	serviceProvider.RegisterService("dummy", NewDummyServiceProxy(s, s.GetQueue()))
}

func (s *dummyService) VoidMethod(msg string) {
	slog.Info("void method", "msg", msg)
}

func (s *dummyService) NonVoidMethod(msg string) string {
	slog.Info("non void method", "msg", msg)
	return fmt.Sprintf("Hello %v", msg)
}
