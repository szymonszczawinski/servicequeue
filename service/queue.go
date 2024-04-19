package service

import (
	"context"
	"log/slog"

	"golang.org/x/sync/errgroup"
)

type Job struct {
	Execute func()
}

type IJobQueue interface {
	Add(Job)
}

type IProxy interface{}

type JobProxy struct {
	queue IJobQueue
}

func NewJobProxy(queue IJobQueue) JobProxy {
	return JobProxy{
		queue: queue,
	}
}

func (jp *JobProxy) ExecuteAsync(job Job) {
	jp.queue.Add(job)
}

type jobQueue struct {
	jobs chan Job
	g    *errgroup.Group
	name string
}

func NeqJobQueue(name string, g *errgroup.Group) *jobQueue {
	slog.Info("new job queue", "name", name)
	return &jobQueue{
		jobs: make(chan Job),
		g:    g,
		name: name,
	}
}

// Start starts a dispatcher.
// This dispatcher will stops when it receive a value from `ctx.Done`.
func (jq *jobQueue) Start(ctx context.Context) {
	slog.Info("job queue start", "name", jq.name)
	jq.g.Go(func() error {
		defer close(jq.jobs)
	Loop:
		for {
			slog.Info("job queue wait for job", "name", jq.name)
			select {
			case <-ctx.Done():
				slog.Info("job queue finish", "name", jq.name)
				break Loop

			case job := <-jq.jobs:

				slog.Info("job queue execute job", "name", jq.name)
				jq.g.Go(func() error {
					job.Execute()
					return nil
				})
			}

		}
		slog.Info("job queue done", "name", jq.name)
		return nil
	})
}

// Add enqueues a job into the queue.
// If the number of enqueued jobs has already reached to the maximum size,
// this will block until the other job has finish and the queue has space to accept a new job.
func (jq *jobQueue) Add(job Job) {
	slog.Info("job queue add job", "name", jq.name)
	jq.g.Go(func() error {
		jq.jobs <- job
		return nil
	})
	slog.Info("job queue job added", "name", jq.name)
}
