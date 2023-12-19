# Centralized Relay
A trusted centralized relay for message transfer between ICON and the various other chains.
Currently, communication with any of the EVM chains is supported.

## Relay installation and configuration
1. **Clone, checkout and install the latest release ([releases page](https://github.com/icon-project/centralised-relay/releases)).**

   *[Go](https://go.dev/doc/install) needs to be installed and a proper Go environment needs to be configured*

    ```shell
    $ git clone git@github.com:icon-project/centralised-relay.git
    $ cd centralised-relay
    $ make install
    ```
   Verify the installation and explore all the commands, sub commands and their usages.
    ```shell
   $ centralized-rly
    ```

2. **Create and store wallets as a json keystore file in proper location**
    Centralized relay requires wallets of each chain to communicate with the bridge
    contracts of each chain. So, create wallets for each chain and store in proper
    location.
    *for example:*
      * for evm: $HOME/wallets/evm/keystore.json
      * for icon: $HOME/wallets/icon/keystore.json
3. **Acquire the bridge contract addresses for each chain**
4. **Configure the wallets that you created earlier in step 2 as an admin in particular
   bridge addresses of the chains.**

5. **Configure the chains you want to relay messages between.**
   Centralized relay by default points to the path: ```$HOME/.centralized-relay/config.yaml``` 
   as a config file for storing the configuration for each chain that you want to relay messages in between.

   To add the chain config files run the following command:
   ```shell
   $ centralized-rly chains add --config-path <your-config-path> --file <chain-config-file-path>
   ```
   If you don't provide the ```--config-path```, the chain config will be added to default config path.
   ```chain-config-file-path``` should be the path of json file where you have the metadata of the 
   chain that you want to add to config. The content of chain config file should be as that of [here](/example/configs). 
   You can list the added chains using the following the command:
   ```shell
   $ centralized-rly chains list
   ```

6. **Run(start) the centralized relay**

   Centralized relay can be started using the command:
   ```shell
   $ centralized-rly start --config-path <your-config-path> --db-path <your-db-path>
   ```
   ```your-db-path``` refers to the path where the centralized-rly stores persistent data(messages)
   that is relayed or to be relayed between chains. The default path is ```$HOME/.centralized-relay/data```
   Once the relay is up and running, you can test for relaying messages from one chain to another chain. For
   testing and demonstration, please refer [here]()
  

  