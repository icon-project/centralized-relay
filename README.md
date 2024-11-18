# Centralized Relay

Centralized Relay serves as a trusted medium for message transfer between different chains.
The following chains are supported:

- ICON
- EVMs
- COSMOS
- STELLER
- SOLANA

## How to use ?

Refer to the [WIKI](https://github.com/icon-project/centralized-relay/wiki).

## Bitcoin Relay

How to run Slave server

### Prerequisites

Go 1.x installed
Set up your environment variables as required

### Start Slave 1

```bash
GO_ENV=slave PORT=8081 API_KEY=your_api_key go run main.go bitcoin
```

### Start Slave 2

```bash
GO_ENV=slave PORT=8082 API_KEY=your_api_key go run main.go bitcoin
```

### Start Master

```bash
 GO_ENV=master PORT=8080 SLAVE_SERVER_1=http://localhost:8081 SLAVE_SERVER_2=http://localhost:8082 API_KEY=your_api_key go run main.go bitcoin IS_PROCESS=1
```

This env to trigger call the slaves

- IS_PROCESS=1
