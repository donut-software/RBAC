package postgresql

import (
	"context"
	"fmt"
	"rbac/internal"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (s *Store) CreateAccountRole(ctx context.Context, accountId string, roleId string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "AccountRole.Create")
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()
	err := s.execTx(ctx, func(q *Queries) error {
		aid, err := uuid.Parse(accountId)
		if err != nil {
			fmt.Println(err)
			return err
		}
		rid, err := uuid.Parse(roleId)
		if err != nil {
			fmt.Println(err)
			return err
		}
		_, err = q.InsertAccountRole(ctx, InsertAccountRoleParams{
			AccountID: aid,
			RoleID:    rid,
		})
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	return err
}
func (s *Store) AccountRole(ctx context.Context, accountRoleId string) (internal.AccountRoles, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "AccountRole.AcountRole")
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()
	accountrole := internal.AccountRoles{}
	err := s.execTx(ctx, func(q *Queries) error {
		acid, err := uuid.Parse(accountRoleId)
		if err != nil {
			fmt.Println(err)
			return err
		}
		ar, err := q.SelectAccountRole(ctx, acid)
		if err != nil {
			fmt.Println(err)
			return err
		}
		accountrole.Id = ar.ID.String()
		accountrole.CreatedAt = ar.CreatedAt

		acc, err := q.SelectAccountsById(ctx, ar.ID)
		if err != nil {
			fmt.Println(err)
			return err
		}
		prof, err := q.SelectProfile(ctx, acc.Profile)
		if err != nil {
			fmt.Println(err)
			return err
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
			fmt.Println(err)
			return err
		}
		role := internal.Roles{
			Id:        account.Id,
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
			fmt.Println(err)
			return err
		}
		rid, err := uuid.Parse(roleId)
		if err != nil {
			fmt.Println(err)
			return err
		}
		err = q.UpdateAccountRole(ctx, UpdateAccountRoleParams{
			AccountID: acid,
			RoleID:    rid,
			ID:        acid,
		})
		if err != nil {
			fmt.Println(err)
			return err
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
			fmt.Println(err)
			return err
		}
		err = q.DeleteAccountRole(ctx, arid)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	return err
}
