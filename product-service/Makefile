
migrate_import:
  go get github.com/golang-migrate/migrate/v4/database/postgres

migrate_create:
  migrate create -ext sql -dir migrations -seq order

migrate_up:
  migrate -database postgres://postgres:1915@localhost:5432/ecommerce?sslmode=disable -path ./migrations up

migrate_down:
  migrate -database postgres://postgres:1915@localhost:5432/ecommerce?sslmode=disable -path ./migrations down

migrate_force:
  migrate -database postgres://postgres:1915@localhost:5432/ecommerce?sslmode=disable -path ./migrations force

client_proto:
  protoc --go_out=. --go-grpc_out=. protos/clientpb/client.proto
