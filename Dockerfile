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

ARG FFMPEG_ENGINE_IMAGE=jrottenberg/ffmpeg:4.4-alpine

FROM golang:1.22 AS builder

LABEL maintainer="Rasbora <[rasbora.support]@openseawave.com>"

WORKDIR /build

# Fetch dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build phase
COPY . ./
RUN go mod verify
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags "-s -X main.Version=`git describe --tags --long` -X main.BuildTime=`TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ'` -X main.GitHash=`git rev-parse HEAD`" -o linux_amd64_`git describe --tags --long` -o rasbora ./cmd/main.go

# Production phase
FROM ${FFMPEG_ENGINE_IMAGE}

WORKDIR /
COPY --from=builder /build/rasbora .

# Expose port for task manager
EXPOSE 3701 

# Run the application
ENTRYPOINT [ "/rasbora"]