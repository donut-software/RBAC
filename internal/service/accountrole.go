package service

import (
	"context"
	"fmt"
	"rbac/internal"
	"strings"

	"go.opentelemetry.io/otel/trace"
)

func (r *RBAC) CreateAccountRole(ctx context.Context, accountRole internal.AccountRoles) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "AccountRole.Create")
	defer span.End()
	id, err := r.repo.CreateAccountRole(ctx, accountRole.Account.Id, accountRole.Role.Id)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}
	ar, err := r.repo.AccountRole(ctx, id)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}
	_ = r.msgBroker.AccountRoleCreated(ctx, ar)
	return nil
}
func (r *RBAC) AccountRole(ctx context.Context, accountRoleId string) (internal.AccountRoles, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "AccountRole.AccountRole")
	defer span.End()
	accRole, err := r.search.GetAccountRole(ctx, accountRoleId)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			accRole, err = r.repo.AccountRole(ctx, accountRoleId)
			if err != nil {
				return internal.AccountRoles{}, fmt.Errorf("get account: %w", err)
			}
			//if you get here means account and profile has value and no error
			err = r.search.IndexAccountRole(ctx, accRole)
			if err != nil {
				return internal.AccountRoles{}, fmt.Errorf("index account: %w", err)
			}
			return accRole, err
		}
	}
	return accRole, err
}
func (r *RBAC) AccountRoleByAccount(ctx context.Context, username string) (internal.AccountRoleByAccountResult, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "AccountRole.AccountRoleByAccount")
	defer span.End()
	role, err := r.search.GetAccountRoleByAccount(ctx, username)
	if err != nil {
		return internal.AccountRoleByAccountResult{}, fmt.Errorf("search: %w", err)
	}
	return role, err
}
func (r *RBAC) AccountRoleByRole(ctx context.Context, id string) (internal.AccountRoleByRoleResult, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "AccountRole.AccountRoleByAccount")
	defer span.End()
	role, err := r.search.GetAccountRoleByRole(ctx, id)
	if err != nil {
		return internal.AccountRoleByRoleResult{}, fmt.Errorf("search: %w", err)
	}
	return role, err
}
func (r *RBAC) UpdateAccountRole(ctx context.Context, accountRole internal.AccountRoles) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "AccountRole.Update")
	defer span.End()
	err := r.repo.UpdateAccountRole(ctx, accountRole.Account.Id, accountRole.Role.Id, accountRole.Id)
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}
	_ = r.msgBroker.AccountRoleUpdated(ctx, accountRole)
	return err
}
func (r *RBAC) DeleteAccountRole(ctx context.Context, id string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "AccountRole.Delete")
	defer span.End()
	err := r.repo.DeleteAccountRole(ctx, id)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}
	_ = r.msgBroker.AccountRoleDeleted(ctx, id)
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
