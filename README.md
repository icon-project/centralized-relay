# Centralized Relay
A trusted centralized relay for message transfer between ICON and the other chains.
Currently, communication with any of the EVM chains is supported.

## Getting Started

### Prerequisites
The following tools and environment are required to set up and run the project. 

#### Go Programming Language.
Go is an open source programming language that makes it easy to build simple, reliable, and efficient software.
Please visit [this link](https://go.dev/doc/install) for installation guide.

#### Foundry 
Foundry is a blazing fast, portable and modular toolkit for Ethereum application development written in Rust.
Please find the installation and usage guide [here](https://book.getfoundry.sh/getting-started/installation).

#### Goloop CLI
Goloop CLI is a command-line interface tool designed for managing nodes, executing transactions, deploying smart contracts, 
and interacting with the ICON blockchain network.
Run the following command to install Goloop CLI.
```
go install github.com/icon-project/goloop/cmd/goloop@latest
```
More about Goloop can be found [here](https://docs.icon.community/concepts/computational-utilities/goloop).  
                                             

### Centralized Relay Installation
- Clone the repository:
    ```
    git@github.com:icon-project/centralised-relay.git
    ```
- Build and install the binary(coding standards and style conventions are also checked):
    ```
    make all
    ```                
  Please run ```centralized-rly``` to verify the installation.

### Build and deploy the contracts


### Transfer messages
  