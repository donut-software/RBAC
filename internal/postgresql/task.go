package postgresql

import (
	"context"
	"fmt"
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
			fmt.Println(err)
			return err
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
			fmt.Println(err)
			return err
		}
		t, err := q.SelectTask(ctx, tid)
		if err != nil {
			fmt.Println(err)
			return err
		}
		// ht, err := q.SelectHelpTextByTasks(ctx, tid)
		// if err != nil {
		// 	fmt.Println(err)
		// 	return err
		// }
		// helptext := internal.HelpText{
		// 	Id:        ht.ID.String(),
		// 	Task_id:   ht.TaskID.String(),
		// 	HelpText:  ht.Helptext,
		// 	CreatedAt: ht.CreatedAt,
		// }

		// m, err := q.SelectMenuByTask(ctx, tid)
		// if err != nil {
		// 	fmt.Println(err)
		// 	return err
		// }
		// menu := []internal.Menu{}
		// for _, value := range m {
		// 	menu = append(menu, internal.Menu{
		// 		Id:        value.ID.String(),
		// 		Name:      value.Name,
		// 		Task_id:   value.TaskID.String(),
		// 		CreatedAt: value.CreatedAt,
		// 	})
		// }

		// n, err := q.SelectNavigationByTask(ctx, tid)
		// if err != nil {
		// 	fmt.Println(err)
		// 	return err
		// }
		// nav := []internal.Navigation{}
		// for _, value := range n {
		// 	nav = append(nav, internal.Navigation{
		// 		Id:        value.ID.String(),
		// 		Name:      value.Name,
		// 		Task_id:   value.TaskID.String(),
		// 		CreatedAt: value.CreatedAt,
		// 	})
		// }
		tasks.Id = t.ID.String()
		tasks.Task = t.Task
		tasks.CreatedAt = t.CreatedAt
		// tasks.HelpText = helptext
		// tasks.Menu = menu
		// tasks.Navigation = nav
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
			fmt.Println(err)
			return err
		}
		err = q.UpdateTask(ctx, UpdateTaskParams{
			Task: taskname,
			ID:   tid,
		})
		if err != nil {
			fmt.Println(err)
			return err
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
			fmt.Println(err)
			return err
		}
		err = q.DeleteTask(ctx, tid)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	return err
}
