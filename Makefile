proto_folder:
	mkdir -p ./services/proto/transaction_proto

build_proto:
	protoc --go_out=./services/proto/transaction_proto --go_opt=paths=source_relative \
    --go-grpc_out=./services/proto/transaction_proto --go-grpc_opt=paths=source_relative \
    services/proto/transaction.proto

run_proto:
	proto_folder
	build_proto