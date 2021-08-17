package service

import (
	"context"
	"fmt"
	"rbac/internal"

	"go.opentelemetry.io/otel/trace"
)

func (r *RBAC) CreateTask(ctx context.Context, taskname string) (string, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Task.Create")
	defer span.End()
	id, err := r.repo.CreateTask(ctx, taskname)
	if err != nil {
		return id, fmt.Errorf("repo: %w", err)
	}
	task, err := r.repo.Task(ctx, id)
	if err != nil {
		return id, fmt.Errorf("repo: %w", err)
	}
	_ = r.msgBroker.TaskCreated(ctx, task)
	return id, nil
}
func (r *RBAC) Task(ctx context.Context, id string) (internal.Tasks, error) {
	fmt.Println("get task")
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Task.Task")
	defer span.End()
	// role, err := r.repo.Task(ctx, id)
	role, err := r.search.GetTask(ctx, id)
	if err != nil {
		return internal.Tasks{}, fmt.Errorf("search: %w", err)
	}
	return role, err
}
func (r *RBAC) UpdateTask(ctx context.Context, task internal.Tasks) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Task.Update")
	defer span.End()
	err := r.repo.UpdateTask(ctx, task.Id, task.Task)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}
	rt, err := r.repo.Task(ctx, task.Id)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}
	_ = r.msgBroker.TaskUpdated(ctx, rt)
	return err
}
func (r *RBAC) DeleteTask(ctx context.Context, id string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Task.Delete")
	defer span.End()
	err := r.repo.DeleteTask(ctx, id)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}
	_ = r.msgBroker.TaskDeleted(ctx, id)
	return err
}

func (r *RBAC) ListTask(ctx context.Context, args internal.ListArgs) (internal.ListTask, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Task.List")
	defer span.End()
	lr, err := r.search.ListTask(ctx, args)
	if err != nil {
		return internal.ListTask{}, fmt.Errorf("search: %w", err)
	}
	return lr, nil

}
