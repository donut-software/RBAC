package internal

import "time"

const (
	//default tasks
	CREATE_ACCOUNT = "create account"
	GET_ACCOUNT    = "get account"
	UPDATE_ACCOUNT = "update account"
	DELETE_ACCOUNT = "delete account"
	LIST_ACCOUNT   = "list account"

	CREATE_ROLE = "create role"
	GET_ROLE    = "get role"
	UPDATE_ROLE = "update role"
	DELETE_ROLE = "delete role"
	LIST_ROLE   = "list role"

	CREATE_ACCOUNT_ROLE = "create account role"
	GET_ACCOUNT_ROLE    = "get account role"
	UPDATE_ACCOUNT_ROLE = "update account role"
	DELETE_ACCOUNT_ROLE = "delete account role"
	LIST_ACCOUNT_ROLE   = "list account role"

	CREATE_TASK = "create task"
	GET_TASK    = "get task"
	UPDATE_TASK = "update task"
	DELETE_TASK = "delete task"
	LIST_TASK   = "list task"

	CREATE_ROLE_TASK = "create role task"
	GET_ROLE_TASK    = "get role task"
	UPDATE_ROLE_TASK = "update role task"
	DELETE_ROLE_TASK = "delete role task"
	LIST_ROLE_TASK   = "list role task"

	//events
	EVENT_ACCOUNT_CREATED = "rbac.accounts.event.created"
	EVENT_ACCOUNT_UPDATED = "rbac.accounts.event.updated"
	EVENT_ACCOUNT_DELETED = "rbac.accounts.event.deleted"

	EVENT_PROFILE_CREATED = "rbac.profiles.event.created"
	EVENT_PROFILE_UPDATED = "rbac.profiles.event.updated"
	EVENT_PROFILE_DELETED = "rbac.profiles.event.deleted"

	EVENT_TASK_CREATED = "rbac.tasks.event.created"
	EVENT_TASK_UPDATED = "rbac.tasks.event.updated"
	EVENT_TASK_DELETED = "rbac.tasks.event.deleted"

	EVENT_ROLE_CREATED = "rbac.roles.event.created"
	EVENT_ROLE_UPDATED = "rbac.roles.event.updated"
	EVENT_ROLE_DELETED = "rbac.roles.event.deleted"

	EVENT_ACCOUNTROLE_CREATED = "rbac.accountRole.event.created"
	EVENT_ACCOUNTROLE_UPDATED = "rbac.accountRole.event.updated"
	EVENT_ACCOUNTROLE_DELETED = "rbac.accountRole.event.deleted"

	EVENT_ROLETASK_CREATED = "rbac.roleTasks.event.created"
	EVENT_ROLETASK_UPDATED = "rbac.roleTasks.event.updated"
	EVENT_ROLETASK_DELETED = "rbac.roleTasks.event.deleted"
)

type Profile struct {
	Id                 string
	Profile_Picture    string
	Profile_Background string
	First_Name         string
	Last_Name          string
	Mobile             string
	Email              string
	CreatedAt          time.Time
}

func (p *Profile) Validate() error {
	// Todo
	return nil
}

type Account struct {
	Id             string
	UserName       string
	HashedPassword string
	Profile        Profile
	IsBlocked      bool
	CreatedAt      time.Time
}

func (a *Account) Validate() error {
	// Todo
	return nil
}

type Roles struct {
	Id        string
	Role      string
	CreatedAt time.Time
}

func (r *Roles) Validate() error {
	// Todo
	return nil
}

type AccountRoles struct {
	Id        string
	Account   Account
	Role      Roles
	CreatedAt time.Time
}

func (ar *AccountRoles) Validate() error {
	// Todo
	return nil
}

type Tasks struct {
	Id         string
	Task       string
	HelpText   HelpText
	Menu       []Menu
	Navigation []Navigation
	CreatedAt  time.Time
}

func (t *Tasks) Validate() error {
	// Todo
	return nil
}

type RoleTasks struct {
	Id        string
	Task      Tasks
	Role      Roles
	CreatedAt time.Time
}

func (rt *RoleTasks) Validate() error {
	// Todo
	return nil
}

type HelpText struct {
	Id        string
	HelpText  string
	Task_id   string
	CreatedAt time.Time
}

func (h *HelpText) Validate() error {
	// Todo
	return nil
}

type Menu struct {
	Id        string
	Name      string
	Task_id   string
	CreatedAt time.Time
}

func (m *Menu) Validate() error {
	// Todo
	return nil
}

type Navigation struct {
	Id        string
	Name      string
	Task_id   string
	CreatedAt time.Time
}

func (n *Navigation) Validate() error {
	// Todo
	return nil
}
