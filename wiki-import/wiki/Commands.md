# In-Game Commands

Commands are typed in the in-game chat. The prefix is configurable in `config.json` (default: `!`). Each command can be enabled/disabled individually in the config.

Operators (users with the appropriate `rights` bitmask) can use disabled commands.

> **Looking for a player-friendly guide?** See [Player Commands](Player-Commands) for usage examples and troubleshooting.

## General Commands

| Command | Description | Default |
|---------|-------------|---------|
| `!help` | Show enabled commands | Enabled |
| `!playtime` | Show account playtime | Enabled |
| `!timer` | Toggle quest timer display | Enabled |
| `!reload` | Reload all players in current Land (forces client refresh) | Enabled |
| `!psn <username>` | Link a PSN ID to your account | Enabled |
| `!discord` | Generate a token for Discord account linking (use with `/link` in Discord) | Enabled |
| `!course <name>` | Toggle a subscription course on/off | Enabled |

## Admin Commands

| Command | Description | Default |
|---------|-------------|---------|
| `!rights <integer>` | Set the rights bitmask on your account | Disabled |
| `!tele <x> <y>` | Teleport to coordinates | Disabled |
| `!kqf get` | Show current Key Quest Flag hex value (G10+ only) | Disabled |
| `!kqf set <hex>` | Set Key Quest Flags (16-char hex, requires Land switch to apply) | Disabled |
| `!ban <character_id> [duration]` | Ban a user, optionally with duration | Disabled |

### Ban Duration Format

Duration uses `<number><unit>`:

| Unit | Meaning |
|------|---------|
| `s` | Seconds |
| `m` | Minutes |
| `h` | Hours |
| `d` | Days |
| `mo` | Months |
| `y` | Years |

Example: `!ban 12345 30d` bans character 12345 for 30 days. Omitting the duration bans permanently.

## Raviente Commands

These control the Raviente large-scale raid. Some subcommands require the ZZ client.

| Command | Description |
|---------|-------------|
| `!ravi start` | Start Great Slaying event |
| `!ravi cm` | Check current damage multiplier |
| `!ravi ss` | Send sedation support (ZZ only) |
| `!ravi sr` | Send resurrection support (ZZ only) |
| `!ravi rs` | Request sedation support (ZZ only) |

## Special Chat

| Input | Effect |
|-------|--------|
| `@dice` | Roll a random number 1-100 |

## Security Classification

### High-Risk (disable in production)

| Command | Risk |
|---------|------|
| `rights` | Grants admin privileges and all courses |
| `kqf` | Unlocks restricted HR progression content |
| `ban` | Can permanently remove players |
| `tele` | Bypasses normal movement restrictions |

### Medium-Risk (use with caution)

| Command | Risk | Mitigation |
|---------|------|------------|
| `course` | Players toggle courses without payment | Limit which courses are available via `Courses` in [config](Erupe-Configuration#courses-subscription-types) (set `Enabled: false` to lock) |

### Safe

| Command | Reason |
|---------|--------|
| `help` | Read-only, shows available commands |
| `reload` | Only affects visual state, no persistent changes |
| `ravi` | Only affects raid event state |
| `timer` | User preference only |
| `playtime` | Read-only display |
| `psn` | Account linking, low risk |
| `discord` | Token generation, requires Discord bot |

## Configuration

Commands are defined in the `Commands` array in `config.json`. Each entry has three fields:

| Property | Type | Description |
|----------|------|-------------|
| `Name` | string | Command identifier (must match internal name) |
| `Enabled` | boolean | Whether the command is active for non-operators |
| `Prefix` | string | Chat text that triggers the command |

### Example: Minimal Safe Config

```json
"Commands": [
  { "Name": "help", "Enabled": true, "Prefix": "help" },
  { "Name": "reload", "Enabled": true, "Prefix": "reload" },
  { "Name": "ravi", "Enabled": true, "Prefix": "ravi" },
  { "Name": "timer", "Enabled": true, "Prefix": "timer" },
  { "Name": "playtime", "Enabled": true, "Prefix": "playtime" }
]
```

### Example: Community Server

```json
"Commands": [
  { "Name": "help", "Enabled": true, "Prefix": "help" },
  { "Name": "reload", "Enabled": true, "Prefix": "reload" },
  { "Name": "ravi", "Enabled": true, "Prefix": "ravi" },
  { "Name": "course", "Enabled": true, "Prefix": "course" },
  { "Name": "timer", "Enabled": true, "Prefix": "timer" },
  { "Name": "playtime", "Enabled": true, "Prefix": "playtime" },
  { "Name": "psn", "Enabled": true, "Prefix": "psn" },
  { "Name": "discord", "Enabled": true, "Prefix": "discord" },
  { "Name": "rights", "Enabled": false, "Prefix": "rights" },
  { "Name": "tele", "Enabled": false, "Prefix": "tele" },
  { "Name": "kqf", "Enabled": false, "Prefix": "kqf" },
  { "Name": "ban", "Enabled": false, "Prefix": "ban" }
]
```

### Example: Full Development Config

```json
"Commands": [
  { "Name": "help", "Enabled": true, "Prefix": "help" },
  { "Name": "reload", "Enabled": true, "Prefix": "reload" },
  { "Name": "ravi", "Enabled": true, "Prefix": "ravi" },
  { "Name": "course", "Enabled": true, "Prefix": "course" },
  { "Name": "rights", "Enabled": true, "Prefix": "rights" },
  { "Name": "tele", "Enabled": true, "Prefix": "tele" },
  { "Name": "kqf", "Enabled": true, "Prefix": "kqf" },
  { "Name": "ban", "Enabled": true, "Prefix": "ban" },
  { "Name": "timer", "Enabled": true, "Prefix": "timer" },
  { "Name": "playtime", "Enabled": true, "Prefix": "playtime" },
  { "Name": "psn", "Enabled": true, "Prefix": "psn" },
  { "Name": "discord", "Enabled": true, "Prefix": "discord" }
]
```
