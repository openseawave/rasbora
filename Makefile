# Copyright (c) 2022-2023 https://rasbora.openseawave.com
#
# This file is part of Rasbora Distributed Video Transcoding
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.

BINARY_NAME=rasbora

run:
	go run ./cmd/main.go

build:
	GOARCH=amd64 GOOS=darwin go build -ldflags "-s -X main.Version=`git describe --tags --long` -X main.BuildTime=`TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ'` -X main.GitHash=`git rev-parse HEAD`" -o linux_amd64_`git describe --tags --long` -o .build/${BINARY_NAME}-darwin-arm64 ./cmd/main.go
	GOARCH=amd64 GOOS=linux go build -ldflags "-s -X main.Version=`git describe --tags --long` -X main.BuildTime=`TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ'` -X main.GitHash=`git rev-parse HEAD`" -o linux_amd64_`git describe --tags --long` -o .build/${BINARY_NAME}-linux-amd64 ./cmd/main.go
	GOARCH=amd64 GOOS=windows go build -ldflags "-s -X main.Version=`git describe --tags --long` -X main.BuildTime=`TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ'` -X main.GitHash=`git rev-parse HEAD`" -o linux_amd64_`git describe --tags --long` -o .build/${BINARY_NAME}-windows-amd64 ./cmd/main.go

swagger:
	swag init --pd -g ./src/taskmanager/taskmanger_restful_protocol.go -o ./src/taskmanager/docs

install:
	docker compose up -d

stop:
	docker compose down

restart:
	docker compose down
	docker compose up -d	

clean:
	go clean
	rm -rf .build/
