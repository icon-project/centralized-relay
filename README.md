# Centralized Relay

Centralized Relay serves as a trusted medium for message transfer between different chains.
The following chains are supported:

- ICON
- AVALANCHE

## Prerequisites

- **Go Language**: Version 1.21.0 or higher. Make sure to have Go properly installed and configured on your system.
  For detailed installation instructions, please visit the [Go Installation Guide](https://go.dev/doc/install).

- **GNU Make**: This is an essential build automation tool used to compile and build the project from source files.
  GNU Make streamlines the process of generating executables and other non-source components of the program.
  To install GNU Make on your system, refer to the [GNU Make Official Website](https://www.gnu.org/software/make/).

## Contract Addresses

| Chain     | xCall Address                               | Connection Address                         | Networks |
|-----------|---------------------------------------------|--------------------------------------------|----------|
| ICON      | cx15a339fa60bd86225050b22ea8cd4a9d7cd8bb83  | cxb2b31a5252bfcc9be29441c626b8b918d578a58b | lisbon   |
| AVALANCHE | 0xAF180CDFE82578dB128088C4D8aa60a38E5CF505  | 0x2500986cCD5e804B206925780e66628e88fE49f3 | fuji     |

## Relay installation and configuration

1. **Clone, checkout and install the latest release ([releases page](https://github.com/icon-project/centralised-relay/releases)).**

    ```shell
    git clone git@github.com:icon-project/centralised-relay.git
    cd centralised-relay
    make install
    ```

   Verify the installation and explore all the commands, sub commands and their usages.

    ```shell
   crly
    ```

2. **Wallet Creation**

   Centralized Relay necessitates the use of specific wallets in a JSON keystore file format for each chain to facilitate communication with their respective bridge(connection) contracts.
   Use the following commands to create wallet for each chain:

    - For evm:

      ```shell
      geth account new --keystore <directory-for-keystore> --password <keystore-password>
      ```

    - For icon:

      ```shell
      goloop ks gen --out <directory-for-keystore> --password <keystore-password>
      ```

   Also make sure to load sufficient fake test balance in these wallets to perform operations like contract deployment and sending transactions.
   Balance can be loaded using faucets for each testnet. Look for the following faucets for specific testnet.
    - [Mumbai testnet faucet](https://mumbaifaucet.com/).
    - [Icon testnet faucet](https://faucet.iconosphere.io/).

   Next store these keystore files in a proper location. The wallet keystore file needs to be referenced in the chain configuration file. Follow next steps for chain configuration.

4. **Configure the chains you want to relay messages between.**

   Centralized Relay uses a configuration file to manage the settings for each chain. By default, it looks for the configuration at `$HOME/.crly/config.yaml`.
   Run the following command to create a configuration file at your preferred location specified by flag `--home`. If the path is not provided
   config file will be created in the default config path.

   ```shell
   centralized-rly config init --home $HOME/.crly
   ```

   **Adding Chain Configurations**:

    - To include a new chain, execute the command:

      ```shell
      centralized-rly chains add --config-path <your-config-path> --file <chain-config-file-path>
      ```

        - If `config-path` is not specified, the configuration will default to the standard path.
        - The `<chain-config-file-path>` should be the path to a JSON file containing the chain's metadata. The structure of this file can be modeled similar to the example found [here](/example/configs).

   **Listing Configured Chains**:
    - To view the list of configured chains, use the command:

      ```shell
      centralized-rly chains list
      ```

    - To view the config file, use the command:

       ```shell
      centralized-rly config show
      ```

   Ensure that each chain's configuration is accurately set up to facilitate smooth message relay between the specified networks.

5. **Run(start) the centralized relay**

   To start the Centralized Relay, use the following command:

   ```shell
   centralized-rly start --config-path <your-config-path> --db-path <your-db-path>
   ```

    - `db-path` refers to the path where the centralized-rly stores data(messages)
      that is relayed or to be relayed between chains. The default path is ```$HOME/.crly/data```

6. **Testing and Demo**

   Once the relay is up and running, you can test for relaying messages from one chain to another chain. For
   testing and demonstration, please refer [here]()

## Messages and Blocks Query

### 1. Message Queries

These commands are used for managing and querying messages in the `centralized-rly` database.

#### i. List Messages

- **Command:** `$ centralized-rly database messages list --page 1 --limit 2 --chain 0x2.icon`
- **Description and Parameters:**
  - List messages from the database.
  - `--page 1`: Specifies the page number for pagination.
  - `--limit 2`: Limits the number of messages displayed to 2.
  - `--chain 0x2.icon`: Filters messages belonging to the blockchain or network identified by `0x2.icon`.

#### ii. Remove Message

- **Command:** `$ centralized-rly database messages rm --sn 1 --chain 0x2.icon`
- **Description and Parameters:**
  - Removes a specific message from the database.
  - `--sn 1`: Identifies the serial number of the message to be removed.
  - `--chain 0x2.icon`: Specifies the blockchain or network of the message.

#### iii. Relay Message

- **Command:** `$ centralized-rly database messages rly --height 32359902 --sn 2 --chain 0x2.icon`
- **Description and Parameters:**
  - Relays or processes a specific message.
  - `--height 32359902`: Specifies the blockchain height for the message relay.
  - `--sn 2`: Identifies the sequence number of the message to be relayed.
  - `--chain 0x2.icon`: Indicates the blockchain or network for the message.

### 2. Block Queries

This command is used for viewing information about blocks in the database.

#### i. View Blocks

- **Command:** `$ centralized-rly database blocks view --chain 0x2.icon`
- **Description and Parameter:**
  - Displays information about blocks stored in the database.
  - `--chain 0x2.icon`: Filters blocks to those belonging to the specified blockchain or network.

### 3. Prune Database

- **Command:** `$ centralized-rly database prune`
- **Description:**
  - Deletes every blocks and messages from database.
