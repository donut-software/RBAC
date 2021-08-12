package service

import (
	"context"
	"fmt"
	"rbac/internal"

	"go.opentelemetry.io/otel/trace"
)

func (r *RBAC) CreateMenu(ctx context.Context, menu internal.Menu) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Menu.Create")
	defer span.End()
	err := r.repo.CreateMenu(ctx, menu)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
func (r *RBAC) Menu(ctx context.Context, id string) (internal.Menu, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Menu.Menu")
	defer span.End()
	menu, err := r.repo.Menu(ctx, id)
	if err != nil {
		fmt.Println(err)
		return internal.Menu{}, err
	}
	return menu, err
}
func (r *RBAC) UpdateMenu(ctx context.Context, menu internal.Menu) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Menu.Update")
	defer span.End()
	err := r.repo.UpdateMenu(ctx, menu)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}
func (r *RBAC) DeleteMenu(ctx context.Context, id string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Menu.Delete")
	defer span.End()
	err := r.repo.DeleteMenu(ctx, id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}
