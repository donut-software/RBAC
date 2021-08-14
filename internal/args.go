package internal

type ListArgs struct {
	From *int
	Size *int
}
type ListAccount struct {
	Accounts []Account
	Total    int64
}
type ListAccountRole struct {
	AccountRoles []AccountRoles
	Total        int64
}

type ListRole struct {
	Roles []Roles
	Total int64
}
type ListTask struct {
	Task  []Tasks
	Total int64
}
type ListRoleTask struct {
	RoleTasks []RoleTasks
	Total     int64
}
type ListHelpText struct {
	HelpText []HelpText
	Total    int64
}
type ListMenu struct {
	Menu  []Menu
	Total int64
}
type ListNavigation struct {
	Navigation []Navigation
	Total      int64
}

type AccountRoleByAccountResult struct {
	Account Account
	Roles   []Roles
}
type AccountRoleByRoleResult struct {
	Role    Roles
	Account []Account
}

type RoleTaskByRole struct {
	Role  Roles
	Tasks []Tasks
}

type RoleTaskByTask struct {
	Task  Tasks
	Roles []Roles
}
type HelpTextByTask struct {
	Task     Tasks
	HelpText HelpText
}

type MenuByTask struct {
	Task Tasks
	Menu []Menu
}
type NavigationByTask struct {
	Task       Tasks
	Navigation []Navigation
}
