package postgresql

import (
	"context"
	"rbac/internal"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (s *Store) CreateTask(ctx context.Context, taskname string) (string, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Task.Create")
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()
	var tid string
	err := s.execTx(ctx, func(q *Queries) error {
		id, err := q.InsertTask(ctx, taskname)
		if err != nil {
			return handleError(err, "create task", internal.ErrorCodeUnknown, "")
		}
		tid = id.String()
		return nil
	})
	return tid, err
}
func (s *Store) Task(ctx context.Context, id string) (internal.Tasks, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Task.Task")
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()
	tasks := internal.Tasks{}
	err := s.execTx(ctx, func(q *Queries) error {
		tid, err := uuid.Parse(id)
		if err != nil {
			return handleError(err, "parse id", internal.ErrorCodeInvalidArgument, "")
		}
		t, err := q.SelectTask(ctx, tid)
		if err != nil {
			return handleError(err, "get task", internal.ErrorCodeUnknown, "task not found")
		}
		tasks.Id = t.ID.String()
		tasks.Task = t.Task
		tasks.CreatedAt = t.CreatedAt
		return nil
	})
	return tasks, err
}
func (s *Store) UpdateTask(ctx context.Context, id string, taskname string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Task.Update")
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()
	err := s.execTx(ctx, func(q *Queries) error {
		tid, err := uuid.Parse(id)
		if err != nil {
			return handleError(err, "parse id", internal.ErrorCodeInvalidArgument, "")
		}
		err = q.UpdateTask(ctx, UpdateTaskParams{
			Task: taskname,
			ID:   tid,
		})
		if err != nil {
			return handleError(err, "update task", internal.ErrorCodeUnknown, "")
		}
		return nil
	})
	return err
}

func (s *Store) DeleteTask(ctx context.Context, id string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Task.Delete")
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()
	err := s.execTx(ctx, func(q *Queries) error {
		tid, err := uuid.Parse(id)
		if err != nil {
			return handleError(err, "parse id", internal.ErrorCodeInvalidArgument, "")
		}
		err = q.DeleteRoleTaskByTask(ctx, tid)
		if err != nil {
			return handleError(err, "delete role by task", internal.ErrorCodeUnknown, "")
		}
		err = q.DeleteHelpTextByTask(ctx, tid)
		if err != nil {
			return handleError(err, "delete helptext by task", internal.ErrorCodeUnknown, "")
		}
		err = q.DeleteMenuByTask(ctx, tid)
		if err != nil {
			return handleError(err, "delete menu by task", internal.ErrorCodeUnknown, "")
		}
		err = q.DeleteNavigationByTask(ctx, tid)
		if err != nil {
			return handleError(err, "delete navigation by task", internal.ErrorCodeUnknown, "")
		}
		err = q.DeleteTask(ctx, tid)
		if err != nil {
			return handleError(err, "delete task", internal.ErrorCodeUnknown, "")
		}
		return nil
	})
	return err
}
