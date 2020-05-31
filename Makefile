
lint:
	golangci-lint run --fix --verbose --concurrency 4 --timeout 5m

codegen:
	protoc --proto_path=./api/v1 --go_out=plugins=grpc:./api/v1 api/v1/apis.proto

