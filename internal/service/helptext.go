package service

import (
	"context"
	"fmt"
	"rbac/internal"
)

func (r *RBAC) CreateHelpText(ctx context.Context, helptext internal.HelpText) error {
	err := r.repo.CreateHelpText(ctx, helptext)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
func (r *RBAC) HelpText(ctx context.Context, id string) (internal.HelpText, error) {
	helptext, err := r.repo.HelpText(ctx, id)
	if err != nil {
		fmt.Println(err)
		return internal.HelpText{}, err
	}
	return helptext, err
}
func (r *RBAC) UpdateHelpText(ctx context.Context, helptext internal.HelpText) error {
	err := r.repo.UpdateHelpText(ctx, helptext)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}
func (r *RBAC) DeleteHelpText(ctx context.Context, id string) error {
	err := r.repo.DeleteHelpText(ctx, id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}
