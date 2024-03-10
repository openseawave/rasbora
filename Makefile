BINARY_NAME=rasbora

run:
	go run ./cmd/main.go

build:
	GOARCH=amd64 GOOS=darwin go build -ldflags "-s -X main.Version=`git describe --tags --long` -X main.BuildTime=`TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ'` -X main.GitHash=`git rev-parse HEAD`" -o linux_amd64_`git describe --tags --long` -o .build/${BINARY_NAME}-darwin-arm64 ./cmd/main.go
	GOARCH=amd64 GOOS=linux go build -ldflags "-s -X main.Version=`git describe --tags --long` -X main.BuildTime=`TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ'` -X main.GitHash=`git rev-parse HEAD`" -o linux_amd64_`git describe --tags --long` -o .build/${BINARY_NAME}-linux-amd64 ./cmd/main.go
	GOARCH=amd64 GOOS=windows go build -ldflags "-s -X main.Version=`git describe --tags --long` -X main.BuildTime=`TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ'` -X main.GitHash=`git rev-parse HEAD`" -o linux_amd64_`git describe --tags --long` -o .build/${BINARY_NAME}-windows-amd64 ./cmd/main.go

swagger:
	swag init --pd -g ./src/taskmanager/taskmanger_restful_protocol.go -o ./src/taskmanager/docs

clean:
	go clean
	rm -rf .build/
