# Deployment Checklist

Follow these steps to deploy your application:

## 1. Machine Setup

- [ ] Create EC2 instance
  - [ ] Ubuntu 20.04
  - [ ] t2.medium
  - [ ] 10 GB
  - [ ] 1 GB swap space
  - [ ] 1 GB RAM
  - [ ] 2 vCPUs
  - [ ] 1 static IP
  - [ ] `sudo apt update && sudo apt upgrade -y`

- [ ] Install the necessary software:
  - [ ] GO
    - [ ] `sudo apt install go`
  - [ ] Git
    - [ ] `sudo apt install git`
  - [ ] Make
    - [ ] `sudo apt install make`

- [ ] Clone the repository.
  - [ ] `git clone git@github.com:icon-project/centralised-relay.git`
- [ ] Build the application.
  - [ ] `make install`
- [ ] Create a systemd service file.
  - [ ] `sudo nano /etc/systemd/system/centralized-rly.service`
  - [ ] Add the following:

    ```
    [Unit]
    Description=Centralized Relay
    After=network.target

    [Service]
    Type=simple
    Restart=always
    RestartSec=5s
    ExecStart=/usr/bin/centralized-rly start

    [Install]
    WantedBy=multi-user.target
    ```

## 2. Configuring Config Files

- [ ] Update the application config file with the correct settings.
  - [ ] `centralized-rly config init`
  - [ ] RPC urls
  - [ ] Chain NIDs
- [ ] Set the environment variables.
- [ ] Check the database connection.

## 3. AWS KMS Related Auth

- [ ] Set up an IAM EC2 role with the necessary permissions.
  - [ ] `aws iam create-role --role-name <role-name> --assume-role-policy-document file://assume-role.json`
  file contents:

  ```
  {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Sid": "",
        "Effect": "Allow",
        "Principal": {
          "Service": "ec2.amazonaws.com"
        },
        "Action": "sts:AssumeRole"
      }
    ]
  }
  ```

  - [ ] `aws iam create-policy --policy-name <policy-name> --policy-document file://kms-policy.json`
  file contents:

  ```
  {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Sid": "Stmt1234567890",
        "Effect": "Allow",
        "Action": [
          "kms:Decrypt",
          "kms:Encrypt",
          "kms:GenerateDataKey",
          "kms:DescribeKey",
          "kms:ListKeys",
          "kms:ReEncrypt*",
          "kms:CreateGrant",
          "kms:ListGrants",
          "kms:RevokeGrant",
          "kms:GenerateRandom"
        ],
        "Resource": "*"
      }
    ]
  }
  ```

  - [ ] `aws iam attach-role-policy --role-name <role-name> --policy-arn <policy-arn>`

  - [ ] Create a KMS key.
    - [ ] `centralized-rly keystore init`

## 4. Wallets and Funds

- [ ] Create or import a wallet.
  - [ ] `centralized-rly keystore create --chain icon --password <password>`
  - [ ] `centralized-rly keystore create --chain avalanche --password <password>`

- [ ] Ensure the wallet has sufficient funds.
  - [ ] `centralized-rly keys list --chain icon`
  - [ ] `centralized-rly keys list --chain avalanche`
  - [ ] Check balances on explorer.
- [ ] Set the wallet address in the config file.
  - [ ] `centralized-rly keystore use --chain icon --address <address>`
  - [ ] `centralized-rly keystore use --chain avalanche --address <address>`

## 5. Final Checks

- [ ] Run the application in a test environment.
- [ ] Check the logs for any errors.
- [ ] If everything is working as expected, deploy the application to the production environment.
