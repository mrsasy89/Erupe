# Server Operations

## Quest & Scenario Files

The `bin/` directory (configured via `BinPath` in config.json) holds game data files the server sends to clients.

```
bin/
├── quests/              # Quest files (.bin or .json)
├── scenarios/           # Scenario files (.bin or .json)
├── events/              # Organized quest archives (not read by server)
├── rengoku_data.bin     # Hunting Road config (.bin or .json)
└── save_override.bin    # Optional: forces savedata values on character creation
```

### JSON Format Support

All three file types support a human-readable `.json` alternative. The server always tries `.bin` first and falls back to `.json` automatically — existing binary files work unchanged.

**Quest** (`bin/quests/<name>.json`): covers all binary sections — quest text (UTF-8, converted to Shift-JIS on the wire), objectives, monster spawns, reward tables, supply box, stages, map sections, gathering points, and more. Use `ParseQuestBinary` to convert an existing `.bin` to JSON.

**Scenario** (`bin/scenarios/<name>.json`): supports sub-header chunks with UTF-8 strings (flags 0x01/0x02), inline episode listings (flag 0x08), and opaque base64 blobs for JKR-compressed chunks (flags 0x10/0x20). See `docs/scenario-format.md` in the Erupe repository for the full format reference.

**Hunting Road** (`bin/rengoku_data.json`): define multi-road and solo-road floors, spawn tables, and monster IDs in plain JSON. The server assembles and ECD-encrypts the binary at startup.

### Quest File Naming

Files in `quests/` follow this pattern:

```
[questID][time][season].bin   (or .json)
```

- **questID**: 5-digit hex ID (e.g., `23045`)
- **time**: `d` = day, `n` = night
- **season**: `0` = warm, `1` = cold, `2` = breeding

Example: `23045d1.bin` is quest 23045, daytime, cold season.

If `SeasonOverride` is enabled in the config, the server converts quest filenames to match the current in-game season automatically.

### Scenario File Naming

Files in `scenarios/` follow this pattern:

```
[CategoryID]_0_0_0_S[MainID]_T[Flags]_C[ChapterID].bin   (or .json)
```

Missing scenario files will crash the client — make sure all referenced scenarios have corresponding files.

### Quest Caching

Decrypted quest data is cached in memory. The cache expires after `QuestCacheExpiry` seconds (default: 300). Encrypted quest files are auto-decrypted on load.

### Quest Backporting

If `ClientMode` is set to Z1 or earlier, quest data is automatically backported for compatibility with older clients.

## Hunting Tournament

The hunting tournament system (`公式狩猟大会`) hosts periodic competitive events where players submit timed hunts and fishing runs. The tournament lobby, cup listings, leaderboards, registration, and run submission are all server-driven.

### Lifecycle phases

Each tournament row has four Unix timestamps that drive the phase shown to the client:

| Phase | State byte | Description |
|-------|:---:|-------------|
| Before start | 0 | No active tournament shown |
| Registration open | 1 | Players can register (`EntryTournament`) |
| Hunting active | 2 | Players can enter and submit runs |
| Ranking / reward | 3 | Results visible; rewards claimable |

### Database tables

| Table | Purpose |
|-------|---------|
| `tournaments` | One row per tournament instance (schedule) |
| `tournament_cups` | Competition categories (speed hunt, fishing, etc.) linked to a tournament |
| `tournament_sub_events` | Specific hunt/fish targets grouped by `cup_group` |
| `tournament_entries` | Player registrations |
| `tournament_results` | Submitted run records (submission order = rank proxy) |

### Seed data

`TournamentDefaults.sql` seeds cup and sub-event data from live tournament #150 (Brachydios time-attack + fishing cups). A demo tournament starting immediately is also inserted so the lobby is populated on a fresh install.

### Scheduling a tournament

Update the demo row or insert a new one:

```sql
-- Update the existing demo tournament
UPDATE tournaments SET
    name        = '第151回公式狩猟大会',
    start_time  = EXTRACT(epoch FROM '2026-04-01 14:00:00+09'::timestamptz)::bigint,
    entry_end   = EXTRACT(epoch FROM '2026-04-04 14:00:00+09'::timestamptz)::bigint,
    ranking_end = EXTRACT(epoch FROM '2026-04-13 14:00:00+09'::timestamptz)::bigint,
    reward_end  = EXTRACT(epoch FROM '2026-04-20 14:00:00+09'::timestamptz)::bigint
WHERE id = 1;
```

To add new cup types or hunt targets, insert rows into `tournament_cups` (with the correct `tournament_id`) and `tournament_sub_events` (linked by `cup_group`).

### Known limitation

The `Unk2` field in `MsgMhfEnterTournamentQuest` is stored as `event_id` in `tournament_results`. Its exact semantics are unconfirmed — if leaderboards appear empty in-game, this is the likely cause. See [#184](https://github.com/Mezeporta/Erupe/issues/184).

## Save Data

### Character Saves

Character data is stored in the `characters` table across ~17 bytea columns: `savedata`, `decomyset`, `hunternavi`, `otomoairou`, `partner`, `platebox`, `platedata`, `platemyset`, `rengokudata`, `savemercenary`, `minidata`, `gacha_items`, `house_info`, `login_boost`, `skin_hist`, `scenariodata`, `savefavoritequest`.

Save data offsets vary by game version (the server reads different offsets depending on `ClientMode`).

### Save Backups

Configured in `config.json`:

```json
"SaveDumps": {
  "Enabled": true,
  "RawEnabled": false,
  "OutputDir": "save-backups"
}
```

When enabled, every save operation writes a backup to `save-backups/<characterID>/`. Set `RawEnabled` to also dump the uncompressed save data.

### Save Corruption

If `DeleteOnSaveCorruption` is enabled, characters with corrupt save data are flagged as deleted. When disabled (default), the character remains in the database but may be unplayable.

## API Server

The optional REST API (default port 8080) provides launcher integration and account management.

### Launcher Endpoints

| Method | Path | Purpose |
|--------|------|---------|
| GET | `/launcher` | Launcher splash page (banners, messages, links) |
| POST | `/login` | Authenticate (username/password → token + character list) |
| POST | `/register` | Create account |
| POST | `/character/create` | Create character (requires token) |
| POST | `/character/delete` | Delete character (requires token + charId) |
| POST | `/character/export` | Export character save data |

### Operator Endpoints (`/v2/admin`)

Everything an operator used to do with `psql`, behind the API. Sign in
normally (`POST /v2/login`) with an account whose `users.op` flag is set and
use the returned token; non-operators get `403`. Every change is logged with
the operator's user ID.

| Method | Path | Purpose |
|--------|------|---------|
| GET | `/v2/admin/users?q=&limit=` | Search accounts (rights, op, ban state) |
| GET | `/v2/admin/users/{id}` | Account detail with characters |
| PATCH | `/v2/admin/users/{id}` | `{"op": bool}` and/or `{"rights": n}`; you cannot remove your own op flag |
| POST | `/v2/admin/users/{id}/password` | `{"password": "…"}` reset |
| POST | `/v2/admin/users/{id}/ban` | Optional `{"expires": "<RFC 3339>"}`; permanent otherwise |
| DELETE | `/v2/admin/users/{id}/ban` | Unban |
| GET / POST | `/v2/admin/notices` | Login notices shown after the static `LoginNotices`; `{"body", "expires"?}`; no restart needed |
| DELETE | `/v2/admin/notices/{id}` | Remove a notice |
| GET / POST | `/v2/admin/distributions` | Gift-box distributions; `POST` takes the header fields plus `items: [{itemType, itemId, quantity}]` |
| GET / DELETE | `/v2/admin/distributions/{id}` | Detail with items / delete |
| GET | `/v2/admin/events` | Rows of the `events` table |
| PUT | `/v2/admin/events/{festa\|diva\|vs\|mezfes}` | Start or restart a cycle at `{"startTime"}` (default now) |
| DELETE | `/v2/admin/events/{type}` | Stop (remove the rows) |
| POST | `/v2/admin/content/reload` | Synchronise the [content files](Database-Setup#content-files) (`game-data/content/`) into the database without a restart |

Full request/response shapes: `docs/openapi.yaml` in the repository.

**First operator.** Nothing can grant the first `op` flag but the database:

```sql
UPDATE users SET op = true WHERE username = 'your-account';
```

After that, further operators are promoted with `PATCH /v2/admin/users/{id}`.

Example — announce a maintenance and restart Diva Defense tonight at midnight JST:

```bash
TOKEN=$(curl -s -X POST http://localhost:8080/v2/login \
  -H 'Content-Type: application/json' -d '{"username":"admin","password":"…"}' | jq -r .user.token)

curl -s -X POST http://localhost:8080/v2/admin/notices -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"body":"~C05Maintenance tonight 02:00–03:00 CEST","expires":"2026-09-20T03:00:00+02:00"}'

curl -s -X PUT http://localhost:8080/v2/admin/events/diva -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' -d '{"startTime":"2026-09-19T15:00:00Z"}'
```

### Screenshot Endpoints

| Method | Path | Purpose |
|--------|------|---------|
| POST | `/api/ss/bbs/upload.php` | Upload screenshot (JPEG, requires token) |
| GET | `/api/ss/bbs/{id}` | Retrieve screenshot |

Screenshot config:
```json
"Screenshots": {
  "Enabled": true,
  "OutputDir": "screenshots",
  "UploadQuality": 100
}
```

## Discord Integration

### Setup

1. Create a bot on the [Discord Developer Portal](https://discord.com/developers/applications) with the `applications.commands` scope
2. Configure in `config.json`:

```json
"Discord": {
  "Enabled": true,
  "BotToken": "YOUR_BOT_TOKEN",
  "RelayChannel": {
    "Enabled": true,
    "MaxMessageLength": 183,
    "RelayChannelID": "CHANNEL_ID"
  }
}
```

### Slash Commands

Registered automatically when the bot starts:

| Command | Description |
|---------|-------------|
| `/link <token>` | Link Discord account to Erupe account (token from `!discord` in-game) |
| `/password <password>` | Change password (requires linked account) |

### Chat Relay

When `RelayChannel` is enabled:
- **Game → Discord**: World chat and stage chat are relayed to the configured Discord channel
- **Discord → Game**: Messages from the Discord channel appear in-game with a `[D]` prefix

Messages are normalized: Discord mentions become `@username`, emoji become `:emoji_name:`, and non-ASCII characters are stripped. Long messages are split into 61-character chunks to fit the in-game chat.

## Docker

### Services

| Service | Port | Purpose |
|---------|------|---------|
| `db` | 5432 | PostgreSQL (auto-inits schema on first run via `setup.sh`) |
| `pgadmin` | 5050 | Database admin UI (`user@pgadmin.com` / `password`) |
| `server` | 53312, 8080, 53310, 54001-54008 | Erupe |

### Workflow

```bash
cd docker/
docker compose up db pgadmin       # 1. Start DB (auto-applies schemas on first run)
# Configure config.json (use "db" as database host)
# Place quest/scenario files in bin/
docker compose up server           # 2. Start Erupe
```

### Data Persistence

- Database: `docker/db-data/`
- Save data: `docker/savedata/`
- Full wipe: `docker compose down` then delete these directories

### Troubleshooting

If PostgreSQL doesn't initialize properly, check that `docker/init/setup.sh` has LF line endings (not CRLF).

## Firewall Configuration

Open the sign, entrance, and channel server ports for incoming TCP connections. Adjust the channel port range to match your `Entrance.Entries` configuration.

### Linux (ufw)

```bash
sudo ufw allow 53312/tcp    # Sign server
sudo ufw allow 53310/tcp    # Entrance server
sudo ufw allow 54001:54010/tcp  # Channel servers
```

### Linux (iptables)

```bash
sudo iptables -A INPUT -p tcp --dport 53312 -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 53310 -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 54001:54010 -j ACCEPT
```

### Windows PowerShell

```powershell
New-NetFirewallRule -DisplayName "Erupe Sign" -Direction Inbound -Protocol TCP -LocalPort 53312 -Action Allow
New-NetFirewallRule -DisplayName "Erupe Entrance" -Direction Inbound -Protocol TCP -LocalPort 53310 -Action Allow
New-NetFirewallRule -DisplayName "Erupe Channels" -Direction Inbound -Protocol TCP -LocalPort 54001-54010 -Action Allow
```

## Monitoring

Quick one-liners to check server health:

```bash
# Check if all servers are listening
ss -tlnp | grep -E '(53312|53310|54001)'

# Count established connections per channel port
netstat -an | grep ESTABLISHED | grep -c 54001

# Watch connections in real time (refreshes every second)
watch -n 1 'ss -tn state established | grep -E "54[0-9]{3}" | wc -l'
```
