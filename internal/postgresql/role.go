package postgresql

import (
	"context"
	"rbac/internal"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (s *Store) CreateRole(ctx context.Context, rolename string) (string, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Role.Create")
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()
	var rid string
	err := s.execTx(ctx, func(q *Queries) error {
		id, err := q.InsertRole(ctx, rolename)
		if err != nil {
			return handleError(err, "create role", internal.ErrorCodeUnknown, "")
		}
		rid = id.String()
		return nil
	})
	return rid, err
}

func (s *Store) Role(ctx context.Context, id string) (internal.Roles, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Role.Role")
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()
	role := internal.Roles{}
	err := s.execTx(ctx, func(q *Queries) error {
		rid, err := uuid.Parse(id)
		if err != nil {
			return handleError(err, "parse id", internal.ErrorCodeInvalidArgument, "")
		}
		r, err := q.SelectRole(ctx, rid)
		if err != nil {
			return handleError(err, "get role", internal.ErrorCodeInvalidArgument, "role not found")
		}
		role.Id = r.ID.String()
		role.Role = r.Role
		role.CreatedAt = r.CreatedAt
		return nil
	})
	return role, err
}

func (s *Store) UpdateRole(ctx context.Context, id string, rolename string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Role.Update")
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()

	err := s.execTx(ctx, func(q *Queries) error {
		rid, err := uuid.Parse(id)
		if err != nil {
			return handleError(err, "parse id", internal.ErrorCodeInvalidArgument, "")
		}
		err = q.UpdateRole(ctx, UpdateRoleParams{
			Role: rolename,
			ID:   rid,
		})
		if err != nil {
			return handleError(err, "update role", internal.ErrorCodeUnknown, "")
		}
		return nil
	})
	return err
}

func (s *Store) DeleteRole(ctx context.Context, id string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Role.Delete")
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()
	err := s.execTx(ctx, func(q *Queries) error {
		rid, err := uuid.Parse(id)
		if err != nil {
			return handleError(err, "parse id", internal.ErrorCodeInvalidArgument, "")
		}
		err = q.DeleteAccountRoleByRole(ctx, rid)
		if err != nil {
			return handleError(err, "delete account role by role", internal.ErrorCodeUnknown, "")
		}
		err = q.DeleteRoleTaskByRole(ctx, rid)
		if err != nil {
			return handleError(err, "delete task by role", internal.ErrorCodeUnknown, "")
		}
		err = q.DeleteRole(ctx, rid)
		if err != nil {
			return handleError(err, "delete role", internal.ErrorCodeUnknown, "")
		}
		return nil
	})
	return err
}
