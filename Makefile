migrateup:
	migrate -path db/migrations/ -database postgres://user:user@localhost:5432/rbac?sslmode=disable up
migratedown:
	migrate -path db/migrations/ -database postgres://user:user@localhost:5432/rbac?sslmode=disable down
rest-server:
	go run cmd/rest-server/main.go --env env.example
indexer-redis:
	go run cmd/elasticsearch-indexer-redis/main.go --env env.example
seed:
	go run cmd/seeder/main.go --env env.example
docker-migrateup:
	docker-compose run rest-server migrate -path /api/migrations/ -database postgres://user:user@postgres:5432/rbac?sslmode=disable up
docker-seed:
	docker-compose run rest-server seeder --env /api/env.example