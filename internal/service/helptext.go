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
	err := r.repo.CreateHelpText(ctx, helptext)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
func (r *RBAC) HelpText(ctx context.Context, id string) (internal.HelpText, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "HelpText.HelpText")
	defer span.End()
	helptext, err := r.repo.HelpText(ctx, id)
	if err != nil {
		fmt.Println(err)
		return internal.HelpText{}, err
	}
	return helptext, err
}
func (r *RBAC) UpdateHelpText(ctx context.Context, helptext internal.HelpText) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "HelpText.Update")
	defer span.End()
	err := r.repo.UpdateHelpText(ctx, helptext)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}
func (r *RBAC) DeleteHelpText(ctx context.Context, id string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "HelpText.Delete")
	defer span.End()
	err := r.repo.DeleteHelpText(ctx, id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}
