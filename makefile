PROTO_DIR = ./protos

.PHONY: proto
proto:
	protoc --proto_path=$(PROTO_DIR) --go_out=$(PROTO_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_DIR) --go-grpc_opt=paths=source_relative $(PROTO_DIR)/*.proto