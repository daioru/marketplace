PROTO_DIR=api/proto
OUT_DIR=internal/generated
AUTH_DNS=postgres://auth_user:auth_pass@localhost:5432/auth_db?sslmode=disable

generate:
	protoc -I api/proto --go_out=. --go-grpc_out=. --grpc-gateway_out=. --grpc-gateway_opt logtostderr=true api/proto/auth.proto

migrate_auth_up:
	goose -dir migrations/auth postgres "$(AUTH_DNS)" up

migrate_auth_down:
	goose -dir migrations/auth postgres "$(AUTH_DNS)" down