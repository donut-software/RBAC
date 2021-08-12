package postgresql

import (
	"context"
	"fmt"
	"rbac/internal"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (s *Store) CreateRoleTasks(ctx context.Context, taskid string, roleid string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "RoleTask.Create")
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()
	err := s.execTx(ctx, func(q *Queries) error {
		tid, err := uuid.Parse(taskid)
		if err != nil {
			fmt.Println(err)
			return err
		}
		rid, err := uuid.Parse(roleid)
		if err != nil {
			fmt.Println(err)
			return err
		}
		_, err = q.InsertRoleTask(ctx, InsertRoleTaskParams{
			RoleID: rid,
			TaskID: tid,
		})
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	return err
}
func (s *Store) RoleTask(ctx context.Context, roleTaskId string) (internal.RoleTasks, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "RoleTask.RoleTask")
	span.SetAttributes(attribute.String("db.system", "postgresql"))
	defer span.End()
	roletask := internal.RoleTasks{}
	err := s.execTx(ctx, func(q *Queries) error {
		rtId, err := uuid.Parse(roleTaskId)
		if err != nil {
			fmt.Println(err)
		}
		rt, err := q.SelectRoleTask(ctx, rtId)
		if err != nil {
			fmt.Println(err)
			return err
		}
		roletask.Id = rt.ID.String()

		t, err := q.SelectTask(ctx, rt.TaskID)
		if err != nil {
			fmt.Println(err)
		}
		ht, err := q.SelectHelpTextByTasks(ctx, rt.TaskID)
		if err != nil {
			fmt.Println(err)
			return err
		}
		helptext := internal.HelpText{
			Id:        ht.ID.String(),
			Task_id:   ht.TaskID.String(),
			HelpText:  ht.Helptext,
			CreatedAt: ht.CreatedAt,
		}

		m, err := q.SelectMenuByTask(ctx, rt.TaskID)
		if err != nil {
			fmt.Println(err)
			return err
		}
		menu := []internal.Menu{}
		for _, value := range m {
			menu = append(menu, internal.Menu{
				Id:        value.ID.String(),
				Name:      value.Name,
				Task_id:   value.TaskID.String(),
				CreatedAt: value.CreatedAt,
			})
		}

		n, err := q.SelectNavigationByTask(ctx, rt.TaskID)
		if err != nil {
			fmt.Println(err)
			return err
		}
		nav := []internal.Navigation{}
		for _, value := range n {
			nav = append(nav, internal.Navigation{
				Id:        value.ID.String(),
				Name:      value.Name,
				Task_id:   value.TaskID.String(),
				CreatedAt: value.CreatedAt,
			})
		}
		task := internal.Tasks{
			Id:         t.ID.String(),
			Task:       t.Task,
			HelpText:   helptext,
			Menu:       menu,
			Navigation: nav,
			CreatedAt:  t.CreatedAt,
		}
		roletask.Task = task

		r, err := q.SelectRole(ctx, rt.RoleID)
		if err != nil {
			fmt.Println(err)
			return err
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
			fmt.Println(err)
		}
		rid, err := uuid.Parse(roleId)
		if err != nil {
			fmt.Println(err)
		}
		rtId, err := uuid.Parse(id)
		if err != nil {
			fmt.Println(err)
		}
		err = q.UpdateRoleTask(ctx, UpdateRoleTaskParams{
			TaskID: tid,
			RoleID: rid,
			ID:     rtId,
		})
		if err != nil {
			fmt.Println(err)
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
			fmt.Println(err)
		}
		err = q.DeleteRoleTask(ctx, rtId)
		if err != nil {
			fmt.Println(err)
		}
		return nil
	})
	return err
}
