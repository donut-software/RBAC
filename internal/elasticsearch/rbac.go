package elasticsearch

import (
	esv7 "github.com/elastic/go-elasticsearch/v7"
)

type RBAC struct {
	client *esv7.Client
}

func NewRBAC(client *esv7.Client) *RBAC {
	return &RBAC{
		client: client,
	}
}

const (
	INDEX_ACCOUNT = "rbac_account"
	INDEX_PROFILE = "rbac_profile"
)
