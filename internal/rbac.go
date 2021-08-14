package internal

import "time"

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
