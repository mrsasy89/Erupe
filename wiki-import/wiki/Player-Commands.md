# Player Commands Reference

This guide covers all chat commands available to players on Erupe servers.

## Using Commands

Type commands in the chat window. All commands start with `!` (your server may use a different prefix).

---

## General Commands

### !help

Display available commands.

**Usage:** `!help`

Shows a list of all enabled commands on your server.

### !reload

Reload all players in your current area.

**Usage:** `!reload`

Forces a refresh of all player data in your Land/Stage. Useful if you see visual glitches or invisible players.

### !course

Toggle subscription courses on your account.

**Usage:** `!course <course_name>`

Available courses (server-dependent). Names are case-insensitive; aliases shown in parentheses.

| ID | Course Name | Aliases | Description |
|----|-------------|---------|-------------|
| 1 | `trial` | `tl` | Trial Course — limited free trial access |
| 2 | `hunterlife` | `hl` | Hunter Life Course — **base subscription required to play the game** |
| 3 | `extra` | `extraa`, `ex` | Extra Course — additional storage, extra character slots, My Series sets |
| 5 | `mobile` | | Mobile Course — mobile-linked features |
| 6 | `premium` | | Premium Course — Poogie outfits, gallery, cosmetic features |
| 7 | `pallone` | `extrac` | Pallone Course — Pallone/balloon features |
| 8 | `assist` | `legend`, `rasta` | Assist Course — hire Legendary Rasta (partner NPCs) for quests |
| 9 | `n` | | N Course — Xbox/NetCafe time-based access |
| 10 | `hiden` | `secret` | Hiden Course — increases the amount of 魂 obtained from Hardcore Quests，relaxes Skill Rank penalties/requirements and increases Rank point/Zenny gains |
| 11 | `huntersupport` | `hunteraid`, `support`, `aid`, `royal` | Hunter Support Course — Hunter Navigation and support features |
| 12 | `nboost` | `netcafeboost`, `boost` | N Boost Course — point and reward multipliers |
| 26 | `netcafe` | `cafe`, `official` | NetCafe Course — official Net Cafe boost time and bonus rewards |
| 27 | `hlrenewing` | `hlr`, `cardhl` | HL Renewing — Hunter Life auto-renewal |
| 28 | `exrenewing` | `exr`, `cardex` | EX Renewing — Extra Course auto-renewal |

Example:
```
!course hunterlife    → Toggles Hunter Life Course on/off
!course hl            → Same thing (alias)
!course assist        → Toggles Assist/Rasta Course
```

Available courses depend on server configuration (`Courses` in config.json). If a course is not enabled on the server, you'll see a "locked" message. See [Enumerations](Enumerations#course-types) for the full list of course IDs.

### !timer

Toggle the quest timer display.

**Usage:** `!timer`

Shows or hides the quest countdown timer during hunts. Your preference is saved to your account.

### !playtime

Show your total playtime.

**Usage:** `!playtime`

Displays how long you've played on your current character.

---

## Account Linking

### !psn

Link a PlayStation Network ID to your account.

**Usage:** `!psn <your_psn_id>`

Example: `!psn MyPSNUsername`

### !discord

Generate a Discord account linking token.

**Usage:** `!discord`

1. Type `!discord` in-game
2. Copy the token shown
3. Use `/link <token>` in your server's Discord bot

---

## Raviente Siege Commands

These commands are used during the Great Slaying (Raviente) siege event.

| Command | Description |
|---------|-------------|
| `!ravi start` | Start the Great Slaying event |
| `!ravi cm` | Check current damage multiplier |
| `!ravi sr` | Send resurrection support |
| `!ravi ss` | Send sedation support |
| `!ravi rs` | Request sedation support from others |

Commands only work while participating in a Raviente siege.

---

## Special Chat

| Input | Effect |
|-------|--------|
| `@dice` | Roll a random number 1-100, visible to all players in your area |

This is not a `!` command — just type `@dice` in regular chat.

---

## Admin Commands

These commands are disabled by default and require operator permissions (`rights` bitmask). Server administrators may enable some for staff use.

| Command | Description |
|---------|-------------|
| `!tele <x> <y>` | Teleport to coordinates |
| `!rights <value>` | Set account permission level |
| `!kqf get` | View current Key Quest Flag value |
| `!kqf set <hex>` | Set KQF to specific 16-char hex value (takes effect after switching Land) |
| `!ban <character_id> [duration]` | Ban a player (optionally with duration) |

### Ban Duration Format

| Unit | Meaning |
|------|---------|
| `s` | Seconds |
| `m` | Minutes |
| `h` | Hours |
| `d` | Days |
| `mo` | Months |
| `y` | Years |

Examples: `!ban ABC123 7d` (7-day ban), `!ban ABC123` (permanent).

---

## Troubleshooting

| Problem | Solution |
|---------|----------|
| Command not recognized | Check spelling; try `!help` to see what's available |
| "Command is disabled" | The server operator has disabled this command |
| Command does nothing | Make sure you're in the correct chat channel; some commands only work in specific situations |
| Course says "locked" | That course is not available on this server — contact your server admin |

---

## For Server Operators

See [Commands](Commands) for configuration details and security classification.
