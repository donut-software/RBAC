package service

import (
	"context"
	"fmt"
	"rbac/internal"
)

func (r *RBAC) CreateMenu(ctx context.Context, menu internal.Menu) error {
	err := r.repo.CreateMenu(ctx, menu)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
func (r *RBAC) Menu(ctx context.Context, id string) (internal.Menu, error) {
	menu, err := r.repo.Menu(ctx, id)
	if err != nil {
		fmt.Println(err)
		return internal.Menu{}, err
	}
	return menu, err
}
func (r *RBAC) UpdateMenu(ctx context.Context, menu internal.Menu) error {
	err := r.repo.UpdateMenu(ctx, menu)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}
func (r *RBAC) DeleteMenu(ctx context.Context, id string) error {
	err := r.repo.DeleteMenu(ctx, id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}
