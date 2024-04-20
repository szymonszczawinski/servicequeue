package user

import (
	"context"
	"log/slog"
	"servicequeue/dummyapi"
	"servicequeue/service"
	"servicequeue/userapi"

	"golang.org/x/sync/errgroup"
)

type userServiceProxy struct {
	service.JobProxy
	impl userapi.IUserService
}

func NewUserServiceProxy(impl userapi.IUserService, queue service.IJobQueue) userapi.IUserService {
	return &userServiceProxy{
		impl:     impl,
		JobProxy: service.NewJobProxy(queue),
	}
}

func (p *userServiceProxy) VoidMethod() {
	p.ExecuteAsync(service.Job{
		Execute: func() {
			p.impl.VoidMethod()
		},
	})
}

type userService struct {
	service.Service
}

func NewUserService(eg *errgroup.Group, ctx context.Context) *userService {
	return &userService{
		Service: service.NewService("user", eg, ctx),
	}
}

func (s *userService) Start() {
	s.Service.Start()
	serviceProvider := service.GetServiceProvider()
	serviceProvider.RegisterService("user", NewUserServiceProxy(s, s.GetQueue()))
}

func (s *userService) VoidMethod() {
	slog.Info("user service void")
	serviceProvider := service.GetServiceProvider()
	d := serviceProvider.GetService("dummy")
	dummyService, ok := d.(dummyapi.IDummyService)
	if ok {
		dummyService.RegisterUserListener(NewUserListenerProxy(s, s.GetQueue()))
	} else {
		slog.Error("user service void not ok")
	}
}

func (s *userService) OnUser(user string) {
	slog.Info("user received", "user", user)
}

type userListenerProxy struct {
	service.JobProxy
	impl userapi.IUserListener
}

func NewUserListenerProxy(impl userapi.IUserListener, queue service.IJobQueue) userapi.IUserListener {
	return &userListenerProxy{
		impl:     impl,
		JobProxy: service.NewJobProxy(queue),
	}
}

func (p *userListenerProxy) OnUser(user string) {
	p.ExecuteAsync(service.Job{
		Execute: func() {
			p.impl.OnUser(user)
		},
	})
}
