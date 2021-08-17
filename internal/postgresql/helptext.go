package postgresql

import (
	"context"
	"rbac/internal"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (s *Store) CreateHelpText(ctx context.Context, helptext internal.HelpText) (string, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Helptext.Create")
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()
	var htid string
	err := s.execTx(ctx, func(q *Queries) error {
		tid, err := uuid.Parse(helptext.Task_id)
		if err != nil {
			return handleError(err, "parse id", internal.ErrorCodeInvalidArgument, "")
		}
		id, err := q.InsertHelpText(ctx, InsertHelpTextParams{
			Helptext: helptext.HelpText,
			TaskID:   tid,
		})
		if err != nil {
			return handleError(err, "create help text", internal.ErrorCodeUnknown, "")
		}
		htid = id.String()
		return nil
	})
	return htid, err
}
func (s *Store) HelpText(ctx context.Context, id string) (internal.HelpText, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Helptext.HelpText")
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()
	helptext := internal.HelpText{}
	err := s.execTx(ctx, func(q *Queries) error {
		htId, err := uuid.Parse(id)
		if err != nil {
			return handleError(err, "parse id", internal.ErrorCodeInvalidArgument, "")
		}
		ht, err := q.SelectHelpText(ctx, htId)
		if err != nil {
			return handleError(err, "get helptext", internal.ErrorCodeUnknown, "helptext not found")
		}
		helptext.Id = ht.ID.String()
		helptext.HelpText = ht.Helptext
		helptext.Task_id = ht.TaskID.String()
		helptext.CreatedAt = ht.CreatedAt
		return nil
	})
	return helptext, err
}
func (s *Store) UpdateHelpText(ctx context.Context, helptext internal.HelpText) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Helptext.Update")
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()
	err := s.execTx(ctx, func(q *Queries) error {
		tid, err := uuid.Parse(helptext.Task_id)
		if err != nil {
			return handleError(err, "parse task id", internal.ErrorCodeInvalidArgument, "")
		}
		id, err := uuid.Parse(helptext.Id)
		if err != nil {
			return handleError(err, "parse helptext id", internal.ErrorCodeUnknown, "")
		}
		err = q.UpdateHelpText(ctx, UpdateHelpTextParams{
			TaskID:   tid,
			Helptext: helptext.HelpText,
			ID:       id,
		})
		if err != nil {
			return handleError(err, "update helptext", internal.ErrorCodeUnknown, "")
		}
		return nil
	})
	return err
}
func (s *Store) DeleteHelpText(ctx context.Context, id string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Helptext.Delete")
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()
	err := s.execTx(ctx, func(q *Queries) error {
		hid, err := uuid.Parse(id)
		if err != nil {
			return handleError(err, "parse id", internal.ErrorCodeInvalidArgument, "")
		}
		err = q.DeleteHelpText(ctx, hid)
		if err != nil {
			return handleError(err, "delete helptext", internal.ErrorCodeUnknown, "")
		}
		return nil
	})
	return err
}
