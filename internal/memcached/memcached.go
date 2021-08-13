package memcached

import (
	"fmt"
	"rbac/internal"
)

func newKey(key string, args internal.ListArgs) string {
	var (
		from int
		size int
	)

	if args.From != nil {
		from = *args.From
	}
	if args.Size != nil {
		size = *args.Size
	}

	// if args.Role != nil {
	// 	role = *args.Role
	// }

	return fmt.Sprintf("%s_%d_%d", key, from, size)
}
