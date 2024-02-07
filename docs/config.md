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

The example below is a JSON file.

```yaml
global:
    timeout: 10s
    kms-key-id: aa3b3e3e-3e3e-3e3e-3e3e-3e3e3e3e3e3e
chains:
    avalanche:
        type: evm
        value:
            rpc-url: https://api.avax.network/ext/bc/C/rpc
            verifier-rpc-url: ""
            start-height: 29743735
            keystore: 0xB89596d95b2183722F16d4C30B347dadbf8C941a
            gas-price: 100056000
            gas-limit: 200000
            contracts:
              xcall: 0x3f6391be658E9e163DA476b6ed1F6135cc29a376
              connection: 0x2500986cCD5e804B206925780e66628e88fE49f3
            concurrency: 0
            finality-block: 10
            nid: 0xa869.fuji
    icon:
        type: icon
        value:
            rpc-url: https://lisbon.net.solidwallet.io/api/v3/
            keystore: hxb6b5791be0b5ef67063b3c10b840fb81514db2fd
            start-height: 34035738
            contracts:
              xcall: cx15a339fa60bd86225050b22ea8cd4a9d7cd8bb83
              connection: cxb2b31a5252bfcc9be29441c626b8b918d578a58b
            network-id: 2
            nid: 0x2.icon
```

## Explantion

- `global`: The global configuration for the chains.
  - `timeout`: The timeout for the chains.
  - `kms-key-id`: The KMS key ID used for keystore encryption.

- `chains`: The chains that will be used.
  - `evm`: The configuration for the evm compataible chains.
    - `type`: The type of the chain.
    - `value`: The configuration for the chain.
      - `rpc-url`: The RPC URL for the chain.
      - `verifier-rpc-url`: The verifier RPC URL for the chain.
        Used for the chains that have a verifier RPC URL.

      - `start-height`: The past start height for the chain when starting the relayer.

      This is the past chain height for the chain when starting the relayer.

      If the start height is set to 0, then the relayer will start from the latest block height.

      If the start height is set to a specific block height, then the relayer will start from that block height.

      If the future block height is set, then the relayer will refuse to start.

      - `keystore`: The keystore for the chain currently being used.
      - `gas-price`: The gas price set for the transcation.
      - `gas-limit`: The gas limit for the transcation.
      - `contracts`: The contracts for the chain.
        - `xcall`: The xcall contract for the chain.
        - `connection`: The connection contract for the chain.
      - `concurrency`: The concurrency for the chain.

      This is the number of concurrent transactions that can be sent when it is started.
      For example if the remaining chain height to sync is 1000 and the concurrency is 10, then 10 transactions will be sent concurrently.

      The next 10 transactions will be sent when the previous 10 transactions are confirmed.

      If the concurrency is set more than the remaining chain height to sync, then the remaining chain height to sync will be used as the concurrency.

      - `finality-block`: The finality block for the chain.
      - `nid`: The NID for the chain, derived from connection contract when deploying.

  - `icon`: The configuration for the ICON chain.
    - `type`: The type of the chain.
    - `value`: The configuration for the chain.
      - `rpc-url`: The RPC URL for the chain.
      - `keystore`: The keystore for the chain in use.
      - `start-height`: The past start height for the chain when starting the relayer.

      This is the past chain height for the chain when starting the relayer.

      If the start height is set to 0, then the relayer will start from the latest block height.

      If the start height is set to a specific block height, then the relayer will start from that block height.

      If the future block height is set, then the relayer will refuse to start.
      - `contracts`: The contracts for the chain.
        - `xcall`: The xcall contract for the chain.
        - `connection`: The connection contract for the chain.
      - `network-id`: The network ID for the chain.

      This is the network ID for the chain when deploying the connection contract.

      Fo example, the network ID for the ICON mainnet is 1 and the network ID for the ICON lisbon testnet is 2.

      - `nid`: The NID for the chain, derived from connection contract when deploying.

```
