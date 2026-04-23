protoc:
	del /q pb\*.go
	protoc --proto_path=proto \
	--go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	proto/*.proto
	
server:
	go run cmd/server/main.go -port=8080

client:
	go run cmd/client/main.go -addr=0.0.0.0:8080

login_admin:
	grpcurl -plaintext -d "{\"username\": \"admin\", \"password\": \"admin\"}" localhost:8080 pcbook.AuthService/Login

login_user:
	grpcurl -plaintext -d "{\"username\": \"user\", \"password\": \"user\"}" localhost:8080 pcbook.AuthService/Login

test:
	go test -cover ./...

.PHONY: protoc server client login_admin login_user test