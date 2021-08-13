package internal

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
