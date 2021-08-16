package service

import (
	"fmt"
	"rbac/internal"

	"go.opentelemetry.io/otel/trace"
	"golang.org/x/net/context"
)

func (r *RBAC) Logout(ctx context.Context) error {
	_, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Account.Logout")
	defer span.End()
	return nil
}
func (r *RBAC) Login(ctx context.Context, username string, password string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Account.Login")
	defer span.End()
	err := r.repo.Login(ctx, username, password)
	if err != nil {
		return fmt.Errorf("repo login: %w", err)
	}
	return nil
}
func (r *RBAC) IsAllowed(ctx context.Context, username string, task string) (bool, error) {
	acrole, err := r.search.GetAccountRoleByAccount(ctx, username)
	if err != nil {
		return false, fmt.Errorf("search: %w", err)
	}
	var tasks []internal.Tasks
	for _, value := range acrole.Roles {
		rt, err := r.search.GetRoleTaskByRole(ctx, value.Id)
		if err != nil {
			return false, fmt.Errorf("search: %w", err)
		}
		tasks = append(tasks, rt.Tasks...)
	}
	for _, value := range tasks {
		if value.Task == task {
			return true, nil
		}
	}
	return false, nil
}
func (r *RBAC) CreateAccount(ctx context.Context, account internal.Account, password string) (string, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Account.Create")
	defer span.End()
	id, err := r.repo.CreateAccount(ctx, account, password)
	if err != nil {
		return id, fmt.Errorf("repo create account: %w", err)
	}
	acc, err := r.repo.Account(ctx, account.UserName)
	if err != nil {
		return id, fmt.Errorf("repo create account: %w", err)
	}
	// err = r.search.IndexAccount(ctx, acc)
	// if err != nil {
	// 	return id, fmt.Errorf("search indexed account: %w", err)
	// }
	_ = r.msgBroker.AccountCreated(ctx, acc) // XXX: Ignoring errors on purpose
	// err = r.search.IndexProfile(ctx, acc.Profile)
	// if err != nil {
	// 	return id, fmt.Errorf("search indexed account: %w", err)
	// }
	_ = r.msgBroker.ProfileCreated(ctx, acc.Profile)
	return id, nil
}
func (r *RBAC) Account(ctx context.Context, username string) (internal.Account, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Account.Account")
	defer span.End()
	// account, err := r.repo.Account(ctx, username)
	account, err := r.search.GetAccount(ctx, username)
	if err != nil {
		return internal.Account{}, fmt.Errorf("get account: %w", err)
	}
	account.Profile, err = r.search.GetProfile(ctx, account.Profile.Id)
	if err != nil {
		return internal.Account{}, fmt.Errorf("get profile: %w", err)
	}
	return account, nil
}
func (r *RBAC) AccountByID(ctx context.Context, id string) (internal.Account, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Account.Account")
	defer span.End()
	// account, err := r.repo.Account(ctx, username)
	account, err := r.search.GetAccountById(ctx, id)
	if err != nil {
		return internal.Account{}, fmt.Errorf("get account: %w", err)
	}
	account.Profile, err = r.search.GetProfile(ctx, account.Profile.Id)
	if err != nil {
		return internal.Account{}, fmt.Errorf("get profile: %w", err)
	}
	return account, nil
}
func (r *RBAC) UpdateProfile(ctx context.Context, profile internal.Profile) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Account.Update")
	defer span.End()
	err := r.repo.UpdateProfile(ctx, profile)
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}
	err = r.search.DeleteProfile(ctx, profile.Id)
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}
	err = r.search.IndexProfile(ctx, profile)
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}
	return nil
}
func (r *RBAC) ChangePassword(ctx context.Context, username string, password string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Account.ChangePassword")
	defer span.End()
	err := r.repo.ChangePassword(ctx, username, password)
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}
	return nil
}
func (r *RBAC) DeleteAccount(ctx context.Context, username string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Account.Delete")
	defer span.End()
	err := r.repo.DeleteAccount(ctx, username)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}
	err = r.search.DeleteAccount(ctx, username)
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}
	return nil
}
func (r *RBAC) ListAccount(ctx context.Context, args internal.ListArgs) (internal.ListAccount, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Account.List")
	defer span.End()
	la, err := r.search.ListAccount(ctx, args)
	if err != nil {
		return internal.ListAccount{}, fmt.Errorf("search: %w", err)
	}
	return la, nil

}
