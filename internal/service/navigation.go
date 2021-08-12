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
	err := r.repo.CreateNavigation(ctx, navigation)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
func (r *RBAC) Navigation(ctx context.Context, id string) (internal.Navigation, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Navigation.Navigation")
	defer span.End()
	navigation, err := r.repo.Navigation(ctx, id)
	if err != nil {
		fmt.Println(err)
		return internal.Navigation{}, err
	}
	return navigation, err
}
func (r *RBAC) UpdateNavigation(ctx context.Context, navigation internal.Navigation) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Navigation.Update")
	defer span.End()
	err := r.repo.UpdateNavigation(ctx, navigation)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
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
