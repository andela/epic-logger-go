
all:
	POD_NAME=golang-service-123-456 go test -v

run:
	POD_NAME=golang-service-123-456 go run example/main.go