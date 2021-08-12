package postgresql

import (
	"context"
	"fmt"
	"rbac/internal"

	"github.com/google/uuid"
)

func (s *Store) CreateHelpText(ctx context.Context, helptext internal.HelpText) error {
	err := s.execTx(ctx, func(q *Queries) error {
		htId, err := uuid.Parse(helptext.Task_id)
		if err != nil {
			fmt.Println(err)
			return err
		}
		_, err = q.InsertHelpText(ctx, InsertHelpTextParams{
			Helptext: helptext.HelpText,
			TaskID:   htId,
		})
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	return err
}
func (s *Store) HelpText(ctx context.Context, id string) (internal.HelpText, error) {
	helptext := internal.HelpText{}
	err := s.execTx(ctx, func(q *Queries) error {
		htId, err := uuid.Parse(id)
		if err != nil {
			fmt.Println(err)
			return err
		}
		ht, err := q.SelectHelpText(ctx, htId)
		if err != nil {
			fmt.Println(err)
			return err
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
	err := s.execTx(ctx, func(q *Queries) error {
		tid, err := uuid.Parse(helptext.Task_id)
		if err != nil {
			fmt.Println(err)
			return err
		}
		id, err := uuid.Parse(helptext.Id)
		if err != nil {
			fmt.Println(err)
			return err
		}
		err = q.UpdateHelpText(ctx, UpdateHelpTextParams{
			TaskID:   tid,
			Helptext: helptext.HelpText,
			ID:       id,
		})
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	return err
}
func (s *Store) DeleteHelpText(ctx context.Context, id string) error {
	err := s.execTx(ctx, func(q *Queries) error {
		hid, err := uuid.Parse(id)
		if err != nil {
			fmt.Println(err)
			return err
		}
		err = q.DeleteHelpText(ctx, hid)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	return err
}
