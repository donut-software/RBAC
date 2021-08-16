package internal

import (
	"fmt"
	"rbac/internal/envvar"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

func NewMemcached(conf *envvar.Configuration) (*memcache.Client, error) {
	host, err := conf.Get("MEMCACHED_HOST")
	if err != nil {
		return nil, fmt.Errorf("conf.Get MEMCACHED_HOST %w", err)
	}

	// XXX Assuming environment variable contains only one server
	client := memcache.New(host)

	if err := client.Ping(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	client.Timeout = 100 * time.Millisecond
	// client.Timeout = 1 * time.Second
	client.MaxIdleConns = 100

	return client, nil
}
