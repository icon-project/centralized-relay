# End-to-End Testing Setup and Demo

## Prerequisites

To run the demo, the following software needs to be installed.

* Docker compose \[[download](https://docs.docker.com/compose/install/)\]

## Setting up the Environment

1. Create an `e2e-tests` folder and clone the `centralised-relay` repository:

    ```bash
    mkdir e2e-tests
    cd e2e-tests
    git clone https://github.com/icon-project/centralised-relay.git
    cd centralised-relay
    make build-docker
    cd -  # Back to the root folder
    ```

2. Build an `icon-chain` image

   ```bash
    git clone https://github.com/icon-project/goloop.git 
    cd goloop
    make gochain-icon-image
    cd -  # Back to the root folder
   ```

## Running e2e Tests

To conduct tests for the IBC integration system, follow these steps:

#### 1. Configure Environment Variables

Before initiating the tests, configure essential environment variables:

- **`TEST_CONFIG_PATH`**: Set this variable to the absolute path of your chosen configuration file. You can create these configuration files using the sample files provided in the `centralised-relay` source folder. Sample configuration files are available at the following locations:
    - sample config : `centralised-relay/test/testsuite/sample-config.yaml`

Here's an example of environment variable configuration:

```bash
export TEST_CONFIG_PATH=/path/to/config.yaml
```

ℹ️ Please note that most of the config content can be used same as it in sample config however you may need to update the image name and version for Archway, Neutron, and Icon in the configuration file you create.


After configuring these variables, navigate to the `centralised-relay` source folder:

```bash
cd centralised-relay
```

#### 2. Run the Test Script

Use the appropriate command to run the test suite. Depending on your specific testing requirements, you can use the following command:

```bash
./scripts/execute-test.sh [options]
```

Replace `[options]` with any command-line options or arguments that the test script supports. Here's an option block to help you:

```markdown
Options:
 --clean: Clean contract directories (true/false, default: false).
 --build-xcall: Build xCall contracts (true/false, default: false).
 --xcall-branch <branch>: Specify the xCall branch to build (default: main).
 --use-docker: Use Docker for building contracts(true/false, default: false).
 --test <test_type>: Specify the type of test (e2e, default: e2e).
```

To perform an end-to-end (e2e) test with all the necessary builds, execute the following command:
```bash
./scripts/execute-test.sh --build-xcall --use-docker --test e2e
```
This command covers building IBC and xCall contracts while utilizing Docker and running an end-to-end test.

Once you've initially built the contracts using the command above, you can easily execute the e2e test by using the following simplified command:
```bash
./scripts/execute-test.sh  --test e2e
```
