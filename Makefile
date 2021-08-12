migrateup:
	migrate -path db/migrations/ -database postgres://user:user@localhost:5432/rbac?sslmode=disable up
migratedown:
	migrate -path db/migrations/ -database postgres://user:user@localhost:5432/rbac?sslmode=disable down
rest:
	go run cmd/rest-server/main.go --env env.example
Phony: migrateup migratedown