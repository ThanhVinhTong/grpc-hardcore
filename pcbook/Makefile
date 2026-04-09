protoc:
	del /q pb\*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	proto/*.proto

run:
	go run main.go

test:
	go test -cover ./...

.PHONY: protoc