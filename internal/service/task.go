package service

import (
	"context"
	"fmt"
	"rbac/internal"
)

func (r *RBAC) CreateTask(ctx context.Context, taskname string) error {
	err := r.repo.CreateTask(ctx, taskname)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
func (r *RBAC) Task(ctx context.Context, id string) (internal.Tasks, error) {
	role, err := r.repo.Task(ctx, id)
	if err != nil {
		fmt.Println(err)
		return internal.Tasks{}, err
	}
	return role, err
}
func (r *RBAC) UpdateTask(ctx context.Context, id string, taskname string) error {
	err := r.repo.UpdateTask(ctx, id, taskname)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}
func (r *RBAC) DeleteTask(ctx context.Context, id string) error {
	err := r.repo.DeleteTask(ctx, id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}
