package service

import (
	"context"
	"fmt"
	"rbac/internal"

	"go.opentelemetry.io/otel/trace"
)

func (r *RBAC) CreateRoleTask(ctx context.Context, taskid string, roleid string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "RoleTask.Create")
	defer span.End()
	id, err := r.repo.CreateRoleTasks(ctx, taskid, roleid)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}
	rt, err := r.repo.RoleTask(ctx, id)
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}
	err = r.search.IndexRoleTask(ctx, rt)
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}
	return nil
}
func (r *RBAC) RoleTask(ctx context.Context, roleTaskId string) (internal.RoleTasks, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "RoleTask.RoleTask")
	defer span.End()
	// role, err := r.repo.RoleTask(ctx, roleTaskId)
	role, err := r.search.GetRoleTask(ctx, roleTaskId)
	if err != nil {
		return internal.RoleTasks{}, fmt.Errorf("search: %w", err)
	}
	return role, err
}
func (r *RBAC) UpdateRoleTask(ctx context.Context, taskId string, roleId string, id string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "RoleTask.Update")
	defer span.End()
	err := r.repo.UpdateRoleTask(ctx, taskId, roleId, id)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}
	rt, err := r.repo.RoleTask(ctx, id)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}
	err = r.search.DeleteRoleTask(ctx, rt.Id)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}
	err = r.search.IndexRoleTask(ctx, rt)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}
	return err
}
func (r *RBAC) DeleteRoleTask(ctx context.Context, id string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "RoleTask.Delete")
	defer span.End()
	err := r.repo.DeleteRoleTask(ctx, id)
	if err != nil {
		fmt.Println(err)
		return err
	}
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
