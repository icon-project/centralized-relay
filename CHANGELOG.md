# Changelog

All notable changes to this project will be documented in this file.

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

```shell
aws kms encrypt --key-id <kms-key-id> --plaintext fileb://path/to/keystore.json --output text --query CiphertextBlob | base64 -d > path/to/keystore/address
```

Example when migrating the icon chain keystore file where its nid is `0x2.icon` and the wallet address is `0x0B958dd815195F73d6B9B91bFDF1639457678FEb`:

verify keystore exists:

```shell
ls $HOME/.centralized-relay/keystore/0x2.icon/0x0B958dd815195F73d6B9B91bFDF1639457678FEb.json
```

Encrypt the keystore file:

```shell
aws kms encrypt --key-id <insert-key-id-here> --plaintext fileb://$HOME/.centralized-relay/keystore/0x2.icon/0x0B958dd815195F73d6B9B91bFDF1639457678FEb.json --output text --query CiphertextBlob | base64 -d > "$HOME/keystore/0x2.icon/0x0B958dd815195F73d6B9B91bFDF1639457678FEb"
```

Move the encrypted wallet passphrase to the new location:

  ```shell
  mv $HOME/keystore/0x2.icon/0x0B958dd815195F73d6B9B91bFDF1639457678FEb.password $HOME/.centralized relay/keystore/0x2.icon/0x0B958dd815195F73d6B9B91bFDF1639457678FEb.pass
  ```

### Additional Information

- All the keystore relayer files are located in the `keystore` directory.
  `ls $HOME/.centralized-relay/keystore`

- The version `1.0.0` keystore files for chain are located in the inside the its `nid` directory in a following format:
  `keystore/<nid>/<wallet-address>.json`

## [1.1.1] - 2024-03-21

### Added

- Websocket support for evm chain

### Fixed

- AWS Region detection
- Static binary build

## [1.1.2] - 2024-03-22

### Fixed

- Region detection for AWS
- Priority 0 (high) for `start-height` evm
- Panic too many packets map access

## [1.1.3] - 2024-03-27

### Added

- Route manually from height (on chain)

### Fixed

- Increase delivery failure by trying for per 15 seconds after initial failures.
- Panics when subscribing to the event result.
- AWS ec2 instance profile detection.
- Other improvements and bug fixes.

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

## [1.2.1] - 2024-04-30

### Fixed

- Websocket connection disconnect issue with icon chain
- Use `eth_gasPrice` for the gas price calculation all the time
- Other improvements and bug fixes
- Use block mined timeout instead of polling when waiting for transcation

### Removed

- Icon redunant polling code

## [1.2.2] - 2024-05-01

### Fixed

- Avoid nonce increment when fixing the nonce error while sending the transaction

# [1.2.3] - 2024-05-14

### Added

- Gas adjustment from config

## [1.2.4] - 2024-05-17

### Added

- Support for the injective chain

### Fixed

- Gas Estimation for the evm chain
- Cosmos sdk global config bech32 prefixes
- Other improvements and bug fixes
