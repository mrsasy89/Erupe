# Database Setup

Erupe uses PostgreSQL and manages the schema itself. There is nothing to
restore or apply by hand: the binary carries its migrations and default
content, and runs them on startup.

## How the schema is managed

```
server/migrations/
├── sql/    0001_init.sql, 0002_….sql, …   schema migrations, embedded in the binary
└── seed/   <table>/*.json, *.sql          default content, loaded once on a fresh database
```

On every start Erupe:

1. Creates the `schema_version` table if it is missing.
2. Detects a database that predates the migration system (tables present, no
   `schema_version` rows) and marks the `0001_init.sql` baseline as applied, so
   an install from before 9.2 upgrades in place.
3. Applies every `sql/NNNN_*.sql` not yet recorded in `schema_version`, in
   numerical order, each in its own transaction. A failing migration stops
   the server with the error in the log; nothing after it runs.
4. If the database was **empty** before this start, loads every seed file
   (shops, event quests, gacha, Festa/Diva/Tournament defaults, campaigns).
   Seeds are not tracked in `schema_version` and are not re-run on later
   starts — an existing database keeps whatever content its operator has
   curated.
5. If a `game-data/content/` directory exists, synchronises its
   [content files](#content-files) into the database. This is how an
   operator keeps shop rows and the like as files instead of SQL; without
   the directory nothing happens.

The version number is simply the highest applied migration. `0001` is the
consolidated baseline (formerly the `init.sql` dump plus the `patch-schema/`
series); numbers are not necessarily contiguous.

### Fresh install

```bash
createdb -U postgres erupe      # or let the setup wizard do it
./erupe-ce                      # applies 0001…, then the seeds
```

Log lines to expect on the first run:

```
Applying migration 0001: 0001_init.sql
…
Database: Applied N migration(s), now at version V
Database: Applied 17 seed data file(s)
```

### Upgrading

Stop the server, replace the binary, start it. Pending migrations are
applied before it accepts connections. Back up first (below): migrations
can rewrite data, and there is no automatic rollback.

To see where a database stands:

```sql
SELECT version, filename, applied_at FROM schema_version ORDER BY version;
```

### Adding a migration (developers)

Create `server/migrations/sql/NNNN_short_description.sql` with the next free
number. It must be plain SQL that PostgreSQL can run in one transaction, and
should be idempotent where practical (`IF NOT EXISTS`, `ON CONFLICT`). Schema
changes belong here; **content does not** — a pricing tweak or a new shop row
shipped as a migration reaches every operator's live database and can only be
undone with another migration (see `0024`/`0026` for the cautionary example).
Default content goes in `seed/`, which only affects new installs; an
operator's own content goes in [content files](#content-files).

## Default content (seeds)

Plain tabular content is JSON, one directory per table
(`seed/<table>/*.json`); seeds that need real SQL (subqueries, procedural
logic) stay `.sql`.

| File | Content |
|------|---------|
| `event_quests/default.json` | Event quest rotation (~290 entries) |
| `shop_items/road.json` | Hunting Road shop (`shop_type` 10) |
| `shop_items/diva.json` | Diva Defense shop (`shop_type` 8) |
| `shop_items/other.json` | GCP and gacha-coin exchange shops (`shop_type` 5, 6, 7) |
| `fpoint_items/default.json` | Frontier Point exchange list |
| `diva_prizes/default.json` | Diva Defense prize ladder |
| `festa_prizes/default.json`, `festa_trials/default.json` | Festa prize tiers and trial goals |
| `campaigns/`, `campaign_categories/`, `campaign_category_links/`, `campaign_rewards/` | Example campaigns |
| `cafebonus/default.json` | Net Cafe bonuses |
| `scenario_counter/default.json` | Scenario counter initialisation |
| `TournamentDefaults.sql` | Hunting Tournament defaults |
| `GachaDemo.sql` | Three example gachas (normal, step-up, box) |
| `DistributionDemo.sql` | Example gift-box distributions |

Each JSON file starts with a `comment` explaining what every column means,
read from the handler that consumes it.

They are demo-quality: enough to see every system working, not a curated
live-server economy. Operators edit the content afterwards — as
[content files](#content-files), through the
[operator API](Server-Operations#operator-endpoints-v2admin) for
distributions and events, or directly in the tables; all three survive
upgrades because seeds never re-run.

## Content files

Content files are the operator's counterpart of the seeds: the same JSON
format, in a directory Erupe reads on every start and on request, and
**synchronised** rather than inserted once. The file is the source of
truth for the rows it covers — edit it, reload, and the table matches.

```
game-data/
└── content/
    └── shop_items/
        └── road.json      ← a copy of seed/shop_items/road.json, edited
```

The directory is `ContentPath` in `config.json`, `game-data/content/` (next
to the quest data) when unset. It does not have to exist; an install
without it behaves exactly as before.

### Starting from a seed

```bash
mkdir -p game-data/content/shop_items
cp server/migrations/seed/shop_items/road.json game-data/content/shop_items/
$EDITOR game-data/content/shop_items/road.json
```

Then either restart, or reload without downtime:

```bash
curl -s -X POST http://localhost:8080/v2/admin/content/reload -H "Authorization: Bearer $TOKEN"
# {"dir":"game-data/content","files":[{"path":"shop_items/road.json","table":"shop_items","inserted":1,"updated":3,"deleted":0}]}
```

Shop rows are read per request, so the change is visible in game at once.

### Format

```json
{
  "comment": "Free text; what the columns mean.",
  "key":   ["shop_type", "shop_id", "item_id"],
  "scope": {"shop_type": 10},
  "rows": [
    {"shop_id": 4, "item_id": 11664, "cost": 20000, "quantity": 1, "min_hr": 0, "min_sr": 0, "min_gr": 1, "store_level": 1, "max_quantity": 0, "road_floors": 0, "road_fatalis": 0}
  ]
}
```

- The **table** is the directory name. A file can hold rows for one table.
- **`key`** (required) names the columns that identify a row — a natural key
  such as shop + item, never the serial `id`, which the database assigns.
  An existing row with the same key is updated in place and keeps its `id`,
  so per-character purchase counters (`shop_items_bought`) survive a price
  change. A row with a new key is inserted. Two rows with the same key are
  an error.
- **`scope`** (optional) says which rows the file *owns*, as column → value
  or column → list of values. Rows in scope that the file no longer lists
  are **deleted**; a row outside the scope is refused. Rows may leave a
  scalar scope column out (it is filled in). `{}` means the whole table.
  Without `scope`, the file only inserts and updates.
- Every row must have the same columns. Values are JSON; `{"raw": "NOW()"}`
  passes a SQL expression through.
- `truncate` and `onConflict` are seed-only and ignored here, so a seed
  file copies over unchanged.

Each file is applied in its own transaction. A malformed file stops the
run (the server keeps running on what the database holds; the reload
endpoint answers `422` with the file name). Several files may target the
same table with different scopes, e.g. `road.json` for `shop_type` 10 and
`event.json` for a temporary tab.

### Which tables

Any table whose rows are plain data works; the shipped seeds declare a
`key` for `shop_items`, `fpoint_items`, `diva_prizes`, `festa_prizes`,
`festa_trials`, `event_quests`, `cafebonus`, `scenario_counter` and the
campaign tables. Distributions and event cycles are better driven through
the [operator API](Server-Operations#operator-endpoints-v2admin), which
knows their side tables.

## Key Tables

### Users & Characters

| Table | Purpose |
|-------|---------|
| `users` | Accounts (bcrypt password, rights bitmask, PSN ID, Discord token) |
| `characters` | Character data with ~17 bytea savedata fields |
| `sign_sessions` | Active authentication tokens |
| `bans` | Ban records with optional expiry |

### Game Content

| Table | Purpose |
|-------|---------|
| `distribution` + `distribution_items` | Gift box system (login bonuses, GM gifts) |
| `event_quests` | Event quest rotation (quest_type, quest_id, cycling timestamps) |
| `gacha_shop` + `gacha_entries` + `gacha_items` | Gacha lottery (normal, step-up, box) |
| `shop_items` | Shop inventory (cost, HR/SR/GR requirements) |
| `fpoint_items` | Frontier Point shop |
| `cafebonus` | Net Cafe time-based bonuses |
| `scenario_counter` | Scenario/cutscene unlock tracking |

### Guilds (8 tables)

| Table | Purpose |
|-------|---------|
| `guilds` | Guild info (name, leader, RP, icon, recruiting) |
| `guild_characters` | Membership and roles |
| `guild_applications` | Pending invites/applications |
| `guild_posts` | Bulletin board |
| `guild_adventures` | Guild adventure tracking |
| `guild_hunts` | Guild hunt records |
| `guild_meals` | Guild meal cooking |
| `guild_alliances` | Inter-guild alliances |

### Events

| Table | Purpose |
|-------|---------|
| `events` | Active events (type: festa/diva/vs/mezfes) |
| `festa_registrations` | Festa sign-ups |
| `festa_trials` + `festa_prizes` | Festa scores and rewards |
| `feature_weapon` | Daily featured weapon rotation |
| `tower` | Tower mode progress (TR, TRP, TSP) |
| `rengoku_score` | Rengoku survival scores |
| `kill_logs` | Monster kill tracking |

## Common Admin Operations

Most of these have an API equivalent since 9.5 — see
[Server Operations → Operator Endpoints](Server-Operations#operator-endpoints-v2admin).
The SQL remains valid for scripts and one-offs.


### Add an event quest
```sql
INSERT INTO event_quests (max_players, quest_type, quest_id, mark)
VALUES (4, 28, 23045, 2);
```
See [Enumerations](Enumerations#event-quests) for quest type and mark values.

### Create a gift distribution
```sql
INSERT INTO distribution (type, deadline, event_name, description, times_acceptable)
VALUES (1, '2099-12-31', 'Welcome Gift', 'A gift for new hunters', 1)
RETURNING id;

INSERT INTO distribution_items (distribution_id, item_type, item_id, quantity)
VALUES (<returned_id>, 7, 482, 10);
```
See [Enumerations](Enumerations#distribution-types) for type/item_type values.

### Ban a user
```sql
INSERT INTO bans (user_id, expires_at)
VALUES (<user_id>, NOW() + INTERVAL '7 days');
```

### Reset Raviente
```sql
SELECT ravireset(0);
```

### Check active sessions
```sql
SELECT u.username, s.token FROM sign_sessions s
JOIN users u ON u.id = s.user_id;
```

## Backup & Restore

### Backup

```bash
# Full backup (custom format, recommended)
pg_dump -U postgres -Fc erupe > erupe_backup.dump

# Plain SQL backup
pg_dump -U postgres erupe > erupe_backup.sql

# Compressed plain SQL
pg_dump -U postgres erupe | gzip > erupe_backup.sql.gz
```

### Restore

```bash
# From custom format
dropdb -U postgres erupe
createdb -U postgres erupe
pg_restore -U postgres -d erupe erupe_backup.dump

# From plain SQL
dropdb -U postgres erupe
createdb -U postgres erupe
psql -U postgres -d erupe -f erupe_backup.sql

# From compressed SQL
gunzip -c erupe_backup.sql.gz | psql -U postgres -d erupe
```

### Maintenance

```bash
# Reclaim space and update query planner statistics
psql -U postgres -d erupe -c "VACUUM ANALYZE;"

# Check database size
psql -U postgres -d erupe -c "SELECT pg_size_pretty(pg_database_size('erupe'));"
```

## Troubleshooting

| Error | Cause | Fix |
|-------|-------|-----|
| `Connection refused` | PostgreSQL not running or wrong port | Verify service is running: `systemctl status postgresql`. Check `Port` in config. |
| `password authentication failed` | Wrong password or user | Verify `Database.Password` in `config.json`. Check user exists: `psql -U postgres -c "\du"` |
| `database "erupe" does not exist` | Database not created | Run `createdb -U postgres erupe` |
| `permission denied for table` | User lacks privileges | `GRANT ALL ON ALL TABLES IN SCHEMA public TO <user>;` |
| `Database migration failed: applying NNNN_….sql` | A migration hit an error (usually a manual schema edit that conflicts with it) | Read the SQL error in the log, fix the conflicting object, restart; the migration retries |
| Server starts but shops/quests are empty | Database was not empty on first start, so seeds were skipped | Load the wanted `server/migrations/seed/*.sql` files with `psql` |
| Old install: "Detected existing database without schema_version tracking" | Expected once when upgrading from before 9.2 | Nothing to do |

## Performance Tuning

For servers with many concurrent players, consider tuning PostgreSQL:

```conf
# /etc/postgresql/*/main/postgresql.conf

shared_buffers = 2GB          # 25% of available RAM
work_mem = 16MB               # Per-operation memory
maintenance_work_mem = 512MB  # For VACUUM, CREATE INDEX
```

For high connection counts, consider [PgBouncer](https://www.pgbouncer.org/) for connection pooling. Restart PostgreSQL after changing `postgresql.conf`.
