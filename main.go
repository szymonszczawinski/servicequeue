package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"servicequeue/dummy"
	"servicequeue/service"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func main() {
	baseContext, cancel := context.WithCancel(context.Background())
	signalChannel := registerShutdownHook(cancel)
	mainGroup, groupContext := errgroup.WithContext(baseContext)
	d := dummy.NewDummyService(mainGroup, groupContext)

	d.Start()

	serviceProvider := service.GetServiceProvider()
	s := serviceProvider.GetService("dummy")
	dummyService, ok := s.(dummy.IDummyService)
	if ok {
		dummyService.VoidMethod("Hello")
		result := dummyService.NonVoidMethod("World")
		slog.Info("result", "msg", result)
	}

	if err := mainGroup.Wait(); err == nil {
		slog.Info("FINISH CORE")
	} else {
		slog.Error("error", "err", err)
	}

	defer close(signalChannel)
}

func registerShutdownHook(cancel context.CancelFunc) chan os.Signal {
	sigCh := make(chan os.Signal, 1)

	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGINT)
	go func() {
		// wait until receiving the signal
		<-sigCh
		slog.Info("shutdown signal")
		cancel()
	}()

	return sigCh
}
