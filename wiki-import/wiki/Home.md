## What is Erupe?

Erupe is a community-maintained server emulator for Monster Hunter Frontier, written in Go. It implements the sign, entrance, and channel servers needed to play the game after Capcom shut down the official servers in December 2019. All code is clean-room reverse-engineered from protocol analysis.

Most core systems are functional: quests, guilds, multiplayer, shops/gacha, mail, character saves, and housing. Some endgame event systems (tournaments, tower) have minimal implementations. See [Client Versions](Client-Versions) for supported game versions (Season 6.0 through ZZ).

## Getting Started

### Prerequisites

- **PostgreSQL** (14+ recommended) — the only thing Erupe needs installed. Go 1.25+ only if you build from source.
- **Quest/scenario data** — see [Server Operations](Server-Operations#quest--scenario-files); `binsync` fetches it.
- Players need a game client; see [Client Distribution](Client-Distribution).

### 1. Get Erupe

Download `erupe-ce` (Linux) or `erupe.exe` (Windows) from the
[latest release](https://github.com/Houmgaor/Erupe/releases/latest), or build it:

```bash
go build -o erupe-ce
```

### 2. First run

Run the binary with no `config.json` next to it and it opens a **setup wizard**
at `http://localhost:8080`: it asks for your PostgreSQL credentials, creates
the `erupe` database if it does not exist, writes `config.json`, and can fetch
the quest data ("Sync Now"). `./erupe-ce --setup` reopens it later.

Prefer files? Create the database yourself and copy the example config:

```bash
createdb -U postgres erupe
cp config.example.json config.json   # set Host, Database.*, ClientMode
```

There is no SQL to run by hand. On startup Erupe applies its embedded schema
migrations, and on a brand-new database also loads the default content (shops,
event quests, gacha, Festa/Diva/Tournament defaults). See
[Database Setup](Database-Setup) for how that works and what to back up.

### 3. Configure

At minimum, in `config.json`:
- `Host` — your public IP or hostname (`127.0.0.1` for local play)
- `Database` — PostgreSQL connection details
- `ClientMode` — must match the game client version your players use (default `ZZ`)

See [Erupe Configuration](Erupe-Configuration) for the full reference.

### 4. Run

```bash
./erupe-ce
```

### Docker alternative

From the `docker/` directory (see `docker/README.md`):

```bash
cp ../config.example.json config.json   # Database.Host = "db", Database.Password = "password"
docker compose up                        # PostgreSQL + pgAdmin (http://localhost:5050) + Erupe
```

The schema is applied automatically on first start, the same way as a native install.

## Wiki Pages

### Setup & Operations
- **[Database Setup](Database-Setup)** — schema management, table reference, bundled demo data, common admin SQL
- **[Server Operations](Server-Operations)** — quest/scenario files, save system, REST API, Discord bot, Docker
- **[Client Distribution](Client-Distribution)** — what to host so players can connect: mhf-outpost vs original launcher, patch tree (game-file sync), launcher page, ports
- **[Erupe Configuration](Erupe-Configuration)** — full config.json reference (general, debug, gameplay, worlds, courses)
- **[Commands](Commands)** — in-game chat commands (general, admin, Raviente)
- **[Player Commands](Player-Commands)** — player-friendly command reference with usage examples

### Development
- **[Contributing](Contributing)** — adding packet handlers, debugging, testing, CI

### Reference
- **[Enumerations](Enumerations)** — distribution types, item types, course types, quest types and flags
- **[Client Versions](Client-Versions)** — all MHF versions from Season 1.0 to Z3.1 with dates, platforms, and ClientMode IDs
- **[Scenarios](Scenarios)** — scenario quest database with IDs and dates (Forward through end-of-service)
- **[Hard Core Monsters Schedule Map](Hard-Core-Monsters-Schedule-Map)** — HC/GHC monster schedule data from mhfsch.bin


## Local quest-type corrections

See [Enumerations — Quest Types](Enumerations#quest-types) and [Quest Type Client Mapping](Quest-Type-Client-Mapping) for the client-verified category names, including **54 → Gear Acquisition Quests** and **9 → Promotional**.

Documentation copied from [Houmgaor/Erupe wiki](https://github.com/Houmgaor/Erupe/wiki), snapshot `69ade90d207aacd7cb674b0feb4ca8b266e22e70`. Original documentation credit remains with its contributors. Local quest-type corrections are documented separately and do not imply all upstream instructions have been revalidated against this fork.
