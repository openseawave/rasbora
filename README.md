<p align="center">
  <a href="https://github.com/openseawave/rasbora">
   <img src="https://github.com/openseawave/rasbora/blob/main/docs/banner.png?raw=true" alt="Rasbora Logo">
  </a>
</p>

<p align="center">
    <a href="https://rasbora.openseawave.com/"><b>Website</b></a> •
    <a href="https://discord.com/invite/xevmJcDPKH"><b>Discord</b></a> •
    <a href="https://www.linkedin.com/company/openseawave"><b>Linkedin</b></a> •
    <a href="https://twitter.com/openseawave"><b>Twitter</b></a> •
    <a href="https://www.reddit.com/r/openseawave/"><b>Reddit</b></a> •
    <a href="https://github.com/orgs/openseawave/projects/1"><b>Roadmap</b></a> •
    <a href="https://github.com/openseawave/rasbora/discussions"><b>Community</b></a>
</p>

# Rasbora

Distributed Scalable Video Transcoding Cluster with Hardware Acceleration for Universal Playback - An affordable on-premise or self-hosted alternative to cloud transcoders like Amazon Elastic Transcoder, Google Transcoder API, Wowza, and Azure Media Services.

## Components

The Rasboa philosophy is built on the idea of components, meaning that each part of the system can operate independently or be combined with a few others.

| Name              | Status     | Responsibility                                           | Edition|
|-------------------|------------|----------------------------------------------------------|--------|
| VideoTranscoding  | ✅ Ready   | Encode/decode and transcode video files.                 | Community |
| TaskManagement   | ✅ Ready   | Receive and manage tasks through the API protocol.        | Community |
| CallbackManager   | ✅ Ready   | Manage callbacks when tasks are ready or encounter issues.| Community |
| SystemRadar       | ✅ Ready   | Monitor and collect system status information.           | Community |
| HouseKeeper       | ⬜️ In Progress | Work as a supervisor to ensure tasks in the transcoder do not stack indefinitely, addressing issues in case of errors. | Enterprise |
| Dashboard | ⬜️ In Progress | Centralized control panel to manage or monitoring Rasbora cluster.  | Enterprise |
| SafeGuard | ⬜️ In Progress | Offers protection on all Rasbora systems, with alerts for abuse, authentication, permissions, and authorization. | Enterprise |
| CrashReporter | ⬜️ In Progress| Reporting crash details and performance alerts in the event of system crashes.| Enterprise |

## Transcoding Strategies

Currently, we are focused on developing support for various encode/decode strategies tailored to different workload needs:

| Name             | Status          | Strategy                                                                   | Edition   |
|------------------|-----------------|----------------------------------------------------------------------------------|-----------|
| A-Z/V         | ✅ Ready        | Processes a single video file, decoding and encoding the entire file to produce multiple quality versions.                         | Community |
| S-C/V         | ⬜️ In Progress  | Divides video files into small chunks; multiple instances of VideoTranscoding work on different chunks, later reassembling them into one video file with multiple qualities. | Enterprise |
| F-F/V         | ⬜️ In Progress  | Analyzes video files and generates frames; each VideoTranscoding instance works on frames, later reassembling them into a video file with multiple qualities. | Enterprise |

## Prerequisites

Here is what you need to be able to run Rasbora.

    Git (Version: >=2.x.x)
    Docker Compose (Version: >=2.24.x)
    Docker (Version: >=25.x.x)

## Deployment (Installation)

| Method      | Type |Status  | Docs                                                           |
|-------------|-------|---------|----------------------------------------------------------------|
| Standalone   |Self-Hosted| ✅ Ready | [Deploy Rasbora as Standalone](https://github.com/openseawave/rasbora/blob/main/docs/standalone.md) |
| Distributed Cluster Bare Metal |Self-Hosted| ⬜️ In Progress | |
| [<img width="44" height="44" src="https://github.com/openseawave/rasbora/blob/main/docs/kubernetes.png?raw=true">](https://github.com/openseawave/rasbora) Kubernetes Cluster | Self-Hosted | ⬜️ In Progress  | |
| [<img width="44" height="44" src="https://github.com/openseawave/rasbora/blob/main/docs/aws.png?raw=true">](https://github.com/openseawave/rasbora) AWS |EC2| ⬜️ In Progress  | |
| [<img width="44" height="44" src="https://github.com/openseawave/rasbora/blob/main/docs/hetnzer.png?raw=true">](https://github.com/openseawave/rasbora) Hetnzer | Bare Metal Server| ⬜️ In Progress  | |
| [<img width="44" height="44" src="https://github.com/openseawave/rasbora/blob/main/docs/ovhcloud.jpg?raw=true">](https://github.com/openseawave/rasbora) OVHCloud | Bare Metal Server| ⬜️ In Progress | |
| [<img width="44" height="44" src="https://github.com/openseawave/rasbora/blob/main/docs/azure.png?raw=true">](https://github.com/openseawave/rasbora) Azure |VM |⬜️ In Progress  ||
| [<img width="44" height="44" src="https://github.com/openseawave/rasbora/blob/main/docs/digitalocean.png?raw=true">](https://github.com/openseawave/rasbora) Digitalocean |Droplets| ⬜️ In Progress  | |

## Supported Storage System

The supported storage systems and their current status:

| Method     | Type | Supported |Status|
|--------------|-----------|-------|----|
| SSD/HHD     | LocalStorage |✅  Ready       |Done  |
| S3/Ceph     | ObjectStorage |✅  Ready        |Done  |
| Gluster    | Network | ⬜️ In Progress | In Progress |
| FreeNAS | Network| ⬜️ In Progress | In Progress |

## Supported Transcoder Engines

The video transcoder engines that are currently supported:

| Engine     | Supported |Status|
|--------------|-----------|-------|
| [<img width="44" height="44" src="https://github.com/openseawave/rasbora/blob/main/docs/ffmpeg.png?raw=true">](https://ffmpeg.org/ffmpeg.html) ffmpeg     | ✅  Ready        |Done  |
| [<img width="44" height="44" src="https://github.com/openseawave/rasbora/blob/main/docs/gstreamer.png?raw=true">](https://gstreamer.freedesktop.org/documentation/tutorials/index.html?gi-language=c) gstreamer | ⬜️In Progress | In Progress |

## Supported Hardware Acceleration

Rasbora's will support in future many hardware acceleration options:

| Hardware     |Engine |Supported| Image | Handler | Status |
|--------------|-------|---------|-------|---------|--------|
| CPU/x86-64| ffmpeg | ✅  Ready|jrottenberg/ffmpeg:4.4-alpine | src/videotranscoder/handlers/default.handler | Done |
| CPU/ARM64 | ffmpeg | ⬜️In Progress | In Progress| In Progress | In Progress |
| GPU/Apple Silicon| ffmpeg | ⬜️In Progress | In Progress | In Progress | In Progress |
| GPU/Nvidia| ffmpeg | ⬜️In Progress | In Progress | In Progress | In Progress |
| GPU/AMD| ffmpeg| ⬜️In Progress | In Progress |  In Progress | In Progress|

## Supported API Communications

Communicate seamlessly with Rasbora through RESTful APIs and explore upcoming gRPC and Websocket integrations:

| Method     | Type | Supported |Status| Docs |
|--------------|-----------|-------|----|-----|
| Restful     | ✅  Ready| application/json        |Done  | [<img width="44" height="44" src="https://github.com/openseawave/rasbora/blob/main/docs/swagger.png?raw=true">](http://localhost:3701/swagger/index.html) [<img width="44" height="44" src="https://github.com/openseawave/rasbora/blob/main/docs/postman.png?raw=true">](https://github.com/openseawave/rasbora/blob/main/postman.json)|
| gRPC         | ⬜️ In Progress       | application/protobuf |In Progress  |In Progress |
| Websocket    | ⬜️ In Progress       | application/json |In Progress  |In Progress |

Note: To access the Swagger documentation, ensure that Rasbora is running, if you change the port of Task Manager component adjust the documentation URL accordingly.

## Supported Callback Methods

Current supported callback methods and their current status:

| Protocol     | Supported | Type |Status|
|--------------|-----------|------|-------|
| HTTP/1.1     | ✅  Ready      | application/json |Done  |
| gRPC         | ⬜️  In Progress      | application/protobuf |In Progress |
| Websocket    | ⬜️  In Progress      | application/json |In Progress  |

## Supported Queue/Database Systems

The supported queue/database systems and their current status:

| Engine     | Supported |Status|
|--------------|-----------|-------|
| [<img width="44" height="44" src="https://github.com/openseawave/rasbora/blob/main/docs/redis.png?raw=true">](https://redis.io/docs/data-types/sorted-sets/) Redis     | ✅  Ready       |Done  |
| [<img width="44" height="44" src="https://github.com/openseawave/rasbora/blob/main/docs/rabbitmq.png?raw=true">](https://www.rabbitmq.com/docs/queues) RabbitMQ  | ⬜️In Progress        |In Progress|
| [<img width="44" height="44" src="https://github.com/openseawave/rasbora/blob/main/docs/activemq.png?raw=true">](https://activemq.apache.org/components/classic/documentation/) ActiveMQ  | ⬜️In Progress |In Progress|
| [<img width="44" height="44" src="https://github.com/openseawave/rasbora/blob/main/docs/awssqs.png?raw=true">](https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/welcome.html) Amazon SQS| ⬜️In Progress |In Progress|

## Supported Log Collectors

Manage and analyze your logs with Rasbora, we will supporting various log collectors such as Logstash, Grafana Loki, and more:

| Engine     | Supported |Status|
|--------------|-----------|-------|
| STDOUT     | ✅   Ready      |Done  |
| Log Files  | ✅  Ready     |Done|
| Database | ✅   Ready    |Done|
| Logstash | ⬜️In Progress |In Progress|
| Grafana Loki| ⬜️ In Progress|In Progress|
| Logwatch | ⬜️ In Progress|In Progress|

## Supported Crash Reporting Methods

The current supported crash methods and their status:

| Method    | Supported |Status|
|--------------|-----------|-------|
| [<img width="44" height="44" src="https://github.com/openseawave/rasbora/blob/main/docs/sentry.png?raw=true">](https://github.com/openseawave/Rasbora) Sentry| ⬜️In Progress |In Progress|
| [<img width="44" height="44" src="https://github.com/openseawave/rasbora/blob/main/docs/datadog.jpg?raw=true">](https://github.com/openseawave/Rasbora) Datadog| ⬜️In Progress |In Progress|

## Supported Monitoring Methods

The supported monitoring methods and their current status:

| Method    | Supported |Status|
|--------------|-----------|-------|
| [<img width="44" height="44" src="https://github.com/openseawave/rasbora/blob/main/docs/redis.png?raw=true">](https://redis.io/docs/data-types/streams/) Redis/Streams|✅   Ready      |Done  |
| [<img width="44" height="44" src="https://github.com/openseawave/rasbora/blob/main/docs/grafana.png?raw=true">](https://grafana.com/docs/grafana/latest/) Grafana   | ⬜️In Progress |In Progress|
| [<img width="44" height="44" src="https://github.com/openseawave/rasbora/blob/main/docs/prometheus.png?raw=true">](https://prometheus.io/docs/introduction/overview/) Prometheus| ⬜️In Progress |In Progress|

## Support

We offer different types of support depending on the project size. You can choose the level of support that suits your expertise and requirements:

| Type                  | Limitations             | Methods                                       | Status       |
|-----------------------|-------------------------|-----------------------------------------------|--------------|
| Community Edition     | No limits on CE-components| Github Issues, Community Forum, Reddit, Discord| ✅ Ready      |
| Enterprise Edition    | CE-components & EE-components | Calls, Email, Chat, Panic Button, Direct 1-to-1 Engineer| ⬜️ In Progress|

To request support for the Enterprise Edition, please contact us at rasbora.enterprise@openseawave.com

## TODO

⬜️ Full Testing Coverage

⬜️ Start Community Forum

⬜️ Q/A Documentation Site

## License

Copyright (c) 2022-2023 <https://rasbora.openseawave.com>
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
