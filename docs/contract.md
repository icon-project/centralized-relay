# Contract

## 1. Introduction

The `contract` command is used to deploy and manage the fee on the chain.

## 2. Usage

```bash
centralized-relay contract [command] [flags]

Flags:
  -h, --help           help for contract

```

## 3. Commands

### Fee

1. Get fee for the chain

```bash
fee get [flags]

Flags:
    -c, --chain string      Chain ID
    -n, --network string    Network ID
    -r, --res-fee bool      Include response fee
  ```

2. Get fee for the chain

```bash
fee set [flags]

Flags:
    -c, --chain string       Chain ID
        --network string     Network ID
        --msg-fee string     Message fee
        --res-fee string     Response fee
```

3. Claim fee for the chain

```bash
fee claim [flags]

Flags:
    -c, --chain string   Chain ID
```
