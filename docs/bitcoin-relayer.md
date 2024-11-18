# Bitcoin Relayer

This is a relayer for Bitcoin that enables you to send or receive BTC or Runes from your Bitcoin Taproot wallet to your wallet on other blockchains.

## Prerequisites

- Understand the basics of Bitcoin transactions and multisig Taproot wallets.
- Gain knowledge about Bitcoin RPC and related APIs, such as Unisat, Quicknode, and others.

Note: In this implementation, we use third-party APIs, such as Quicknode, to crawl Bitcoin transactions and Unisat to retrieve BTC/RUNE UTXOs or fee rates. Therefore, it is necessary to prepare API keys for these services.

## Configuration

config.yaml is the main configuration file for the Bitcoin relayer. It includes the following configurations:

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
    unisat-wallet-url: # Unisat OPEN API URL (https://wallet-api.unisat.io for mainnet, https://wallet-api-testnet.unisat.io for testnet)
    request-timeout: 1000 # Request Timeout (ms)
    network-id: # Bitcoin Network ID (1: Bitcoin Mainnet, 2: Bitcoin Testnet)
    op-code: # Bitcoin OP Code (0x5e: default)
    finality-block: # Bitcoin Finality Block (10: default)
    nid: 0x1.btc # Bitcoin NID (0x1.btc: Bitcoin Mainnet, 0x2.btc: Bitcoin Testnet)
    chain-name: bitcoin # Bitcoin Chain Name
    recoveryLockTime: # Recovery Lock Time (recovery lock time of the master wallet)
    start-height: # Start Height (start height)
    mode: master # master or slave
    slave-server-1: # Slave Server 1 URL (only used when mode is master)
    slave-server-2: # Slave Server 2 URL (only used when mode is master)
    port: 8082 # Server Port
    api-key: key # Slave Server API Key (Using to authenticate between the master and slave servers)
    masterPubKey: # Master Public Key (public key of the master wallet)
    slave1PubKey: # Slave 1 Public Key (public key of the slave wallet 1)
    slave2PubKey: # Slave 2 Public Key (public key of the slave wallet 2)
    relayerPrivKey: # Relayer Private Key (private key of the relayer it depends the deployed server that which start for master/slave1/slave2 server)

```

# How it works

The mechanism of the Bitcoin relayer is based on the master-slave architecture.

## Relayer Multisig Taproot Wallet

- Since the Bitcoin network does not support smart contracts, the Relayer Multisig Taproot Wallet serves as the wallet to receive BTC/Rune, securely store them, and the logic for sending and receiving BTC/Rune will be based on this wallet.

- The private key of the Relayer Multisig Taproot Wallet is generated from the `relayerPrivKey` in the config.yaml file

- This wallet address is derived from three different public keys: the master public key, the slave1 public key, and the slave2 public key.

- To spend tokens from this wallet, it requires signatures from two out of the three public keys: the master public key and one of the slave public keys.

**_Note: The order of the public keys in the wallet address is crucial. It must remain consistent between the order used to generate the Relayer Multisig Taproot Wallet and the order used to sign the transaction_**

## Master Server

The master server will crawl Bitcoin transactions from the Bitcoin network and verify if the transaction is valid. A valid transaction must have the recipient as the Relayer Multisig Wallet address and must include an OP_14 data value matching the `op-code` specified in the `config.yaml` file.
The master server is the primary server responsible for:

- Requesting the slave servers to sign the transactions
- Combining the signatures from the slave servers and its own, then broadcasting the transaction to the Bitcoin network.

## Slave Servers

It functions similarly to the master server, but the slave servers do not broadcast transactions. Instead, they crawl transactions, cache them, and wait for requests from the master server to sign the transactions. Once signed, the slave servers send the signatures back to the master server.

## Data Structure

Based on the XCall message structure, the Bitcoin Relayer is designed and implemented to parse messages with the following structure: `OP_14 YOUR _PAYLOAD`.

Since we leverage op codes to send data, the limitation is 40 bytes per transaction due to Bitcoin Core's default standardness rules. Therefore, the payload `(YOUR_PAYLOAD)` will be split into multiple UTXOs (outputs), each with a maximum size of 40 bytes, including a dust amount of 547 sats for each part.

The Bitcoin Relayer will decode the message from the Bitcoin transaction and parse the payload into the `BridgeDecodedMsg` data structure.

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
			TokenAddress: "0:0",
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

Since the Bitcoin Relayer operates based on a master-slave architecture, the master server and at least two slave servers should be deployed separately. Ideally, a total of three servers should run simultaneously. However, for testing purposes, these servers can be deployed on the same machine.

#### Master Server Configuration

Here are some configuration differences between the master and slave servers:

```yaml
# ... config above
mode: master # master or slave
slave-server-1: # Slave Server 1 URL (only used when mode is master)
slave-server-2: # Slave Server 2 URL (only used when mode is master)
port: 8080 # Server Port (master server port)
api-key: key # Slave Server API Key (Using to authenticate between the master and slave servers)
masterPubKey: # Master Public Key (public key of the master wallet)
slave1PubKey: # Slave 1 Public Key (public key of the slave wallet 1)
slave2PubKey: # Slave 2 Public Key (public key of the slave wallet 2)
relayerPrivKey: # Relayer Private Key for master public key
```

#### Slave Server Configuration

For the slave servers, there is no need to configure `slave-server-1` and `slave-server-2`. The mode parameter should be set to `slave`.

```yaml
# ... config above
mode: slave
#slave-server-1
#slave-server-2
port: 8081 or 8082 # Slave Server Port, it depends the deployed server that which start for slave1/slave2 server
api-key: # Same with master server api-key (Using to authenticate between the master and slave servers)
relayerPrivKey: # Relayer Private Key for slave1 or slave2 public key, it depends the deployed server that which start for slave1/slave2 server
```

#### Start the Relayer

```bash
RELAY_HOME="YOUR_SOURCE_CODE_PATH" go run main.go start
```

### Implementation Details:

- Deposit BTC/Runes from Bitcoin to ICON
- Withdraw BTC/Runes from ICON to Bitcoin
- Rollback BTC/Runes in case the deposit fails
- Refund BTC if the bridge message amount does not match the output sent to the Relayer

#### Testing Results:

##### Deposit BTC Successfully

- Bitcoin tx: https://mempool.space/tx/9a9d955dff45c6cef6f4e41a12052dde21179069a2e17fe8f381f6c75e112b6a
- Connection tx:
  - https://tracker.icon.community/transaction/0x20f6364733a64882da22dcc06cc9086e0bbae6ec966197796aa53cdcfb419b26
- Xcall execute:
  - https://tracker.icon.community/transaction/0xc08aea4dde75f5624cf49dd00ba8ec8181f8d9c05db709a24d142e40f108c382

##### Deposit RUNE Successfully

- Bitcoin tx:
  - https://mempool.space/tx/924c7c6bd13f465b0b50cb8ad883544b22bfe54fae42e2ecfc9f9609a1b616f7
- Connection tx:
  - https://tracker.icon.community/transaction/0x72fba555202b3d0b45f70c6ba9ca9a00c162fd20d98e1c7fb93f6555ec7bd0ca
- Xcall execute:
  - https://tracker.icon.community/transaction/0xc240279f1fdd590c18546777692c8821a3fdff795bd514bdd29fe67c3d540e7f

##### Deposit BTC Failed

- Bitcoin tx:
  - https://mempool.space/tx/b84060ce292dd61f8490bae54f8354caa8642de730e5f409b72d67b05617dcb0
- Connection tx:
  - https://tracker.icon.community/transaction/0x09e3a5b9c8dcc3f3436eafcc6a01397e2353f81100aa79fdd4209deac5545b2b
- Xcall tx:
  - https://tracker.icon.community/transaction/0x3f0ed85491f053177c7fb9e308d287137aa8bb97a17c1b5fee5f61ccd6eeaf25
- Rollback tx:
  - https://mempool.space/tx/da35fb5971ee045c35139a8c2a0388ec3a79f22045322a9ac76204728b0bb486

##### Deposit RUNE Failed

- Bitcoin tx:
  - https://mempool.space/tx/89095d016b50644a328667cd5543b0f29c0f2a81242094ea7318bded49cf30a8
- Connection tx:
  - https://tracker.icon.community/transaction/0x27a9f171531d828b3e57f2232cb19be01bbd4ab835511531a2936ea3ed63d3f7
- Xcall tx:
  - https://tracker.icon.community/transaction/0x74a447a1c06f4ebdbfe8ddb222a8e89c0ffddd1c098d9b484e2c7b6b05ce7d97
- Rollback tx:
  - https://mempool.space/address/bc1p2sdwgq7j32j250w8h47fe9v3hyc8fl2rdftwhxp0r7ww89mcwrns5reskh

##### Deposit BTC with wrong amount, and got refund
*The refund amount does not include the transaction fee. If the transfer amount cannot cover the transaction fee, the refund transaction will be ignored.*

- Request tx:
  - https://mempool.space/tx/50aa0c67d8a533d3766bd2076a2bc57bb67de7d61e9f503db271e915f0f75bae
- Refund tx:
  - https://mempool.space/tx/6a976c2d6651020cce3c13f464128b19e7c318d825dda0d47a14025d94179c0a

##### Withdraw BTC Successfully

- ICON tx:
  - https://tracker.icon.community/transaction/0x3854443002829635830e679c83d41303ace0093a78320846aa6f543835ecf751
- Bitcoin tx:
  - https://mempool.space/tx/cf671e0ecc434e2cb06152bae30d35114d7fef8c1c3ec7ae60aea45691edf75b

##### Withdraw RUNE Successfully

- ICON tx:
  - https://tracker.icon.community/transaction/0x2fbb0aca1b99692b24baae68c2b451945db9eb829f09996eac01b5799bf35fc1
- Bitcoin tx:
  - https://mempool.space/tx/21f8ba718ba003e38ef291c1f8a6de7706fb49b2addfbb3eadf4bf1808d83a17

### Known Issues

- If a transaction sent to the Relayer Multisig Wallet contains a BTC amount that does not match the BTC amount defined in the XCall message, the relayer will refund the amount to the sender, minus the transaction fee.
- To stress test the system, you need to prepare a large number of BTC/RUNE UTXOs for the Relayer Multisig Wallet. This ensures the system has sufficient UTXOs to process transactions and avoids issues related to insufficient UTXOs.
- In the case of a rollback transaction, the relayer will refund the exact amount of BTC/RUNE that the user sent to the Relayer Multisig Wallet.

### How to build transaction

To build a transaction, please refer to this file relayer/chains/bitcoin/provider_mainnet_test.go for detailed instructions and guidelines.
