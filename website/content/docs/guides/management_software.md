---
title: "Managament Software"
date: 2023-12-18T10:00:00+02:00
lastmod: 2023-12-18T10:00:00+02:00
draft: false
weight: 30
toc: true
---
## Bandcommand

### Configuration

The MQTT broker can be configured using the constants defined in `cmd/main.go`.

### Building
To build the management software, you need to [install Go](https://go.dev/doc/install).

Inside the `bandcommand` directory, use
```shell
go build -o bandcommand ./cmd
```
which produces a single executable called `bandcommand`.

You can also use
```shell
go run ./cmd
```
to start bandcommand directly.

### Tests

To generate the mocks used in some of the tests, [mockgen](https://github.com/uber-go/mock) is required.

```shell
go generate _mocks/gen.go
go test ./... -v
```

## MQTT Broker

We use MQTT to send messages between the Barmbands and Bandcommand, which requires an MQTT Broker.
If you want to run your own Broker, you can use the Docker Compose setup in the `broker` directory.

```shell
docker compose up
```


