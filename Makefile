runw: main.go
	./gmsess.exe

run: main.go
	go build

	./gmsess

build-proto: ./proto
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./proto/session.proto