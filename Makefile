
all:
	POD_NAME=golang-service-123-456 go test

run:
	POD_NAME=golang-service-123-456 go run mock/main.go