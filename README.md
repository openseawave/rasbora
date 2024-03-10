# Rasbora

Distributed Scalable Video Transcoding Cluster with Hardware Acceleration for Universal Playback - An affordable on-premise or self-hosted alternative to cloud transcoders like Amazon Elastic Transcoder, Google Transcoder API, Wowza, and Azure Media Services.

## Components

The Rasboa philosophy is built on the idea of components, meaning that each part of the system can operate independently or be combined with a few others.

| Name              | Status     | Responsibility                                           |
|-------------------|------------|----------------------------------------------------------|
| VideoTranscoding  | ✅ Ready   | Encode/decode and transcode video files.                 |
| TaskManagement   | ✅ Ready   | Receive and manage tasks through the API protocol.        |
| CallbackManager   | ✅ Ready   | Manage callbacks when tasks are ready or encounter issues.|
| SystemRadar       | ✅ Ready   | Monitor and collect system status information.           |
| HouseKeeper       | ⬜️ In Progress | Work as a supervisor to ensure tasks in the transcoder do not stack indefinitely, addressing issues in case of errors. |

## Prerequisites

Here is what you need to be able to run Rasbora.

    Git (Version: >=2.x.x)
    Docker Compose (Version: >=2.24.x)
    Docker (Version: >=25.x.x)

## Development (Standalone)

Deploy Rasbora as standalone:

```bash
git clone https://github.com/openseawave/rasbora.git

cd rasbora/

# Modify the .env file to fit your system
cp .env.example .env

# Run Docker Compose
docker-compose up -d
```

Note: Please ensure that Docker Compose is installed on your system before running the above commands.

## Supported Storage System

| Method     | Type | Supported |Status|
|--------------|-----------|-------|----|
| SSD/HHD     | LocalStorage |✅  Ready       |Done  |
| S3     | ObjectStorage |✅  Ready        |Done  |
| Gluster    | Network | ⬜️ | In Progress |
| FreeNAS | Network| ⬜️ | In Progress |

## Supported Transcoder Engines

| Engine     | Supported |Status|
|--------------|-----------|-------|
| ffmpeg     | ✅  Ready        |Done  |

## Supported Hardware Acceleration

| Hardware     |Engine |Supported| Image | Handler | Status |
|--------------|-------|---------|-------|---------|--------|
| CPU/x86-64| ffmpeg | ✅  Ready|jrottenberg/ffmpeg:4.4-alpine | src/videotranscoder/handlers/default.handler | Done |
| CPU/ARM64 | ffmpeg | ⬜️ | In Progress| In Progress | In Progress |
| GPU/Apple Silicon| ffmpeg | ⬜️ | In Progress | In Progress | In Progress |
| GPU/Nvidia| ffmpeg | ⬜️ | In Progress | In Progress | In Progress |
| GPU/AMD| ffmpeg| ⬜️ | In Progress |  In Progress | In Progress|

## Supported API Communications

| Method     | Type | Supported |Status|
|--------------|-----------|-------|----|
| Restful     | ✅  Ready| application/json        |Done  |
| gRPC         | ⬜️        | application/protobuf |In Progress  |
| Websocket    | ⬜️        | application/json |In Progress  |

## Supported Callback Methods

| Protocol     | Supported | Type |Status|
|--------------|-----------|------|-------|
| HTTP/1.1     | ✅  Ready      | application/json |Done  |
| gRPC         | ⬜️        | application/protobuf |In Progress |
| Websocket    | ⬜️        | application/json |In Progress  |

## Supported Queue/Database Systems

| Engine     | Supported |Status|
|--------------|-----------|-------|
| Redis     | ✅  Ready       |Done  |
| RabbitMQ  | ⬜️        |In Progress|
| ActiveMQ  | ⬜️ |In Progress|
| ZeroMq | ⬜️ |In Progress|
| Amazon SQS| ⬜️ |In Progress|

## Supported Log Collectors

| Engine     | Supported |Status|
|--------------|-----------|-------|
| STDOUT     | ✅   Ready      |Done  |
| Log Files  | ✅  Ready     |Done|
| Database | ✅   Ready    |Done|
| Logstash | ⬜️ |In Progress|
| Grafana Loki| ⬜️ |In Progress|
| Logwatch | ⬜️ |In Progress|

## Supported Monitoring Methods

| Method    | Supported |Status|
|--------------|-----------|-------|
| Redis/Streams|✅   Ready      |Done  |
| Grafana   | ⬜️ |In Progress|
| Prometheus| ⬜️ |In Progress|

## TODO

⬜️ Full Testing Coverage

⬜️ Start Community Forum

⬜️ Q/A Documentation Site

## License

Copyright (c) 2022-2023 https://rasbora.openseawave.com
This file is part of Rasbora Distributed Video Transcoding
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
