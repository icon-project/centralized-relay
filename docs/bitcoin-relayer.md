# Bitcoin Relayer

This is a relayer for Bitcoin that allows you to send or receive BTC/Rune from your Bitcoin taproot wallet to your wallet on other chains.

## Prerequisites

- Understand the basics of Bitcoin transactions and Multisig taproot wallets
- Understand about Bitcoin RPC and the related APIs (Unisat, Quicknode, etc.)

## Configuration

config.yaml is the main configuration file for the Bitcoin relayer. It includes the following sections:

```yaml
bitcoin:
  type: bitcoin
  value:
    rpc-url: # Bitcoin RPC URL
    rpc-user: # Bitcoin RPC User
    rpc-password: # Bitcoin RPC Password
    address: # Bitcoin Taproot multisig address
    unisat-url: # Unisat API URL
    unisat-key: # Unisat API Key
    network-id: # Bitcoin Network ID (1: Bitcoin Mainnet, 2: Bitcoin Testnet)
    op-code: # Bitcoin OP Code (0x5e: default)
    finality-block: # Bitcoin Finality Block (10: default)
    nid: 0x2.btc # Bitcoin NID (0x1.btc: Bitcoin Mainnet, 0x2.btc: Bitcoin Testnet)
    chain-name: bitcoin # Bitcoin Chain Name
    mode: master # master or slave
    slave-server-1: # Slave Server 1 URL (only used when mode is master)
    slave-server-2: # Slave Server 2 URL (only used when mode is master)
    port: 8082 # Server Port
    api-key: key # Slave Server API Key (Using to authenticate between the master and slave servers)
    masterPubKey: # Master Public Key (public key of the master wallet)
    slave1PubKey: # Slave 1 Public Key (public key of the slave wallet 1)
    slave2PubKey: # Slave 2 Public Key (public key of the slave wallet 2)
    relayerPrivKey: # Relayer Private Key (private key of the relayer it depends the deployed server that which start for master/slave1/slave2 server)
    recoveryLockTime: # Recovery Lock Time (recovery lock time of the master wallet)
    start-height: # Start Height (start height)
```

# How it works

The mechanism of the Bitcoin relayer is based on the master-slave architecture.

## Relayer Multisig Taproot Wallet

- Because in Bitcoin network that not support the smart contract, so the Relayer Multisig Taproot Wallet is the wallet that will receive the BTC/Rune, keep them and the logic to handle send and receive BTC/Rune will process based on this wallet

- The private key of the Relayer Multisig Taproot Wallet was generated from the `relayerPrivKey` in the config.yaml file

- This Wallet Address was combined from 3 different public keys: the master public key, slave1 public key and slave2 public key

- To spend token from this wallet, it needs 3 signatures from the 3 public keys (the master public key and one of the slave public keys)

**_Note: The order of the public keys in the wallet address is important, it must be the same between the order when we generate the Relayer Multisig Taproot Wallet and sign the transaction_**

## Master Server

The master server will crawl Bitcoin transaction from the Bitcoin network and check if the transaction is a valid transaction with recipient is Relayer multigsig wallet address and the condition that contain value of OP_RETURN data is the same as the `op-code` (OP_RETURN OP_PUSHNUM_14) in the config.yaml file.
The master server is the main server that handles:

- Requesting the slave servers to sign the transactions
- Combining the signatures from the slave servers and itself then broadcasting the transaction to the Bitcoin network

## Slave Servers

It works as the same with the master server, but the slave servers will not broadcast transaction instead of they crawl transactions and cache them, and waiting for the master request to sign the transactions and send the signature back to the master server.

## Data Structure

Based on the XCall message structure, The Bitcoin Relayer was designed and implemented to parse message with structure `OP_RETURN OP_PUSHNUM_14 YOUR _PAYLOAD`.

Because OP_RETURN is only limited to 80 bytes by Bitcoin Core's default standardness rules, so the payload `(YOUR_PAYLOAD)` will be split into multiple parts with the maximum size of 76 bytes for each part and to grarentee .

The Bitcoin Relayer will decode the message from the Bitcoin transaction and parse the payload to `BridgeDecodedMsg` data structure.

```golang
type BridgeDecodedMsg struct {
	Message    *XCallMessage
	ChainId    uint8
	Receiver   string
	Connectors []string
}

type XCallMessage struct {
	MessageType  uint8
	Action       string
	TokenAddress string
	From         string
	To           string
	Amount       []byte
	Data         []byte
}
```

**Example:**

```golang
bridgeMsg := BridgeDecodedMsg{
		Message:  XCallMessage{
			MessageType:  1,
			Action:       "Deposit",
			TokenAddress: "0:1",
			To:           "0x2.icon/hx452e235f9f1fd1006b1941ed1ad19ef51d1192f6",
			From:         "tb1pgzx880yfr7q8dgz8dqhw50sncu4f4hmw5cn3800354tuzcy9jx5shvv7su",
			Amount:       new(big.Int).SetUint64(100000).Bytes(),
			Data:         []byte(""),
		},
		ChainId:  1,
		Receiver: "cxfc86ee7687e1bf681b5548b2667844485c0e7192",
		Connectors: []string{
			"cx577f5e756abd89cbcba38a58508b60a12754d2f5",
		},
	}
```

### Deploy the Relayer

Since the Bitcoin relayer works based on the master-slave architecture, so we need to deploy seperate the master server and at least 2 slave servers, in totally there are 3 servers need to be run at the same time.
