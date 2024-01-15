FROM golang:1.21 AS builder
LABEL maintainer="Rasbora <[rasbora.support]@openseawaves.com>"
WORKDIR /build

# Fetch dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build phase
COPY . ./
RUN go mod verify
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags "-s -X main.Version=`git describe --tags --long` -X main.BuildTime=`TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ'` -X main.GitHash=`git rev-parse HEAD`" -o linux_amd64_`git describe --tags --long` -o rasbora ./cmd/main.go

# Production phase
FROM golang:1.21
RUN apt-get update && apt-get install ffmpeg -y
WORKDIR /app
COPY --from=builder /build/rasbora .
EXPOSE 3701
ENTRYPOINT [ "/app/rasbora"]