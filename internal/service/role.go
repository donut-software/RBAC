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
	err := r.repo.CreateRole(ctx, rolename)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
func (r *RBAC) Role(ctx context.Context, id string) (internal.Roles, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Role.Role")
	defer span.End()
	role, err := r.repo.Role(ctx, id)
	if err != nil {
		fmt.Println(err)
		return internal.Roles{}, err
	}
	return role, err
}
func (r *RBAC) UpdateRole(ctx context.Context, id string, rolename string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Role.Update")
	defer span.End()
	err := r.repo.UpdateRole(ctx, id, rolename)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}
func (r *RBAC) DeleteRole(ctx context.Context, id string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Role.Delete")
	defer span.End()
	err := r.repo.DeleteRole(ctx, id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}
