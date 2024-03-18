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

**Additinal Context:**

- All the keystore relayer files are located in the `keystore` directory.
  `ls $HOME/.centralized-relay/keystore`

- The version `1.0.0` keystore files for chain are located in the inside the its `nid` directory in a following format:
  `keystore/<nid>/<wallet-address>.json`