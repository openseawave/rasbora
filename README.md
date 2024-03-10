# Rasbora

Distributed Scalable Video Transcoding Cluster with Hardware Acceleration for Universal Playback - An affordable on-premise or self-hosted alternative to cloud transcoders like Amazon Elastic Transcoder, Google Transcoder API, Wowza, and Azure Media Services.

## Prerequisites

Here is what you need to be able to run Rasbora.

    Docker Compose (Version: >=2.24.x)
    Docker (Version: >=25.x.x)

## Supported Storage System

| Method     | Type | Supported |Status|
|--------------|-----------|-------|
| SSD/HHD     | LocalStorage |✅         |done  |
| S3     | ObjectStorage |✅         |done  |
| Gluster    | Network | ⬜️ | in progress |
| FreeNAS | Network| ⬜️ | in progress |

## Supported Transcoder Engines

| Engine     | Supported |Status|
|--------------|-----------|-------|
| Fmmpeg     | ✅         |done  |

## Supported Hardware Acceleration

| Hardware     |Engine |Supported| Image | Handler | Status |
|--------------|-------|---------|-------|---------|--------|
| CPU/x86-64| ffmpeg | ✅ |jrottenberg/ffmpeg:4.4-alpine | src/videotranscoder/handlers/default.handler | done |
| CPU/ARM64 | ffmpeg | ⬜️ | not ready | not ready | in progress |
| GPU/Apple Silicon| ffmpeg | ⬜️ | not ready | not ready | in progress |
| GPU/Nvidia| ffmpeg | ⬜️ | not ready | not ready | in progress |
| GPU/AMD| ffmpeg| ⬜️ | not ready |  not ready | in progress|

## Supported API Communications

| Method     | Type | Supported |Status|
|--------------|-----------|-------|
| Restful     | application/json |✅         |done  |
| gRPC         | ⬜️        | application/protobuf |in progress  |
| Websocket    | ⬜️        | application/json |in progress  |

## Supported Callback Methods

| Protocol     | Supported | Type |Status|
|--------------|-----------|------|-------|
| HTTP/1.1     | ✅        | application/json |done  |
| gRPC         | ⬜️        | application/protobuf |in progress  |
| Websocket    | ⬜️        | application/json |in progress  |

## Supported Queue/Database Systems

| Engine     | Supported |Status|
|--------------|-----------|-------|
| Redis     | ✅         |done  |
| RabbitMQ  | ⬜️        |in progress|
| ActiveMQ  | ⬜️ |in progress|
| ZeroMq | ⬜️ |in progress|
| Amazon SQS| ⬜️ |in progress|

## Supported Log Collectors

| Engine     | Supported |Status|
|--------------|-----------|-------|
| STDOUT     | ✅         |done  |
| Log Files  | ✅       |done|
| Database | ✅       |done|
| Logstash | ⬜️ |in progress|
| Grafana Loki| ⬜️ |in progress|
| Logwatch | ⬜️ |in progress|

## Supported Monitoring Methods

| Method    | Supported |Status|
|--------------|-----------|-------|
| Redis/Streams|✅         |done  |
| Grafana   | ⬜️ |in progress|
| Prometheus| ⬜️ |in progress|
