package service

import (
	"fmt"
	"rbac/internal"

	"golang.org/x/net/context"
)

func (r *RBAC) CreateAccount(ctx context.Context, account internal.Account, password string) error {
	err := r.repo.CreateAccount(ctx, account, password)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
func (r *RBAC) Account(ctx context.Context, username string) (internal.Account, error) {
	account, err := r.repo.Account(ctx, username)
	if err != nil {
		fmt.Println(err)
		return internal.Account{}, err
	}
	return account, nil
}
func (r *RBAC) UpdateProfile(ctx context.Context, profile internal.Profile) error {
	err := r.repo.UpdateProfile(ctx, profile)
	if err != nil {
		fmt.Print(err)
		return err
	}
	return nil
}
func (r *RBAC) ChangePassword(ctx context.Context, username string, password string) error {
	err := r.repo.ChangePassword(ctx, username, password)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
func (r *RBAC) DeleteAccount(ctx context.Context, username string) error {
	err := r.repo.DeleteAccount(ctx, username)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return nil
}
