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
	id, err := r.repo.CreateMenu(ctx, menu)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}
	m, err := r.repo.Menu(ctx, id)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}
	err = r.search.IndexMenu(ctx, m)
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}
	return nil
}
func (r *RBAC) Menu(ctx context.Context, id string) (internal.Menu, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Menu.Menu")
	defer span.End()
	// menu, err := r.repo.Menu(ctx, id)
	menu, err := r.search.GetMenu(ctx, id)
	if err != nil {
		return internal.Menu{}, fmt.Errorf("search: %w", err)
	}
	return menu, err
}
func (r *RBAC) UpdateMenu(ctx context.Context, menu internal.Menu) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Menu.Update")
	defer span.End()
	err := r.repo.UpdateMenu(ctx, menu)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}
	m, err := r.repo.Menu(ctx, menu.Id)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}
	err = r.search.DeleteMenu(ctx, m.Id)
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}
	err = r.search.IndexMenu(ctx, m)
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}
	return err
}

func (r *RBAC) ListMenu(ctx context.Context, args internal.ListArgs) (internal.ListMenu, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Task.List")
	defer span.End()
	lr, err := r.search.ListMenu(ctx, args)
	if err != nil {
		return internal.ListMenu{}, fmt.Errorf("search: %w", err)
	}
	return lr, nil

}

func (r *RBAC) DeleteMenu(ctx context.Context, id string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Menu.Delete")
	defer span.End()
	err := r.repo.DeleteMenu(ctx, id)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}
	err = r.search.DeleteMenu(ctx, id)
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}
	return err
}
