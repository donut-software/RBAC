package postgresql

import (
	"context"
	"fmt"
	"rbac/internal"

	"github.com/google/uuid"
)

func (s *Store) CreateRole(ctx context.Context, rolename string) error {
	err := s.execTx(ctx, func(q *Queries) error {
		_, err := q.InsertRole(ctx, rolename)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	return err
}

func (s *Store) Role(ctx context.Context, id string) (internal.Roles, error) {
	role := internal.Roles{}
	err := s.execTx(ctx, func(q *Queries) error {
		rid, err := uuid.Parse(id)
		if err != nil {
			fmt.Println(err)
			return err
		}
		r, err := q.SelectRole(ctx, rid)
		if err != nil {
			fmt.Println(err)
			return err
		}
		role.Id = r.ID.String()
		role.Role = r.Role
		role.CreatedAt = r.CreatedAt
		return nil
	})
	return role, err
}

func (s *Store) UpdateRole(ctx context.Context, id string, rolename string) error {

	err := s.execTx(ctx, func(q *Queries) error {
		rid, err := uuid.Parse(id)
		if err != nil {
			fmt.Println(err)
			return err
		}
		err = q.UpdateRole(ctx, UpdateRoleParams{
			Role: rolename,
			ID:   rid,
		})
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	return err
}

func (s *Store) DeleteRole(ctx context.Context, id string) error {
	err := s.execTx(ctx, func(q *Queries) error {
		rid, err := uuid.Parse(id)
		if err != nil {
			fmt.Println(err)
			return err
		}
		err = q.DeleteRole(ctx, rid)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	return err
}
