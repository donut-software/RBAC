package events

import "rbac/internal/memcached"

type RBACEvents struct {
	cache *memcached.RBAC
}

func NewRBACEvents(cache *memcached.RBAC) *RBACEvents {
	return &RBACEvents{
		cache: cache,
	}
}
