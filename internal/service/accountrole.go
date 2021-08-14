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
	id, err := r.repo.CreateAccountRole(ctx, accountid, roleid)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}
	ar, err := r.repo.AccountRole(ctx, id)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}
	err = r.search.IndexAccountRole(ctx, ar)
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}
	return nil
}
func (r *RBAC) AccountRole(ctx context.Context, accountRoleId string) (internal.AccountRoles, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "AccountRole.AccountRole")
	defer span.End()
	// role, err := r.repo.AccountRole(ctx, accountRoleId)
	role, err := r.search.GetAccountRole(ctx, accountRoleId)
	if err != nil {
		fmt.Println(err)
		return internal.AccountRoles{}, err
	}
	return role, err
}
func (r *RBAC) AccountRoleByAccount(ctx context.Context, username string) (internal.AccountRoleByAccountResult, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "AccountRole.AccountRoleByAccount")
	defer span.End()
	// // role, err := r.repo.AccountRole(ctx, accountRoleId)
	// acc, err := r.search.GetAccountById(ctx, accountId)
	// if err != nil {
	// 	return internal.AccountRoleByAccountResult{}, fmt.Errorf("search: %w", err)
	// }
	role, err := r.search.GetAccountRoleByAccount(ctx, username)
	if err != nil {
		return internal.AccountRoleByAccountResult{}, fmt.Errorf("search: %w", err)
	}
	return role, err
}
func (r *RBAC) UpdateAccountRole(ctx context.Context, accountId string, roleId string, id string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "AccountRole.Update")
	defer span.End()
	err := r.repo.UpdateAccountRole(ctx, accountId, roleId, id)
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}
	return err
}
func (r *RBAC) DeleteAccountRole(ctx context.Context, id string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "AccountRole.Delete")
	defer span.End()
	err := r.repo.DeleteAccountRole(ctx, id)
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}
	return err
}
func (r *RBAC) ListAccountRole(ctx context.Context, args internal.ListArgs) (internal.ListAccountRole, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "AccountRole.List")
	defer span.End()
	lacr, err := r.search.ListAccountRole(ctx, args)
	if err != nil {
		return internal.ListAccountRole{}, fmt.Errorf("search: %w", err)
	}
	return lacr, err
}
