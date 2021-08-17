package service

import (
	"context"
	"fmt"
	"rbac/internal"

	"go.opentelemetry.io/otel/trace"
)

func (r *RBAC) CreateRole(ctx context.Context, rolename string) (string, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Role.Create")
	defer span.End()
	id, err := r.repo.CreateRole(ctx, rolename)
	if err != nil {
		return id, fmt.Errorf("repo: %w", err)
	}
	role, err := r.repo.Role(ctx, id)
	if err != nil {
		return id, fmt.Errorf("search: %w", err)
	}
	_ = r.msgBroker.RoleCreated(ctx, role)
	return id, nil
}
func (r *RBAC) Role(ctx context.Context, id string) (internal.Roles, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Role.Role")
	defer span.End()
	role, err := r.search.GetRole(ctx, id)
	if err != nil {
		return internal.Roles{}, fmt.Errorf("search: %w", err)
	}
	return role, err
}
func (r *RBAC) UpdateRole(ctx context.Context, rl internal.Roles) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Role.Update")
	defer span.End()
	err := r.repo.UpdateRole(ctx, rl.Id, rl.Role)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}
	role, err := r.repo.Role(ctx, rl.Id)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}
	_ = r.msgBroker.RoleUpdated(ctx, role)
	return err
}
func (r *RBAC) ListRole(ctx context.Context, args internal.ListArgs) (internal.ListRole, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Role.List")
	defer span.End()
	lr, err := r.search.ListRole(ctx, args)
	if err != nil {
		return internal.ListRole{}, fmt.Errorf("search: %w", err)
	}
	return lr, nil

}
func (r *RBAC) DeleteRole(ctx context.Context, id string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Role.Delete")
	defer span.End()
	err := r.repo.DeleteRole(ctx, id)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}
	_ = r.msgBroker.RoleDeleted(ctx, id)
	return err
}
