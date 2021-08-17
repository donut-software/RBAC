package postgresql

import (
	"context"
	"rbac/internal"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (s *Store) CreateNavigation(ctx context.Context, menu internal.Navigation) (string, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Navigation.Create")
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()
	var nid string
	err := s.execTx(ctx, func(q *Queries) error {
		htId, err := uuid.Parse(menu.Task_id)
		if err != nil {
			return handleError(err, "parse id", internal.ErrorCodeInvalidArgument, "")
		}
		id, err := q.InsertNavigation(ctx, InsertNavigationParams{
			Name:   menu.Name,
			TaskID: htId,
		})
		if err != nil {
			return handleError(err, "create navigation", internal.ErrorCodeUnknown, "")
		}
		nid = id.String()
		return nil
	})
	return nid, err
}
func (s *Store) Navigation(ctx context.Context, id string) (internal.Navigation, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Navigation.Navigation")
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()
	menu := internal.Navigation{}
	err := s.execTx(ctx, func(q *Queries) error {
		htId, err := uuid.Parse(id)
		if err != nil {
			return handleError(err, "parse id", internal.ErrorCodeInvalidArgument, "")
		}
		ht, err := q.SelectNavigation(ctx, htId)
		if err != nil {
			return handleError(err, "get navigation", internal.ErrorCodeUnknown, "navigation not found")
		}
		menu.Id = ht.ID.String()
		menu.Name = ht.Name
		menu.Task_id = ht.TaskID.String()
		menu.CreatedAt = ht.CreatedAt
		return nil
	})
	return menu, err
}
func (s *Store) UpdateNavigation(ctx context.Context, menu internal.Navigation) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Navigation.Update")
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()
	err := s.execTx(ctx, func(q *Queries) error {
		tid, err := uuid.Parse(menu.Task_id)
		if err != nil {
			return handleError(err, "parse task id", internal.ErrorCodeInvalidArgument, "")
		}
		id, err := uuid.Parse(menu.Id)
		if err != nil {
			return handleError(err, "parse menu id", internal.ErrorCodeInvalidArgument, "")
		}
		err = q.UpdateNavigation(ctx, UpdateNavigationParams{
			TaskID: tid,
			Name:   menu.Name,
			ID:     id,
		})
		if err != nil {
			return handleError(err, "update navigation", internal.ErrorCodeUnknown, "")
		}
		return nil
	})
	return err
}
func (s *Store) DeleteNavigation(ctx context.Context, id string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Navigation.Delete")
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()
	err := s.execTx(ctx, func(q *Queries) error {
		hid, err := uuid.Parse(id)
		if err != nil {
			return handleError(err, "parse id", internal.ErrorCodeInvalidArgument, "")
		}
		err = q.DeleteNavigation(ctx, hid)
		if err != nil {
			return handleError(err, "delete navigation", internal.ErrorCodeUnknown, "")
		}
		return nil
	})
	return err
}
