package postgresql

import (
	"context"
	"rbac/internal"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (s *Store) CreateAccountRole(ctx context.Context, accountId string, roleId string) (string, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "AccountRole.Create")
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()
	var arid string
	err := s.execTx(ctx, func(q *Queries) error {
		aid, err := uuid.Parse(accountId)
		if err != nil {
			return handleError(err, "parse account id", internal.ErrorCodeInvalidArgument, "")
		}
		rid, err := uuid.Parse(roleId)
		if err != nil {
			return handleError(err, "parse role id", internal.ErrorCodeInvalidArgument, "")
		}
		id, err := q.InsertAccountRole(ctx, InsertAccountRoleParams{
			AccountID: aid,
			RoleID:    rid,
		})
		if err != nil {
			return handleError(err, "create account role", internal.ErrorCodeUnknown, "")
		}
		arid = id.String()
		return nil
	})
	return arid, err
}
func (s *Store) AccountRole(ctx context.Context, accountRoleId string) (internal.AccountRoles, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "AccountRole.AcountRole")
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()
	accountrole := internal.AccountRoles{}
	err := s.execTx(ctx, func(q *Queries) error {
		acid, err := uuid.Parse(accountRoleId)
		if err != nil {
			return handleError(err, "parse account role id", internal.ErrorCodeInvalidArgument, "")
		}
		ar, err := q.SelectAccountRole(ctx, acid)
		if err != nil {
			return handleError(err, "get account role", internal.ErrorCodeUnknown, "accounr role not found")
		}
		accountrole.Id = ar.ID.String()
		accountrole.CreatedAt = ar.CreatedAt

		acc, err := q.SelectAccountsById(ctx, ar.AccountID)
		if err != nil {
			return handleError(err, "get account", internal.ErrorCodeUnknown, "account not found")
		}
		prof, err := q.SelectProfile(ctx, acc.Profile)
		if err != nil {
			return handleError(err, "get profile", internal.ErrorCodeUnknown, "profile not found")
		}
		profile := internal.Profile{
			Id:                 prof.ID.String(),
			Profile_Picture:    prof.ProfilePicture,
			Profile_Background: prof.ProfileBackground,
			First_Name:         prof.FirstName,
			Last_Name:          prof.LastName,
			Mobile:             prof.Mobile,
			Email:              prof.Email,
			CreatedAt:          prof.CreatedAt,
		}
		account := internal.Account{
			Id:        acc.ID.String(),
			UserName:  acc.Username,
			Profile:   profile,
			IsBlocked: acc.IsBlocked,
			CreatedAt: acc.CreatedAt,
		}
		accountrole.Account = account

		rl, err := q.SelectRole(ctx, ar.RoleID)
		if err != nil {
			return handleError(err, "get role", internal.ErrorCodeUnknown, "role not found")
		}
		role := internal.Roles{
			Id:        rl.ID.String(),
			Role:      rl.Role,
			CreatedAt: rl.CreatedAt,
		}
		accountrole.Role = role
		return nil
	})
	return accountrole, err
}
func (s *Store) UpdateAccountRole(ctx context.Context, accountId string, roleId string, id string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "AccountRole.Update")
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()
	err := s.execTx(ctx, func(q *Queries) error {
		acid, err := uuid.Parse(accountId)
		if err != nil {
			return handleError(err, "parse account id", internal.ErrorCodeInvalidArgument, "")
		}
		rid, err := uuid.Parse(roleId)
		if err != nil {
			return handleError(err, "parse role id", internal.ErrorCodeInvalidArgument, "")
		}
		err = q.UpdateAccountRole(ctx, UpdateAccountRoleParams{
			AccountID: acid,
			RoleID:    rid,
			ID:        acid,
		})
		if err != nil {
			return handleError(err, "update account role", internal.ErrorCodeUnknown, "")
		}
		return nil
	})
	return err
}

func (s *Store) DeleteAccountRole(ctx context.Context, id string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "AccountRole.Delete")
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()
	err := s.execTx(ctx, func(q *Queries) error {
		arid, err := uuid.Parse(id)
		if err != nil {
			return handleError(err, "parse id", internal.ErrorCodeInvalidArgument, "")
		}
		err = q.DeleteAccountRole(ctx, arid)
		if err != nil {
			return handleError(err, "delete accountrole", internal.ErrorCodeUnknown, "")
		}
		return nil
	})
	return err
}
