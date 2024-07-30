# Centralized Relay

Centralized Relay serves as a trusted medium for message transfer between different chains.
The following chains are supported:

- ICON
- AVALANCHE
- COSMOS

## Contract Addresses

| Chain     | xCall Address                               | Connection Address                         | Networks | Wallets  |
|-----------|---------------------------------------------|--------------------------------------------|----------| -------- |
| ICON      | cx15a339fa60bd86225050b22ea8cd4a9d7cd8bb83  | cxb2b31a5252bfcc9be29441c626b8b918d578a58b | lisbon  | hxb6b5791be0b5ef67063b3c10b840fb81514db2fd |
| AVALANCHE | 0xAF180CDFE82578dB128088C4D8aa60a38E5CF505  | 0x2500986cCD5e804B206925780e66628e88fE49f3 | fuji    | 0xB89596d95b2183722F16d4C30B347dadbf8C941a |
| ICON      | cxa07f426062a1384bdd762afa6a87d123fbc81c75  | cxdada6921d08fbf37c6f228816852e58b219cc589 | mainnet | hxda27114a959a3351f3613b055ca96f8f8cb34cbe |
| AVALANCHE | 0xfc83a3f252090b26f92f91dfb9dc3eb710adaf1b  | 0xCC7936eA419516635fC6fEb8AD2d4341b5D0C2B3 | mainnet | 0xebA66Ad34CCEB70669eddbaA8c9Fb927d41fE2d7 |

## How to use ?

Refer to the [WIKI](<https://github.com/icon-project/centralized-relay/wiki>).

## Bitcoin Relay

How to run Slave server 

### Prerequisites

Go 1.x installed
Set up your environment variables as required

### Start Slave

```bash
GO_ENV=master MASTER_SERVER=http://localhost:8080 SLAVE_SERVER=http://localhost:8081 API_KEY=your_api_key go run main.go bitcoin
```

### Start Master

```bash
GO_ENV=slave MASTER_SERVER=http://localhost:8080 SLAVE_SERVER=http://localhost:8081 API_KEY=your_api_key go run main.go bitcoin
```
