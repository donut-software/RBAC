package service

import (
	"context"
	"fmt"
	"rbac/internal"

	"go.opentelemetry.io/otel/trace"
)

func (r *RBAC) CreateRole(ctx context.Context, rolename string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Role.Create")
	defer span.End()
	id, err := r.repo.CreateRole(ctx, rolename)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}
	role, err := r.repo.Role(ctx, id)
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}
	err = r.search.IndexRole(ctx, role)
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}
	return nil
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
func (r *RBAC) UpdateRole(ctx context.Context, id string, rolename string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Role.Update")
	defer span.End()
	err := r.repo.UpdateRole(ctx, id, rolename)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}
	role, err := r.repo.Role(ctx, id)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}
	err = r.search.DeleteRole(ctx, id)
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}
	err = r.search.IndexRole(ctx, role)
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}
	return err
}
func (r *RBAC) ListRole(ctx context.Context, args internal.ListArgs) (internal.ListRole, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Account.List")
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
		return fmt.Errorf("search: %w", err)
	}
	return err
}
