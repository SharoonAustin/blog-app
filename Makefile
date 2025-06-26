GOCMD=go
GINKGO=$(GOCMD) run github.com/onsi/ginkgo/v2/ginkgo

generate-proto:
	protoc --go_out=./proto --go-grpc_out=./proto proto/blog.proto

run-unit-test:
	go test -coverprofile=coverage.out ./...

run-server: generate-proto run-unit-test
	go run cmd/server/main.go

run-client: generate-proto run-unit-test
	go run cmd/client/main.go



