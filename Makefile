docker:
	docker compose up --build
inventory-gen:
	 protoc --go_out=proto/gen --go-grpc_out=proto/gen proto/inventory/inventory_service.proto
orders-gen:
	protoc --go_out=proto/gen --go-grpc_out=proto/gen proto/orders/orders_service.proto
promotion-gen:
	protoc --proto_path=/usr/local/include --proto_path=proto --go_out=proto/gen --go-grpc_out=proto/gen proto/promotion/promotion_service.proto
statistics-gen:
	protoc --go_out=proto/gen --go-grpc_out=proto/gen proto/statistics/statistics_service.proto

proto-gen: inventory-gen orders-gen promotion-gen