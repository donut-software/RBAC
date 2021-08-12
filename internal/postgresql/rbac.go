package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"rbac/internal"

	"github.com/google/uuid"
)

type RBAC interface {
	CreateAccount(ctx context.Context, account internal.Account, password string) error
	Account(ctx context.Context, username string) (internal.Account, error)
	UpdateProfile(ctx context.Context, profile internal.Profile) error
	ChangePassword(ctx context.Context, username string, password string) error
	DeleteAccount(ctx context.Context, username string) error

	CreateRole(ctx context.Context, rolename string) error
	Role(ctx context.Context, id string) (internal.Roles, error)
	UpdateRole(ctx context.Context, id string, rolename string) error
	DeleteRole(ctx context.Context, id string) error

	CreateTask(ctx context.Context, taskname string) error
	Task(ctx context.Context, id string) (internal.Tasks, error)
	UpdateTask(ctx context.Context, id string, taskname string) error
	DeleteTask(ctx context.Context, id string) error
}

func NewRBAC(db *sql.DB) RBAC {
	return &Store{
		db: db,
		q:  New(db),
	}
}
func (s *Store) CreateAccount(ctx context.Context, account internal.Account, password string) error {
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

func (s *Store) CreateTask(ctx context.Context, taskname string) error {
	err := s.execTx(ctx, func(q *Queries) error {
		_, err := q.InsertTask(ctx, taskname)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	return err
}
func (s *Store) Task(ctx context.Context, id string) (internal.Tasks, error) {
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
		ht, err := q.SelectHelpTextByTasks(ctx, tid)
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

		m, err := q.SelectMenuByTask(ctx, tid)
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

		n, err := q.SelectNavigationByTask(ctx, tid)
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
		tasks.Id = t.ID.String()
		tasks.Task = t.Task
		tasks.CreatedAt = t.CreatedAt
		tasks.HelpText = helptext
		tasks.Menu = menu
		tasks.Navigation = nav
		return nil
	})
	return tasks, err
}
func (s *Store) UpdateTask(ctx context.Context, id string, taskname string) error {
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
