# DB Command

The `db` command is used to interact with the database in the centralized relay.

The database is a levelDB database that is used to store the states of the relay.

The database stores the following information:

- Messages received from the chain that are not yet processed and processing.
- The relay failed messages for retries until configured attempts.
- The flagged messages that needs manual intervention.
- The last block that was processed for all the configured chains

## Usage

```bash
centralized-relay db [command] [flags]
```

## Commands

### List all the messages in the database

```bash
messages list  [flags]

Flags:
  -c, --chain   string      Chain ID
  -p, --page    int         Page number
  -l, --limit   int         Page limit
```

### Relay a message manually

```bash
messages relay [flags]

Flags:
  -c, --chain   string      Chain ID
  -s, --sn      int         Sequence number
  -h, --height  int         Block height [optional: fetch messages from chain]
```

### Remove a message from the database

```bash
messages remove [flags]

Flags:
  -c, --chain   string        Chain ID
  -s, --sn      int           Sequence number
```

### Prune the database

```bash
prune [flags]

Flags:

--db-path
```

### Revert Message

```bash
message revert [flags]

Flags:
  -c, --chain   string        Chain ID
  -s, --sn      int           Sequence number
```

## Examples

1. **List all the messages in the database.**

```bash
centralized-relay db messages list --chain 0x2.icon
```

2. **Relay a message from the database manually on chain.**

```bash
centralized-relay db messages relay --chain 0x2.icon --sn 1 --height 100
```

3. **Relay a message from the database manually from database.**

```bash
centralized-relay db messages relay --chain 0xa869.fuji --sn 1
```

4. **Remove a message from the database.**

```bash
centralized-relay db messages remove --chain 0x2.icon --sn 1
```

5. **Revert a message.**

```bash
centralized-relay db messages revert --chain 0x2.icon --sn 1
```

6. **Prune the database.**

```bash
centralized-relay db prune
```
