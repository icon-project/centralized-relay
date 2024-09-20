# Changelog

All notable changes to this project will be documented in this file.

## [1.7.0] - 2024-09-20

### Added
- STELLAR chain implementation.

### Changed
- Store only the latest Xcall-Pkg-ID in SUI config.
- Make transaction block query interval configurable from the config in SUI. Also set query page size to 15.

### Fixed
- Fixed issue with conn reset and default lasttry value  for route message for EVM.

## [1.6.0] - 2024-08-12

### Added
- SUI chain implementation.
- Add `LastProcessedTxStore` to store the chain specific information about the last processed transaction.

### Changed
- While switching in between multiple relayer wallet address, set admin before saving config.
- Add two new fields in the relayer's `Message` type: 
    - `DappModuleCapID` to track to which dapp module the message need to relayed. This is used for `ExecuteCall` and `ExecuteRollback` only.
    - `TxInfo` to store the information of transaction in which the message is contained in the source chain.
        This field is opaque to relayer's core and its type is known only in the chain's implementation.
- Changed the `Listener` method signature in the chain provider interface from `Listener(ctx context.Context, lastSavedHeight uint64, blockInfo chan *types.BlockInfo) error` to `Listener(ctx context.Context, lastProcessedTx types.LastProcessedTx, blockInfo chan *types.BlockInfo) error`.


## [1.5.1] - 2024-08-12

### Changed

- RPC retries and exponential backoffs.
- Websocket healthcheck timout from 1 min to 10 second.
- Initiate startup tasks early.
- EVM websocket healthcheck uses latest query, previosly genesis block query.

### Fixed

- Addressed wasm response format changes.
- WASM chain healthcheck frequency.
- WASM sdk context lockup issue.
- WASM duplicated block when batching requests.
- Recovery start from genesis block when databse is empty.
- Docker build issues.

## [1.5.0-rc1] - 2024-08-03

### Added

- WASM conditional batch polling.

### Changed

- RPC retries and exponential backoffs.
- Websocket healthcheck from 1 min to 10 second.
- Initiate startup tasks early.

### Fixed

- Addressed wasm response format changes.
- WASM chain healthcheck frquency.
- WASM sdk context lockup issue.
- WASM duplicated block when batching requests.
- Recovery start from genesis block when databse is empty.
- Docker build issues.

### Removed

- EVM and wasm unlimited concurrency.

## [1.4.1] - 2024-07-17

### Added

- Support for the rollbackMessage event
- Cosmwasmvm `v2.1.0` support

### Changed

- Use rpc instead of the websocket filterLogs `eth_getLogs` for the evm chain

## [1.3.4] - 2024-07-01

### Fixed

- Restore the keystore when it is not found

### Removed

- Redunant sequence increment for the cosmos chain

## [1.3.3] - 2024-07-01

### Fixed

- Avoid closing connection when errors are encountered
- Cover all tcp connection errors
- Other improvements and bug fixes

## [1.3.2] - 2024-07-01

### Fixed

- EVM RPC connection recovery using the errors message

## [1.3.1] - 2024-06-27

### Fixed

- WS connection recovery for the evm chain
- Poper error checking for the icon and cosmos chain

## [1.3.0] - 2024-06-23

### Fixed

- Polling fix for evm when ws errors

### Added

- Use gas price cap and gas tip for the evm chain

## [1.2.9] - 2024-06-19

### Added

- Use gas price cap and gas tip for the evm chain

### Fixed

- Use pending nonce instead of the latest nonce for the evm chain
- Other improvements and bug fixes
- CPU and memeory usage optimization, dropped by more than 100%
- Retry is more stable
- Exponential backoff for the retry count

### Changed

- mutext on router
- evm past polling for events is optimized, it can batch call now using config option
- cosmwasm block polling using batch size using config option

## [1.2.8] - 2024-06-06

### Changed

- Removed concurrency

## [1.2.7] - 2024-05-28

### Fixed

- Use on/off switch for the polling and subscriptions for recoveries
- Other improvements and bug fixes

### Changed

- Evm block mined is replaced by custom function

## [1.2.6] - 2024-05-26

### Fixed

- Cosmos contracts subscriptions respects the configured contracts
- RPC failures are are handled more elegently, switches to the polling and back to the subscriptions
- Address check validation for the manaul relay on the icon chain
- Other improvements and bug fixes

### Changed

- Icon `progressInterval` notification block is not incremented to handle rpc failures
- Default Block mined wait time is increased to 10 minutes
- Exponential backoff for the retry count

## [1.2.5] - 2024-05-18

### Fixed

- Wrong params sent when estimating gas for the evm chain `executeCall`

## [1.2.4] - 2024-05-17

### Added

- Support for the injective chain

### Fixed

- Gas Estimation for the evm chain
- Cosmos sdk global config bech32 prefixes
- Other improvements and bug fixes

## [1.2.3] - 2024-05-14

### Added

- Gas adjustment from config

## [1.2.2] - 2024-05-01

### Fixed

- Avoid nonce increment when fixing the nonce error while sending the transaction

## [1.2.1] - 2024-04-30

### Fixed

- Websocket connection disconnect issue with icon chain
- Use `eth_gasPrice` for the gas price calculation all the time
- Other improvements and bug fixes
- Use block mined timeout instead of polling when waiting for transcation

### Removed

- Icon redunant polling code

## [1.2.0] - 2024-04-09

### Added

- Full websocket listner support for all chains
- Auto clean expired messages to avoid disk space issues

### Changed

- CallMessage is only retried twice to avoid spamming retries
- Use only one websocket connection for maximum efficiency
- Use `/event` active listener instead of block search leading to significant performance gains

### Fixed

- Error handling for the websocket connection
- Start height for the icon chain
- Manual relay for icon chain using the height (on chain)
- Other improvements and bug fixes

### Removed

- Height sync is no longer necessary.

## [1.1.3] - 2024-03-27

### Added

- Route manually from height (on chain)

### Fixed

- Increase delivery failure by trying for per 15 seconds after initial failures.
- Panics when subscribing to the event result.
- AWS ec2 instance profile detection.
- Other improvements and bug fixes.

## [1.1.2] - 2024-03-22

### Fixed

- Region detection for AWS
- Priority 0 (high) for `start-height` evm
- Panic too many packets map access

## [1.1.1] - 2024-03-21

### Added

- Websocket support for evm chain

### Fixed

- AWS Region detection
- Static binary build

## [1.1.0] - 2024-03-18

### Added

- Support for cosmos chain
- Xcall contract support
- CallMessage event listener for all supported chains
- Fee related operations cmd. `getFee`, `setFee` and `claimFee`
- Structured events for easier event handling

### Changed

- Wallet encryption and decryption
- GO version to 1.22.1

### Fixed

- Incorrect icon chain balance calculation.
- SetAdmin also checks for the admin address to avoid setting the admin address to the same address.
- Retries is less error prone now, only retries after set interval (5s) when failed.
- Fixed the issue when initilizing config file.

### Migration from 1.0.0 to 1.1.0

We have added support for cosmos chains and xcall contract support. To migrate from 1.0.0 to 1.1.0, you need to update the configuration file and add the cosmos chain details. The configuration file is backward compatible, so you can add the cosmos chain details without affecting the existing configuration.

Major changes in this release is the wallet encryption and decryption. Previously we only encrypted the keystore password, now we encrypt the entire keystore file adding an extra layer of security. The relay will automatically decrypt the keystore file and use it to sign the messages.

We have also added the xcall execution contract support. The relay will now listen to the call message event and execute the xcall contract.

Exection will respect the fees set on configuration. The relay will now calculate the fees and execute the contract.

Migrate keystore files to the new format by running the following command:

**important**: Before running the command, make sure you have the AWS KMS key id. You can get the KMS key id by running the `crly config show` command.
