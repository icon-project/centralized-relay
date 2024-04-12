# Keystore Command

The `keystore` command is used to manage your Ethereum keystore.

## Usage

```bash
centralized-relay keystore <command> [options]

Commands:
  create    Create a keystore
  import    Import a keystore
  use       Use a keystore

Options:
  -h, --help   Show help
```

## Commands

### `create`

```bash
keystore create [flags]

Flags:
  -c, --chain string       The chain for which to create the keystore
  -p, --password string    The password to encrypt the keystore
```

### `import`

The `import` command is used to import a KMS key into the relay keystore.

```bash
keystore import [flags]

Flags:
  -c, --chain string           The chain for which to import the keystore
  -k, --keystore string        The path to the keystore
  -p, --password string        The password to encrypt the keystore
```

### `use`

The `use` command is used to set a specific keystore as the active one.

_Warning: contract call: the method changes the relayer address._

```bash
keystore use [flags]

Flags:
  -a, --address string         The address of the keystore to use
  -c, --chain string           The chain for which to use the keystore
```

## Examples

### Create a keystore

```bash
centralized-relay keystore create --chain=0x2.icon --password=12345678
```

### Import a keystore

```bash
centralized-relay keystore import --chain=0x2.icon --keystore=./keystore.json --password=12345678
```

### Use a keystore

```bash
centralized-relay keystore use --chain=0x2.icon --address=0x1234567890
```
