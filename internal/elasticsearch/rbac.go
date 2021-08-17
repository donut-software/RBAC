package elasticsearch

import (
	esv7 "github.com/elastic/go-elasticsearch/v7"
)

type RBAC struct {
	client     *esv7.Client
	searchSize *int
}

func NewRBAC(client *esv7.Client, searchSize int) *RBAC {
	return &RBAC{
		client:     client,
		searchSize: &searchSize,
	}
}

const (
	INDEX_ACCOUNT      = "rbacaccount"
	INDEX_PROFILE      = "rbacprofile"
	INDEX_ROLE         = "rbacrole"
	INDEX_ACCOUNT_ROLE = "rbacaccountrole"
	INDEX_TASK         = "rbactask"
	INDEX_ROLE_TASK    = "rbacroletask"
	INDEX_HELPTEXT     = "rbachelptext"
	INDEX_MENU         = "rbacmenu"
	INDEX_NAVIGATION   = "rbacnavigation"
)
