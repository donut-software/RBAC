package postgresql

import (
	"context"
	"rbac/internal"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (s *Store) CreateRoleTasks(ctx context.Context, taskid string, roleid string) (string, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "RoleTask.Create")
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()
	var rtid string
	err := s.execTx(ctx, func(q *Queries) error {
		tid, err := uuid.Parse(taskid)
		if err != nil {
			return handleError(err, "parse task id", internal.ErrorCodeInvalidArgument, "")
		}
		rid, err := uuid.Parse(roleid)
		if err != nil {
			return handleError(err, "parse role id", internal.ErrorCodeInvalidArgument, "")
		}
		id, err := q.InsertRoleTask(ctx, InsertRoleTaskParams{
			RoleID: rid,
			TaskID: tid,
		})
		if err != nil {
			return handleError(err, "create role task", internal.ErrorCodeUnknown, "")
		}
		rtid = id.String()
		return nil
	})
	return rtid, err
}
func (s *Store) RoleTask(ctx context.Context, roleTaskId string) (internal.RoleTasks, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "RoleTask.RoleTask")
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()
	roletask := internal.RoleTasks{}
	err := s.execTx(ctx, func(q *Queries) error {
		rtId, err := uuid.Parse(roleTaskId)
		if err != nil {
			return handleError(err, "parse id", internal.ErrorCodeInvalidArgument, "")
		}
		rt, err := q.SelectRoleTask(ctx, rtId)
		if err != nil {
			return handleError(err, "get role task", internal.ErrorCodeUnknown, "roletask not found")
		}
		roletask.Id = rt.ID.String()

		t, err := q.SelectTask(ctx, rt.TaskID)
		if err != nil {
			return handleError(err, "get task", internal.ErrorCodeUnknown, "task not found")
		}

		task := internal.Tasks{
			Id:        t.ID.String(),
			Task:      t.Task,
			CreatedAt: t.CreatedAt,
		}
		roletask.Task = task

		r, err := q.SelectRole(ctx, rt.RoleID)
		if err != nil {
			return handleError(err, "get role", internal.ErrorCodeUnknown, "role not found")
		}
		role := internal.Roles{
			Id:        r.ID.String(),
			Role:      r.Role,
			CreatedAt: r.CreatedAt,
		}
		roletask.Role = role
		return nil
	})
	return roletask, err
}
func (s *Store) UpdateRoleTask(ctx context.Context, taskId string, roleId string, id string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "RoleTask.Update")
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()
	err := s.execTx(ctx, func(q *Queries) error {
		tid, err := uuid.Parse(taskId)
		if err != nil {
			return handleError(err, "parse task id", internal.ErrorCodeInvalidArgument, "")
		}
		rid, err := uuid.Parse(roleId)
		if err != nil {
			return handleError(err, "parse role id", internal.ErrorCodeInvalidArgument, "")
		}
		rtId, err := uuid.Parse(id)
		if err != nil {
			return handleError(err, "parse id", internal.ErrorCodeInvalidArgument, "")
		}
		err = q.UpdateRoleTask(ctx, UpdateRoleTaskParams{
			TaskID: tid,
			RoleID: rid,
			ID:     rtId,
		})
		if err != nil {
			return handleError(err, "update roletask", internal.ErrorCodeUnknown, "")
		}
		return nil
	})
	return err
}
func (s *Store) DeleteRoleTask(ctx context.Context, id string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "RoleTask.Delete")
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()
	err := s.execTx(ctx, func(q *Queries) error {
		rtId, err := uuid.Parse(id)
		if err != nil {
			return handleError(err, "parse id", internal.ErrorCodeInvalidArgument, "")
		}
		err = q.DeleteRoleTask(ctx, rtId)
		if err != nil {
			return handleError(err, "delete role task", internal.ErrorCodeUnknown, "")
		}
		return nil
	})
	return err
}
