# Config

## Command Line Options

The `config` command is used to manage the configuration file for the centralized relay.

## Usage

```bash
centralized-relay config [command] [flags]
```

## Flags

- `--config`: The config file to use.

## Commands

### Initialize the config file

```bash
init
```

This command will create the config file if it does not exist.

### Show the config file

```bash
show
```

This command will show the config file contents.

## Introduction

The config file should be either a JSON or YAML file.

The example below is a YAML file.

```yaml
global:
  timeout: 10s
  kms-key-id: f5c550ca-a6f2-4597-895c-4846ab8e4ad2
chains:
  avalanche:
    type: evm
    value:
      rpc-url: https://nd-507-756-317.p2pify.com/b01d81c64a3b7454f7a40a20ef44a1ce/ext/bc/C/rpc
      verifier-rpc-url: ""
      start-height: 0
      address: 0xB89596d95b2183722F16d4C30B347dadbf8C941a
      gas-min: 0
      gas-limit: 100056000
      contracts:
        xcall: 0x3f6391be658E9e163DA476b6ed1F6135cc29a376
        connection: 0x475d58a524ABDCe114847AD11F6172B9558b0af2
      concurrency: 0
      block-interval: 2s
      finality-block: 10
      nid: 0xa869.fuji
  icon:
    type: icon
    value:
      rpc-url: https://lisbon.net.solidwallet.io/api/v3/
      address: hxb6b5791be0b5ef67063b3c10b840fb81514db2fd
      start-height: 0
      step-min: 1
      step-limit: 2000000000000000000
      contracts:
        xcall: cx15a339fa60bd86225050b22ea8cd4a9d7cd8bb83
        connection: cx8d02efb10359105f7e033149556eaea531a3740e
      network-id: 2
      finality-block: 10
      nid: 0x2.icon

  archway:
    type: cosmos
    value:
      chain-id: constantine-3
      nid: archway
      rpc-url: https://rpc.constantine.archway.io:443
      grpc-url: grpc.constantine.archway.io:443
      keyring-backend: memory
      address: archway185jdmecuvmep8puqp0gpjszmy2w8ykes6ecxk8
      account-prefix: archway
      start-height: 0
      contracts:
        xcall: archway1h04c8eqr99dnsw6wqx80juj2vtuxth70eh65cf6pnj4zan6ms4jqshc5wk
        connection: archway1s0lw2w40g76ssvjd7at2en35ed9xhfskzgpsuemasn27tskey7wqyxfslm
      denomination: aconst
      gas-prices: 900000000000aconst
      gas-adjustment: 1.5
      max-gas-amount: 4000000
      min-gas-amount: 20000
      block-interval: 6s
      tx-confirmation-interval: 6s
      broadcast-mode: sync
      sign-mode: SIGN_MODE_DIRECT
      simulate: true
      finality-block: 10
```

## Explantion

Common configuration.

| Field  | Description | Allowed Values | Example | Type |
| -----  | ----------- | -------------- | ------- | ---- |
| timeout | The timeout for the chains. | 10s | 10s | duration |
| kms-key-id | The KMS key ID used for keystore encryption. | --- | --- | uuid |
| chains | The chains that will be used. | --- | cosmos, evm, icon | map |
| type | The type of the chain. | evm | evm | string |
| rpc-url | The RPC URL for the chain. | --- | --- | url |
| verifier-rpc-url | The verifier RPC URL for the chain. Used for the chains that have a verifier RPC URL. | --- | --- | url |
| start-height | This is the past chain height for the chain when starting the relayer. If the start height is set to 0, then the relayer will start from the latest block height. If the start height is set to a specific block height, then the relayer will start from that block height. If the future block height set, then the relayer will refuse to start. | 0  | 123 | int |
| address | The keystore/wallet for the chain currently being used. | --- | --- | string |
| contracts | The contracts for the chain. | --- | --- | map |

The configuration file is divided into two sections: global and chains.

Chain specific configurations.

**EVM**

| Field  | Description | Allowed Values | Example | Type |
| -----  | ----------- | -------------- | ------- | ---- |
| gas-min | The minimum gas price allowed for the transcation. | 0 | 0 | int |
| gas-limit | The maximum allowed gas limit for the transcation. | 100056000 | 100056000 | int |

**ICON**

| Field  | Description | Allowed Values | Example | Type |
| -----  | ----------- | -------------- | ------- | ---- |
| network-id | The network ID for the chain. | 1 | 1 | int |
| nid | The NID for the chain, derived from connection contract when deploying. | 0x2.icon | 0x2.icon | string |
| step-min | The minimum step price for the chain. | 1 | 1 | int |
| step-limit | The maximum step limit for the chain. | 2000000000000000000 | 2000000000000000000 | int |
| finality-block | The finality block for the chain. | 10 | 10 | int |
| rpc-url | The RPC URL for the chain. | any valid rpc url specific to the chain | <https://lisbon.net.solidwallet.io/api/v3/> | url |

**COSMOS**

| Field  | Description | Allowed Values | Example | Type |
| -----  | ----------- | -------------- | ------- | ---- |
| chain-id | The chain ID for the chain. | constantine-3 | constantine-3 | string |
| rpc-url | The RPC URL for the chain. | any valid rpc url | <https://rpc.constantine.archway.io:443> | url |
| grpc-url | The gRPC URL for the chain. | any valid grpc url | grpc.constantine.archway.io:443 | url |
| keyring-backend | The keyring backend for the chain. | memory, test, file | memory | string |
| account-prefix | The account prefix for the chain. | archway | archway | string |
