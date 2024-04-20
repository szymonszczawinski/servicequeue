package dummy

import (
	"context"
	"fmt"
	"log/slog"
	"servicequeue/dummyapi"
	"servicequeue/service"
	"servicequeue/userapi"
	"time"

	"golang.org/x/sync/errgroup"
)

type dummyServiceProxy struct {
	service.JobProxy
	impl dummyapi.IDummyService
}

func NewDummyServiceProxy(impl dummyapi.IDummyService, queue service.IJobQueue) dummyapi.IDummyService {
	return &dummyServiceProxy{
		impl:     impl,
		JobProxy: service.NewJobProxy(queue),
	}
}

func (p *dummyServiceProxy) VoidMethod(msg string) {
	slog.Info("dummy proxy void")
	p.ExecuteAsync(service.Job{
		Execute: func() {
			slog.Info("job Execute", "job", "VoidMethod")
			p.impl.VoidMethod(msg)
		},
		Name: "VoidMethod",
	})
}

func (p *dummyServiceProxy) NonVoidMethod(msg string) string {
	slog.Info("dummy proxy non-void")
	job := service.NewSyncJob("NonVoidMethod", func() any {
		return p.impl.NonVoidMethod(msg)
	})
	result := p.ExecuteSync(*job)
	stringResult, ok := result.(string)
	if ok {
		return stringResult
	}
	return ""
}

// func (p *dummyServiceProxy) NonVoidMethod(msg string) string {
// 	slog.Info("dummy proxy non-void")
// 	resChan := make(chan string)
//
// 	defer close(resChan)
// 	p.ExecuteAsync(service.Job{
// 		Execute: func() {
// 			result := p.impl.NonVoidMethod(msg)
// 			resChan <- result
// 		},
// 		Name: "NonVoidMethod",
// 	})
// 	stored := <-resChan
// 	return stored
// }

func (p *dummyServiceProxy) RegisterUserListener(listener userapi.IUserListener) {
	slog.Info("dummy proxy register listener")
	p.ExecuteAsync(service.Job{
		Execute: func() {
			p.impl.RegisterUserListener(listener)
		},
		Name: "RegisterListener",
	})
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
	time.Sleep(time.Second * 3)
	slog.Info("void method", "msg", msg)
}

func (s *dummyService) NonVoidMethod(msg string) string {
	slog.Info("non void method", "msg", msg)
	return fmt.Sprintf("Hello %v", msg)
}

func (p *dummyService) RegisterUserListener(listener userapi.IUserListener) {
	slog.Info("register listener")
	listener.OnUser("John Doe")
}
