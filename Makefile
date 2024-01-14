BINARY_NAME=rasbora

run:
	go run ./cmd/main.go

build:
	GOARCH=amd64 GOOS=darwin go build -o .build/${BINARY_NAME}-darwin ./cmd/main.go
	GOARCH=amd64 GOOS=linux go build -o .build/${BINARY_NAME}-linux ./cmd/main.go
	GOARCH=amd64 GOOS=windows go build -o .build/${BINARY_NAME}-windows ./cmd/main.go

swagger:
	swag init --pd -g ./src/taskmanager/taskmanger_restful_protocol.go -o ./src/taskmanager/docs

clean:
	go clean
	rm -rf .build/
