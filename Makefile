user:
	@go run user-service/user_server.go

seed:
	@go run scripts/seed.go
	
gate:
	@go run api-gateway/gate.go

test:
	go test -v ./...

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	       types/user.proto types/expense.proto types/budget.proto

.PHONY: user seed gate test proto
