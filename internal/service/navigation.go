package service

import (
	"context"
	"fmt"
	"rbac/internal"

	"go.opentelemetry.io/otel/trace"
)

func (r *RBAC) CreateNavigation(ctx context.Context, navigation internal.Navigation) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Navigation.Create")
	defer span.End()
	id, err := r.repo.CreateNavigation(ctx, navigation)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}
	n, err := r.repo.Navigation(ctx, id)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}
	err = r.search.IndexNavigation(ctx, n)
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}
	return nil
}
func (r *RBAC) Navigation(ctx context.Context, id string) (internal.Navigation, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Navigation.Navigation")
	defer span.End()
	// navigation, err := r.repo.Navigation(ctx, id)
	navigation, err := r.search.GetNavigation(ctx, id)
	if err != nil {
		fmt.Println(err)
		return internal.Navigation{}, fmt.Errorf("search: %w", err)
	}
	return navigation, err
}
func (r *RBAC) UpdateNavigation(ctx context.Context, navigation internal.Navigation) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Navigation.Update")
	defer span.End()
	err := r.repo.UpdateNavigation(ctx, navigation)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}
	n, err := r.repo.Navigation(ctx, navigation.Id)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}
	err = r.search.DeleteNavigation(ctx, n.Id)
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}
	err = r.search.IndexNavigation(ctx, n)
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}
	return err
}

func (r *RBAC) ListNavigation(ctx context.Context, args internal.ListArgs) (internal.ListNavigation, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Navigation.List")
	defer span.End()
	lr, err := r.search.ListNavigation(ctx, args)
	if err != nil {
		return internal.ListNavigation{}, fmt.Errorf("search: %w", err)
	}
	return lr, nil

}
func (r *RBAC) DeleteNavigation(ctx context.Context, id string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Navigation.Delete")
	defer span.End()
	err := r.repo.DeleteNavigation(ctx, id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}
