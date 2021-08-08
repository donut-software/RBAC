migrateup:
	migrate -path db/migrations/ -database postgres://user:user@localhost:5432/rbac?sslmode=disable up
migratedown:
	migrate -path db/migrations/ -database postgres://user:user@localhost:5432/rbac?sslmode=disable down
Phony: migrateup migratedown