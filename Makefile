docker:
	docker compose up --build
inventory-gen:
	 protoc --go_out=proto/gen --go-grpc_out=proto/gen proto/inventory/inventory_service.proto
orders-gen:
	protoc --go_out=proto/gen --go-grpc_out=proto/gen proto/orders/orders_service.proto
proto-gen: inventory-gen orders-gen