package service

import (
	"context"
	"log/slog"

	"golang.org/x/sync/errgroup"
)

type IService interface {
	Start() error
}

type Service struct {
	queue *jobQueue
	eg    *errgroup.Group
	ctx   context.Context
}

func NewService(name string, eg *errgroup.Group, ctx context.Context) Service {
	queue := NeqJobQueue(name, eg)

	return Service{
		queue: queue,
		ctx:   ctx,
		eg:    eg,
	}
}

func (s *Service) Start() error {
	s.queue.Start(s.ctx)
	return nil
}

func (s Service) GetQueue() IJobQueue {
	return s.queue
}

var providerInstance *ServiceProvider

type ServiceProvider struct {
	services map[string]IProxy
}

func GetServiceProvider() *ServiceProvider {
	if providerInstance == nil {
		providerInstance = &ServiceProvider{
			services: map[string]IProxy{},
		}
	}
	return providerInstance
}

func (sp *ServiceProvider) RegisterService(name string, service IProxy) {
	slog.Info("register service", "service", name)
	sp.services[name] = service
}

func (sp *ServiceProvider) GetService(name string) IProxy {
	return sp.services[name]
}
