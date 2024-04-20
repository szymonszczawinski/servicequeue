package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"servicequeue/dummy"
	"servicequeue/dummyapi"
	"servicequeue/service"
	"servicequeue/user"
	"servicequeue/userapi"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func main() {
	baseContext, cancel := context.WithCancel(context.Background())
	signalChannel := registerShutdownHook(cancel)
	mainGroup, groupContext := errgroup.WithContext(baseContext)
	newd := dummy.NewDummyService(mainGroup, groupContext)
	newd.Start()

	newu := user.NewUserService(mainGroup, groupContext)
	newu.Start()
	serviceProvider := service.GetServiceProvider()
	d := serviceProvider.GetService("dummy")
	dummyService, ok := d.(dummyapi.IDummyService)
	if ok {
		go func() {
			slog.Info("main dummy void")
			dummyService.VoidMethod("Hello")
			slog.Info("main dummy non-void")
			// blocking on waiting for result
			result := dummyService.NonVoidMethod("World")
			slog.Info("result", "msg", result)
		}()
	}
	u := serviceProvider.GetService("user")
	userService, ok := u.(userapi.IUserService)
	if ok {
		slog.Info("main user void")
		userService.VoidMethod()
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
