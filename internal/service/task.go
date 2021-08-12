package service

import (
	"context"
	"fmt"
	"rbac/internal"

	"go.opentelemetry.io/otel/trace"
)

func (r *RBAC) CreateTask(ctx context.Context, taskname string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Task.Create")
	defer span.End()
	err := r.repo.CreateTask(ctx, taskname)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
func (r *RBAC) Task(ctx context.Context, id string) (internal.Tasks, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Task.Task")
	defer span.End()
	role, err := r.repo.Task(ctx, id)
	if err != nil {
		fmt.Println(err)
		return internal.Tasks{}, err
	}
	return role, err
}
func (r *RBAC) UpdateTask(ctx context.Context, id string, taskname string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Task.Update")
	defer span.End()
	err := r.repo.UpdateTask(ctx, id, taskname)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}
func (r *RBAC) DeleteTask(ctx context.Context, id string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Task.Delete")
	defer span.End()
	err := r.repo.DeleteTask(ctx, id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}
