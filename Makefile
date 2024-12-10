build_inventory_proto:
	protoc --go_out=paths=source_relative:. \
    --go-grpc_out=paths=source_relative:. \
    services/proto/inventory/inventory.proto

build_order_proto:
	protoc --go_out=paths=source_relative:. \
    --go-grpc_out=paths=source_relative:. \
    services/proto/order/order.proto

build_payment_proto:
	protoc --go_out=paths=source_relative:. \
	--go-grpc_out=paths=source_relative:. \
	services/proto/payment/payment.proto

run_protos: build_inventory_proto build_order_proto build_payment_proto