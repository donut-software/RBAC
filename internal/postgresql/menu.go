package postgresql

import (
	"context"
	"fmt"
	"rbac/internal"

	"github.com/google/uuid"
)

func (s *Store) CreateMenu(ctx context.Context, menu internal.Menu) error {
	err := s.execTx(ctx, func(q *Queries) error {
		htId, err := uuid.Parse(menu.Task_id)
		if err != nil {
			fmt.Println(err)
			return err
		}
		_, err = q.InsertMenu(ctx, InsertMenuParams{
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
func (s *Store) Menu(ctx context.Context, id string) (internal.Menu, error) {
	menu := internal.Menu{}
	err := s.execTx(ctx, func(q *Queries) error {
		htId, err := uuid.Parse(id)
		if err != nil {
			fmt.Println(err)
			return err
		}
		ht, err := q.SelectMenu(ctx, htId)
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
func (s *Store) UpdateMenu(ctx context.Context, menu internal.Menu) error {
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
		err = q.UpdateMenu(ctx, UpdateMenuParams{
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
func (s *Store) DeleteMenu(ctx context.Context, id string) error {
	err := s.execTx(ctx, func(q *Queries) error {
		hid, err := uuid.Parse(id)
		if err != nil {
			fmt.Println(err)
			return err
		}
		err = q.DeleteMenu(ctx, hid)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	return err
}
