package service

import (
	"context"
	"fmt"
	"rbac/internal"

	"go.opentelemetry.io/otel/trace"
)

func (r *RBAC) CreateAccountRole(ctx context.Context, accountid string, roleid string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "AccountRole.Create")
	defer span.End()
	err := r.repo.CreateAccountRole(ctx, accountid, roleid)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
func (r *RBAC) AccountRole(ctx context.Context, accountRoleId string) (internal.AccountRoles, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "AccountRole.AccountRole")
	defer span.End()
	role, err := r.repo.AccountRole(ctx, accountRoleId)
	if err != nil {
		fmt.Println(err)
		return internal.AccountRoles{}, err
	}
	return role, err
}
func (r *RBAC) UpdateAccountRole(ctx context.Context, accountId string, roleId string, id string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "AccountRole.Update")
	defer span.End()
	err := r.repo.UpdateAccountRole(ctx, accountId, roleId, id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}
func (r *RBAC) DeleteAccountRole(ctx context.Context, id string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "AccountRole.Delete")
	defer span.End()
	err := r.repo.DeleteAccountRole(ctx, id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}
