package postgresql

import (
	"context"
	"fmt"
	"rbac/internal"

	"github.com/google/uuid"
)

func (s *Store) CreateNavigation(ctx context.Context, menu internal.Navigation) error {
	err := s.execTx(ctx, func(q *Queries) error {
		htId, err := uuid.Parse(menu.Task_id)
		if err != nil {
			fmt.Println(err)
			return err
		}
		_, err = q.InsertNavigation(ctx, InsertNavigationParams{
			Name:   menu.Name,
			TaskID: htId,
		})
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	return err
}
func (s *Store) Navigation(ctx context.Context, id string) (internal.Navigation, error) {
	menu := internal.Navigation{}
	err := s.execTx(ctx, func(q *Queries) error {
		htId, err := uuid.Parse(id)
		if err != nil {
			fmt.Println(err)
			return err
		}
		ht, err := q.SelectNavigation(ctx, htId)
		if err != nil {
			fmt.Println(err)
			return err
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
	err := s.execTx(ctx, func(q *Queries) error {
		tid, err := uuid.Parse(menu.Task_id)
		if err != nil {
			fmt.Println(err)
			return err
		}
		id, err := uuid.Parse(menu.Id)
		if err != nil {
			fmt.Println(err)
			return err
		}
		err = q.UpdateNavigation(ctx, UpdateNavigationParams{
			TaskID: tid,
			Name:   menu.Name,
			ID:     id,
		})
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	return err
}
func (s *Store) DeleteNavigation(ctx context.Context, id string) error {
	err := s.execTx(ctx, func(q *Queries) error {
		hid, err := uuid.Parse(id)
		if err != nil {
			fmt.Println(err)
			return err
		}
		err = q.DeleteNavigation(ctx, hid)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	return err
}
