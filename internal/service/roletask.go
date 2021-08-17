package service

import (
	"context"
	"fmt"
	"rbac/internal"

	"go.opentelemetry.io/otel/trace"
)

func (r *RBAC) CreateRoleTask(ctx context.Context, roleTask internal.RoleTasks) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "RoleTask.Create")
	defer span.End()
	id, err := r.repo.CreateRoleTasks(ctx, roleTask.Task.Id, roleTask.Role.Id)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}
	rt, err := r.repo.RoleTask(ctx, id)
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}
	_ = r.msgBroker.RoleTaskCreated(ctx, rt)
	return nil
}
func (r *RBAC) RoleTask(ctx context.Context, roleTaskId string) (internal.RoleTasks, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "RoleTask.RoleTask")
	defer span.End()
	role, err := r.search.GetRoleTask(ctx, roleTaskId)
	if err != nil {
		return internal.RoleTasks{}, fmt.Errorf("search: %w", err)
	}
	return role, err
}
func (r *RBAC) UpdateRoleTask(ctx context.Context, roleTask internal.RoleTasks) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "RoleTask.Update")
	defer span.End()
	err := r.repo.UpdateRoleTask(ctx, roleTask.Task.Id, roleTask.Role.Id, roleTask.Id)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}
	rt, err := r.repo.RoleTask(ctx, roleTask.Id)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}
	_ = r.msgBroker.RoleTaskUpdated(ctx, rt)
	return err
}
func (r *RBAC) DeleteRoleTask(ctx context.Context, id string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "RoleTask.Delete")
	defer span.End()
	err := r.repo.DeleteRoleTask(ctx, id)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}
	_ = r.msgBroker.RoleTaskDeleted(ctx, id)
	return err
}

func (r *RBAC) ListRoleTask(ctx context.Context, args internal.ListArgs) (internal.ListRoleTask, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "RoleTask.List")
	defer span.End()
	lacr, err := r.search.ListRoleTask(ctx, args)
	if err != nil {
		return internal.ListRoleTask{}, fmt.Errorf("search: %w", err)
	}
	return lacr, err
}
