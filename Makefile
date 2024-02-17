user:
	@go build -o bin/user ./user-management
	@./bin/user

test:
	go test -v ./...

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative user-management/types/user.proto
