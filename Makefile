PROTO_DIR=api/proto
OUT_DIR=internal/generated

generate:
	protoc --go_out=$(OUT_DIR) --go-grpc_out=$(OUT_DIR) \
		--go_opt=paths=source_relative --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/*.proto
