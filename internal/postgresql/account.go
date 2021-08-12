package postgresql

import (
	"context"
	"fmt"
	"rbac/internal"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (s *Store) CreateAccount(ctx context.Context, account internal.Account, password string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Account.Create")
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()

	err := s.execTx(ctx, func(q *Queries) error {
		profileId, err := q.InsertProfile(ctx, InsertProfileParams{
			ProfilePicture:    account.Profile.Profile_Picture,
			ProfileBackground: account.Profile.Profile_Background,
			LastName:          account.Profile.Last_Name,
			FirstName:         account.Profile.First_Name,
			Mobile:            account.Profile.Mobile,
			Email:             account.Profile.Email,
		})
		if err != nil {
			fmt.Println(err)
		}
		hashedPassword, err := HashPassword(password)
		if err != nil {
			fmt.Println(err)
		}
		_, err = q.InsertAccounts(ctx, InsertAccountsParams{
			Username:       account.UserName,
			Hashedpassword: hashedPassword,
			Profile:        profileId,
		})
		if err != nil {
			fmt.Println(err)
		}
		return nil
	})
	return err
}
func (s *Store) Account(ctx context.Context, username string) (internal.Account, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Account.Account")
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()
	account := internal.Account{}
	err := s.execTx(ctx, func(q *Queries) error {
		acc, err := q.SelectAccounts(ctx, username)
		if err != nil {
			fmt.Println(err)
			return err
		}
		account.Id = acc.ID.String()
		account.UserName = acc.Username
		account.IsBlocked = acc.IsBlocked
		account.CreatedAt = acc.CreatedAt
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
		account.Profile = profile
		return nil
	})
	return account, err
}
func (s *Store) UpdateProfile(ctx context.Context, profile internal.Profile) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Account.Update")
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()
	err := s.execTx(ctx, func(q *Queries) error {
		profId, err := uuid.Parse(profile.Id)
		if err != nil {
			fmt.Println(err)
			return err
		}
		q.UpdateProfile(ctx, UpdateProfileParams{
			ProfilePicture:    profile.Profile_Picture,
			ProfileBackground: profile.Profile_Background,
			FirstName:         profile.First_Name,
			LastName:          profile.Last_Name,
			Mobile:            profile.Mobile,
			Email:             profile.Email,
			ID:                profId,
		})
		return nil
	})
	return err
}
func (s *Store) DeleteAccount(ctx context.Context, username string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Account.Delete")
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()
	err := s.execTx(ctx, func(q *Queries) error {
		err := q.DeleteAccount(ctx, username)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	return err
}
func (s *Store) ChangePassword(ctx context.Context, username string, password string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Account.ChangePassword")
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()
	err := s.execTx(ctx, func(q *Queries) error {
		hashedPassword, err := HashPassword(password)
		if err != nil {
			fmt.Println(err)
			return err
		}
		err = q.ChangePassword(ctx, ChangePasswordParams{
			Hashedpassword: hashedPassword,
			Username:       username,
		})
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	return err
}
