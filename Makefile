migrateup:
	migrate -path db/migrations/ -database postgres://user:user@localhost:5432/rbac?sslmode=disable up
migratedown:
	migrate -path db/migrations/ -database postgres://user:user@localhost:5432/rbac?sslmode=disable down
server:
	go run cmd/rest-server/main.go --env env.example
seed:
	go run cmd/seeder/main.go --env env.example
stop-containers:
	docker stop 8f43578e6195 b053b51c6ed8 c359a6944968 bf6c116fb007 5023cf58b3b1 df8ba6372597
start-containers:
	docker start 8f43578e6195 b053b51c6ed8 c359a6944968 bf6c116fb007 5023cf58b3b1 df8ba6372597

docker-migrateup:
	docker-compose run rest-server migrate -path /api/migrations/ -database postgres://user:user@postgres:5432/rbac?sslmode=disable up
docker-seed:
	docker-compose run rest-server seeder --env /api/env.example
Phony: migrateup migratedown