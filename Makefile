PROTO_DIR=api/proto
OUT_DIR=internal/generated
AUTH_DNS=postgres://auth_user:auth_pass@localhost:5432/auth_db?sslmode=disable

generate:
	protoc --go_out=$(OUT_DIR) --go-grpc_out=$(OUT_DIR) \
		--go_opt=paths=source_relative --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/*.proto

migrate_auth:
	goose -dir migrations/auth postgres "$(AUTH_DNS)" up