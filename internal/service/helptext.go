package service

import (
	"context"
	"fmt"
	"rbac/internal"

	"go.opentelemetry.io/otel/trace"
)

func (r *RBAC) CreateHelpText(ctx context.Context, helptext internal.HelpText) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "HelpText.Create")
	defer span.End()
	id, err := r.repo.CreateHelpText(ctx, helptext)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}
	ht, err := r.repo.HelpText(ctx, id)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}
	err = r.search.IndexHelpText(ctx, ht)
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}
	return nil
}
func (r *RBAC) HelpText(ctx context.Context, id string) (internal.HelpText, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "HelpText.HelpText")
	defer span.End()
	// helptext, err := r.repo.HelpText(ctx, id)
	helptext, err := r.search.GetHelpText(ctx, id)
	if err != nil {
		fmt.Println(err)
		return internal.HelpText{}, fmt.Errorf("repo: %w", err)
	}
	return helptext, err
}
func (r *RBAC) UpdateHelpText(ctx context.Context, helptext internal.HelpText) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "HelpText.Update")
	defer span.End()
	err := r.repo.UpdateHelpText(ctx, helptext)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}
	ht, err := r.repo.HelpText(ctx, helptext.Id)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}
	err = r.search.DeleteHelpText(ctx, ht.Id)
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}
	err = r.search.IndexHelpText(ctx, ht)
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}

	return err
}
func (r *RBAC) DeleteHelpText(ctx context.Context, id string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "HelpText.Delete")
	defer span.End()
	err := r.repo.DeleteHelpText(ctx, id)
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}
	return err
}

func (r *RBAC) ListhelpText(ctx context.Context, args internal.ListArgs) (internal.ListHelpText, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Task.List")
	defer span.End()
	lr, err := r.search.ListHelpText(ctx, args)
	if err != nil {
		return internal.ListHelpText{}, fmt.Errorf("search: %w", err)
	}
	return lr, nil

}
