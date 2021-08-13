package service

import (
	"fmt"
	"rbac/internal"

	"go.opentelemetry.io/otel/trace"
	"golang.org/x/net/context"
)

func (r *RBAC) CreateAccount(ctx context.Context, account internal.Account, password string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Account.Create")
	defer span.End()
	err := r.repo.CreateAccount(ctx, account, password)
	if err != nil {
		return fmt.Errorf("repo create account: %w", err)
	}
	acc, err := r.repo.Account(ctx, account.UserName)
	if err != nil {
		return fmt.Errorf("repo create account: %w", err)
	}
	err = r.search.IndexAccount(ctx, acc)
	if err != nil {
		return fmt.Errorf("search indexed account: %w", err)
	}
	// err = r.search.IndexProfile(ctx, acc.Profile)
	// if err != nil {
	// 	return fmt.Errorf("search indexed profile: %w", err)
	// }
	return nil
}
func (r *RBAC) Account(ctx context.Context, username string) (internal.Account, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Account.Account")
	defer span.End()
	// account, err := r.repo.Account(ctx, username)
	account, err := r.search.GetAccount(ctx, username)
	if err != nil {
		return internal.Account{}, fmt.Errorf("get account: %w", err)
	}
	return account, nil
}
func (r *RBAC) UpdateProfile(ctx context.Context, profile internal.Profile) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Account.Update")
	defer span.End()
	err := r.repo.UpdateProfile(ctx, profile)
	if err != nil {
		fmt.Print(err)
		return err
	}
	return nil
}
func (r *RBAC) ChangePassword(ctx context.Context, username string, password string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Account.ChangePassword")
	defer span.End()
	err := r.repo.ChangePassword(ctx, username, password)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
func (r *RBAC) DeleteAccount(ctx context.Context, username string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Account.Delete")
	defer span.End()
	err := r.repo.DeleteAccount(ctx, username)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return nil
}
