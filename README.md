# Centralized Relay
Centralized Relay serves as a trusted medium for message transfer between different chains.
Currently, supports the following chains:
- ICON
- AVALANCHE
---
## Features
- **Quick integration with DAAP**
- **Prevention from abnormal transactions**

---
## Prerequisites

- **Go Language**: Version 1.20.0 or higher. Make sure to have Go properly installed and configured on your system. 
   For detailed installation instructions, visit the [Go Installation Guide](https://go.dev/doc/install).

- **GNU Make**: This is an essential build automation tool used to compile and build the project from source files. 
   GNU Make streamlines the process of generating executables and other non-source components of the program. 
   To install GNU Make on your system, refer to the [GNU Make Official Website](https://www.gnu.org/software/make/).

---

## Contract Addresses
| Chain     | xCall Address                               | Bridge Address                              | Networks | 
|-----------|---------------------------------------------|---------------------------------------------|----------|
| ICON      | cx15a339fa60bd86225050b22ea8cd4a9d7cd8bb83  | cxb2b31a5252bfcc9be29441c626b8b918d578a58b  | lisbon   | 
| AVALANCHE | 0xAF180CDFE82578dB128088C4D8aa60a38E5CF505  | 0x2500986cCD5e804B206925780e66628e88fE49f3  | fuji     |



---
## Relay installation and configuration
1. **Clone, checkout and install the latest release ([releases page](https://github.com/icon-project/centralised-relay/releases)).**

    ```shell
    $ git clone git@github.com:icon-project/centralised-relay.git
    $ cd centralised-relay
    $ make install
    ```
   Verify the installation and explore all the commands, sub commands and their usages.
    ```shell
   $ centralized-rly
    ```
2. **Wallet Creation and Storage**

   Centralized Relay necessitates the use of specific wallets for each chain to facilitate communication with their respective bridge contracts. Follow these steps to create and properly store the wallets:

   - **Create Wallets**: Generate a new wallet for each chain you intend to interact with.
   - **Store in JSON Keystore Format**: Save the wallet details in a JSON keystore file format.

   *Example Storage Locations:*
   - evm:  `$HOME/wallets/evm/keystore.json`
   - icon: `$HOME/wallets/icon/keystore.json`

3. **Configure the wallets that you created in step 2 as an admin in particular
   bridge contract addresses of the chains.**

4. **Configure the chains you want to relay messages between.**

   The Centralized Relay uses a configuration file to manage the settings for each chain. By default, it looks for the configuration at `$HOME/.centralized-relay/config.yaml`.

   **Adding Chain Configurations**:

   - To include a new chain, execute the command:
     ```shell
     centralized-rly chains add --config-path <your-config-path> --file <chain-config-file-path>
     ```
     - If `--config-path` is not specified, the configuration will default to the standard path.
     - The `<chain-config-file-path>` should be the path to a JSON file containing the chain's metadata. The structure of this file can be modeled after the examples found [here](/example/configs).

   **Listing Configured Chains**:
   - To view the list of configured chains, use the command:
     ```shell
     centralized-rly chains list
     ```

   Ensure that each chain's configuration is accurately set up to facilitate smooth message relay between the specified networks.

5. **Run(start) the centralized relay**

   To start the Centralized Relay, use the following command:
   ```shell
   $ centralized-rly start --config-path <your-config-path> --db-path <your-db-path>
   ```
   - `your-db-path` refers to the path where the centralized-rly stores data(messages)
   that is relayed or to be relayed between chains. The default path is ```$HOME/.centralized-relay/data```
   
6. **Testing and Demo**
   Once the relay is up and running, you can test for relaying messages from one chain to another chain. For
   testing and demonstration, please refer [here]()
  

  