user:
	@go build -o bin/user ./user-management
	@./bin/user

budget:
	@go run budget/main.go

seed:
	@go run scripts/seed.go
	
test:
	go test -v ./...

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative types/user.proto

.PHONY: user budget