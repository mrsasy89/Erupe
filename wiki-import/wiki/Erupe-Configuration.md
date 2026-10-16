# Erupe Config File
### General Configuration

| Variable               | Description                                                                                                         | Default      | Options                                                                               |
|------------------------|---------------------------------------------------------------------------------------------------------------------|--------------|---------------------------------------------------------------------------------------|
| Host                   | The IP or host address to allow connections from                                                                    | 127.0.0.1    |                                                                                       |
| BinPath                | The bin path folder is where you place files needed for various parts of the game such as scenario and quest files  | bin          |                                                                                       |
| ContentPath            | Directory of `<table>/*.json` [content files](Database-Setup#content-files) synchronised into the database on start and on `POST /v2/admin/content/reload`; empty = `<BinPath>/content`| ""           |                                                                                       |
| Language               | Sets the language of server-side strings. Only English `en` and Japanese `jp` are available, contributions welcome  | en           | en/jp                                                                                 |
| DisableSoftCrash       | Disables some crash outputs and in-game shutdown timer                                                              | false        |                                                                                       |
| HideLoginNotice        | Hides the notices that appear on login from `LoginNotices`                                                          | true         |                                                                                       |
| LoginNotices           | Array of MHFML-formatted login notices (see below)                                                                  |              |                                                                                       |
| PatchServerManifest    | Bare host (no scheme) the sign server hands to the original launcher for `/mhf_file.php`; see [Client Distribution](Client-Distribution) |              | Required for PS3 clients                                                              |
| PatchServerFile        | Bare host (no scheme) the original launcher fetches `/mhfdat/…` from; usually the same as above                       |              | Required for PS3 clients                                                              |
| ScreenshotAPIURL       | See [ZeruLight/Servers](https://github.com/ZeruLight/Servers#mhf-screenshot-bbs-api)                                |              |                                                                                       |
| DeleteOnSaveCorruption | Flags characters with corrupted savedata as deleted, can be used as punishment for cheaters                         | false        |                                                                                       |
| ClientMode             | Changes the targeted Client Version, see [Client Versions](https://github.com/mrsasy89/Erupe/wiki/Client-Versions) | ZZ           |                                                                                       |
| QuestCacheExpiry       | Number of seconds decrypted Event Quest data will be cached                                                         | 300          |                                                                                       |
| CommandPrefix          | Sets the in-game Chat Command prefix                                                                                | !            |                                                                                       |
| AutoCreateAccount      | Allows users that don't exist to auto create their account from initial login                                       | true         |                                                                                       |
| DefaultCourses         | Sets the Courses that will always be enabled by their ID                                                            | 1, 23, 24    | See [Enumerations](https://github.com/mrsasy89/Erupe/wiki/Enumerations#course-types) |
| EarthStatus            | Enables Pallone Fest, Tower and Conquest War events                                                                 | 0            |1=Conquest 2=Conquest, 11=Pallone, 12 = Pallone Reward, 21=Tower                                                      |
| EarthID                | A random event ID                                                                                                   | 0            |                                                                                       |
| EarthMonsters          | Sets the IDs of the monsters targeted in the Conquest War                                                           | 0            |                                                                                       |
| SaveDumps.Enabled      | Enables save dumps to a folder that is set at `SaveDumps.OutputDir`                                                 | true         |                                                                                       |
| SaveDumps.RawEnabled   | Enables saving an uncompressed copy of the savedata                                                                 | false        |                                                                                       |
| SaveDumps.OutputDir    | The folder that save dumps are saved to                                                                             | save-backups |                                                                                       |

#### MHFML Tags (Login Notice Formatting)

Login notices use MHFML (Monster Hunter Frontier Markup Language) tags:

| Tag | Effect |
|-----|--------|
| `<BODY>` | Start new text block |
| `<BR>` | Line break |
| `<CENTER>`, `<LEFT>`, `<RIGHT>` | Text alignment |
| `<SIZE_2>` | Normal text size |
| `<SIZE_3>` | Large text size |
| `<C_4>` | Red text |
| `<C_5>` | Yellow text |
| `<C_7>` | White text |

**Example notice:**
```json
"LoginNotices": [
  "<BODY><CENTER><SIZE_3><C_4>Welcome to Our Server!<BR><BODY><LEFT><SIZE_2><C_5>Join our Discord: discord.gg/example"
]
```

### `DebugOptions` Configuration

| Variable            | Description                                                                                 | Default | Options    |
|---------------------|---------------------------------------------------------------------------------------------|---------|------------|
| CleanDB             | Deletes all data in the `users`, `characters`, `guilds` tables permanently                  | false   |            |
| MaxLauncherHR       | Always return HR7 on login to allow you to bypass World HR requirements                     | false   |            |
| LogInboundMessages  | Enables logging inbound messages to stdout                                                  | false   |            |
| LogOutboundMessages | Enables logging outbound messages to stdout                                                 | false   |            |
| LogMessageData      | Enables logging the data transferred as a hexdump                                           | false   |            |
| MaxHexdumpLength    | Maximum number of hex bytes that will be dumped to stdout                                   | 256     |            |
| DivaOverride        | Overrides the Diva Defense event status in game                                             | 0       | 0/1/2/3/-1 |
| FestaOverride       | Overrides the Hunter Festival event status in game                                          | -1      | 0/1/2/3/-1 |
| TournamentOverride  | Overrides the Hunter Tournament event status in game                                        | 0       | 0/1/2/3/-1 |
| DisableTokenCheck   | This disables the random token that is generated at login from being checked, very insecure | false   |            |
| QuestTools          | Enables various quest debug logs                                                            | false   |            |
| ProxyPort           | See [rockisch/mhf-dev-proxy](https://github.com/rockisch/mhf-dev-proxy)                     | 0       |            |
| CapLink.Values      | Sets values that are checked by the game client                                             |         |            |
| CapLink.Key         | A string secure key                                                                         |         |            |
| CapLink.Host        | The host address of the CapLink API                                                         |         |            |
| CapLink.Port        | The port of the CapLink API                                                                 | 80      |            |

### `GameplayOptions` Configuration

| Variable                       | Description                                                                              | Default | Options         |
|--------------------------------|------------------------------------------------------------------------------------------|---------|-----------------|
| FeaturedWeapons                | Number of Active Feature weapons to generate daily                                       | 1       |                 |
| MaximumNP                      | Maximum number of N Points held by a player                                              | 100000  |                 |
| MaximumRP                      | Maximum number of Ryoudan (Clan) Points held by a player                                 | 50000   |                 |
| MaximumFP                      | Maximum number of Festa Points held by a player                                          | 120000  |                 |                                                           
| TreasureHuntExpiry             | The number of seconds until a Clan Treasure Hunt expires                                 | 604800  |                 |
| DisableLoginBoost              | Disables the Login Boost system                                                          | false   |                 |
| DisableBoostTime               | Disables the daily NetCafe Boost Time                                                    | false   |                 |
| BoostTimeDuration              | The number of seconds NetCafe Boost Time lasts for                                       | 7200    |                 |
| ClanMealDuration               | The number of seconds a Clan Meal can be activated for after cooking                     | 3600    |                 |
| ClanMemberLimits               | Sets the max Clan Members after hitting Rank thresholds (Max. 100)                       |         | [Rank, Members] |
| BonusQuestAllowance            | Number of Bonus Point Quests to allow daily                                              | 3       |                 |
| DailyQuestAllowance            | Number of Daily Quests to allow daily                                                    | 1       |                 |
| LowLatencyRaviente             | Experimental mode to instantly update global Raviente hitpoints                          | false   |                 |
| RegularRavienteMaxPlayers      | Changes the maximum players able to sign up to fight Regular Raviente (HR2)              | 8       |                 |
| ViolentRavienteMaxPlayers      | Changes the maximum players able to sign up to fight Violent Raviente (HR5)              | 8       |                 |
| BerserkRavienteMaxPlayers      | Changes the maximum players able to sign up to fight Berserk Raviente (G1)               | 32      |                 |
| ExtremeRavienteMaxPlayers      | Changes the maximum players able to sign up to fight Extreme Raviente (G1+G50 weapon)    | 32      |                 |
| SmallBerserkRavienteMaxPlayers | Changes the maximum players able to sign up to fight Small Berserk Raviente (G1)         | 8       |                 |
| GUrgentRate                    | Adjusts the rate of G Urgent quests spawning                                             | 0.10    |                 |
| GCPMultiplier                  | Adjusts the multiplier of GCP rewarded for quest completion                              | 1.00    |                 |
| HRPMultiplier                  | Adjusts the multiplier of Hunter Rank Points rewarded for quest completion               | 1.00    |                 |
| HRPMultiplierNC                | Adjusts the multiplier of Hunter Rank Points rewarded for quest completion in a NetCafe  | 1.00    |                 |
| SRPMultiplier                  | Adjusts the multiplier of Skill Rank Points rewarded for quest completion                | 1.00    |                 |
| SRPMultiplierNC                | Adjusts the multiplier of Skill Rank Points rewarded for quest completion in a NetCafe   | 1.00    |                 |
| GRPMultiplier                  | Adjusts the multiplier of G Rank Points rewarded for quest completion                    | 1.00    |                 |
| GRPMultiplierNC                | Adjusts the multiplier of G Rank Points rewarded for quest completion in a NetCafe       | 1.00    |                 |
| GSRPMultiplier                 | Adjusts the multiplier of G Skill Rank Points rewarded for quest completion              | 1.00    |                 |
| GSRPMultiplierNC               | Adjusts the multiplier of G Skill Rank Points rewarded for quest completion in a NetCafe | 1.00    |                 |
| ZennyMultiplier                | Adjusts the multiplier of Zenny rewarded for quest completion                            | 1.00    |                 |
| ZennyMultiplierNC              | Adjusts the multiplier of Zenny rewarded for quest completion in a NetCafe               | 1.00    |                 |
| GZennyMultiplier               | Adjusts the multiplier of G Zenny rewarded for quest completion                          | 1.00    |                 |
| GZennyMultiplierNC             | Adjusts the multiplier of G Zenny rewarded for quest completion in a NetCafe             | 1.00    |                 |
| MaterialMultiplier             | Adjusts the multiplier of Monster Materials rewarded for quest completion                | 1.00    |                 |
| MaterialMultiplierNC           | Adjusts the multiplier of Monster Materials rewarded for quest completion in a NetCafe   | 1.00    |                 |
| GMaterialMultiplier            | Adjusts the multiplier of G Monster Materials rewarded for quest completion              | 1.00    |                 |
| GMaterialMultiplierNC          | Adjusts the multiplier of G Monster Materials rewarded for quest completion in a NetCafe | 1.00    |                 |
| ExtraCarves                    | Grant n extra chances to carve ALL carcasses                                             | 0       |                 |
| ExtraCarvesNC                  | Grant n extra chances to carve ALL carcasses in a NetCafe                                | 0       |                 |
| GExtraCarves                   | Grant n extra chances to carve ALL G Rank carcasses                                      | 0       |                 |
| GExtraCarvesNC                 | Grant n extra chances to carve ALL G Rank carcasses in a NetCafe                         | 0       |                 |
| DisableHunterNavi              | Disables the Hunter Navi                                                                 | false   |                 |
| MezFesSoloTickets              | Number of solo tickets given weekly                                                      | 5       |                 |
| MezFesGroupTickets             | Number of group tickets given weekly                                                     | 1       |                 |
| MezFesDuration                 | The duration of MezFes in seconds (start time = end of week - duration)                  | 172800  |                 |
| MezFesSwitchMinigame           | Switches the group minigame from Volpakkun Together to Tokotoko Partnya                  | false   |                 |
| EnableKaijiEvent               | Enabled the Kaiji event in the Rasta Bar (G10 only)                                      | false   |                 |
| EnableHiganjimaEvent           | Enables the Higanjima event in the Rasta Bar                                             | false   |                 |
| EnableNierEvent                | Enables the Nier event in the Rasta Bar                                                  | false   |                 |
| DisableRoad                    | Disables the Hunting Road                                                                | false   |                 |
| SeasonOverride                 | Enables quests to use the current Season and time                                        | false   |                 |

### `API` Configuration

The REST API (default port 8080) used by mhf-outpost, web sign-up pages and save export. Endpoints are listed in [Server Operations → API Server](Server-Operations#api-server).

| Config Item         | Description                                                                                          | Default | Options |
|---------------------|------------------------------------------------------------------------------------------------------|---------|---------|
| Enabled             | Start the API server                                                                                 | true    |         |
| Port                | Listening port                                                                                       | 8080    |         |
| PatchServer         | Origin (`http://host[/subdir]`) of the game-file tree, returned as `patchServer` by `POST /v2/login` so mhf-outpost can sync; see [Client Distribution](Client-Distribution) |         | Leave empty if you host no files |
| PatchTree.Enabled   | Serve the game-file tree from the API port (`/mhf_file.php`, `/mhfdat/`) and regenerate its manifest on startup; sets `PatchServer` to `http://<Host>:<Port>` if empty | false   |         |
| PatchTree.Root      | Directory holding `mhfdat/{exe,dat}` (and the generated `key.txt`)                                   | patch   |         |
| Banners             | `[{"src", "link"}]` images shown by launcher UIs (`GET /launcher`)                                    | []      |         |
| Messages            | `[{"message", "date", "kind", "link"}]` notices shown by launcher UIs (`kind`: 0 default, 1 new)    | []      |         |
| Links               | `[{"name", "icon", "link"}]` shown by launcher UIs                                                   | []      |         |
| LandingPage.Enabled | Serve a browser landing page at `/`                                                                  | true    |         |
| LandingPage.Title   | Landing page title                                                                                   |         |         |
| LandingPage.Content | Landing page body, raw HTML allowed                                                                  |         |         |

### `Logging` Configuration

> **Branch availability:** File logging is available on the `v9.2.x-stable` branch. The `main` branch logs to stdout only.

When enabled, Erupe uses [lumberjack](https://github.com/natefinch/lumberjack) for automatic log rotation and compression.

| Setting | Type | Default | Description |
|---------|------|---------|-------------|
| `LogToFile` | boolean | `true` | Enable file logging (logs to both console and file) |
| `LogFilePath` | string | `"logs/erupe.log"` | Path to log file (directory created automatically) |
| `LogMaxSize` | number | `100` | Maximum log file size in MB before rotation |
| `LogMaxBackups` | number | `3` | Number of old log files to keep |
| `LogMaxAge` | number | `28` | Maximum days to retain old logs |
| `LogCompress` | boolean | `true` | Compress rotated log files with gzip |

#### Log Rotation

When the current log reaches `LogMaxSize` MB:

1. Current log is renamed to `erupe.log.YYYY-MM-DD-HH-MM-SS`
2. If `LogCompress` is `true`, old log is compressed to `.gz`
3. New `erupe.log` is created
4. Logs beyond `LogMaxBackups` count or `LogMaxAge` days are deleted

#### Log Format

The format depends on `DevMode`:

- **`DevMode: true`** — human-readable console format: `2025-11-18T10:30:45.123Z INFO channelserver Player connected {"charID": 12345}`
- **`DevMode: false`** — JSON format: `{"level":"info","ts":"...","logger":"channelserver","msg":"Player connected","charID":12345}`

#### Disk Space Estimate

```
Total ≈ (LogMaxSize × LogMaxBackups) / CompressionRatio
```

Text logs compress ~10:1, JSON logs ~8:1. Default settings (100 MB × 3 backups, compressed) use roughly **30 MB** of disk.

#### Examples

**Standard production:**
```json
"Logging": {
  "LogToFile": true,
  "LogFilePath": "logs/erupe.log",
  "LogMaxSize": 100,
  "LogMaxBackups": 7,
  "LogMaxAge": 30,
  "LogCompress": true
}
```

**Development (no file logging):**
```json
"Logging": {
  "LogToFile": false
}
```

**High-volume server:**
```json
"Logging": {
  "LogToFile": true,
  "LogFilePath": "/var/log/erupe/erupe.log",
  "LogMaxSize": 500,
  "LogMaxBackups": 14,
  "LogMaxAge": 60,
  "LogCompress": true
}
```

### Discord
There is limited Discord capability in Erupe. The feature allows you to relay messages from your server into a channel.
This may be either be removed or revamped in a future version.

### Commands
There are several chat commands that can be turned on and off. Most are for admin or debugging purposes.

| Name     | command       | Description                                    | Default  | Options                                                                               |
|----------|---------------|------------------------------------------------|----------|---------------------------------------------------------------------------------------|
| Help     | help          | Shows enabled Chat Commands                    | enabled  |                                                                                       |
| Rights   | rights VALUE  | Sets the rights integer for your account       | disabled |                                                                                       |
| Teleport | tele X,Y      | Teleports the user to specific x,y coordinates | disabled |                                                                                       |
| Reload   | reload        | Reloads all users and character objects        | enabled  |                                                                                       |
| KeyQuest | kqf FLAGS     | Sets the Key Quest Flag for your character     | disabled |                                                                                       |
| Course   | course OPTION | Enables/Disables a course for your account     | enabled  | See [Enumerations](https://github.com/mrsasy89/Erupe/wiki/Enumerations#course-types) |
| PSN      | psn USERNAME  | Links the specified PSN to your account        | enabled  |                                                                                       |

### Ravi Sub Commands
| Name     | command                        | Description                   |
|----------|--------------------------------|-------------------------------|
| Raviente | ravi start                     | Starts Ravi Event             |
| Raviente | ravi cm / ravi checkmultiplier | Checks Ravi Damage Multiplier |
| Raviente | ravi ss                        | Send Sedation Support         |
| Raviente | ravi sr                        | Send Resurrection Support     |
| Raviente | ravi rs                        | Request Sedation Support      |


### World `Entries` Configuration

| Config Item        | Description                  | Options                                                    |
|--------------------|------------------------------|------------------------------------------------------------|
| Type               | Server type                  | 1=Normal, 2=Cities, 3=Newbie, 4=Tavern, 5=Return, 6=MezFes |
| Recommended        | Changes quest availability   | See below table                                            |
| AllowedClientFlags | Sets clients allowed to join | See below table                                            |

### `Recommended` 
* 0 = All quests
* 1 = Up to 2☆ Quests  
* 2 = Up to 4☆ star Quests 
* 4 = All Quests in HR (Enables G Experience Tab) 
* 5 = Only G Rank Quests 
* 6 = Minigame World

### `AllowedClientFlags`
* 0 = All clients
* 2 = Wii U
* 64 = Blank
* 512 = PC
* 4096 = PC,PS3/PS4
* 4098 = PS3/PS4
* 8258 = PS Vita
* 12288 = PC/PS3(R)

### Courses (Subscription Types)

See [Enumerations](https://github.com/mrsasy89/Erupe/wiki/Enumerations#course-types) for the full course type list.

Courses are stored as a bitfield in `users.rights`. Players toggle available courses with the [`!course` command](Player-Commands#course).

#### Bitfield Encoding

The bit position equals the course ID (`value = 2^ID`), matching [Enumerations](Enumerations#course-types).

| Course | ID (Bit) | Value |
|--------|----------|-------|
| Trial | 1 | 2 |
| HunterLife | 2 | 4 |
| Extra | 3 | 8 |
| ExtraB | 4 | 16 |
| Mobile | 5 | 32 |
| Premium | 6 | 64 |
| Pallone | 7 | 128 |
| Assist | 8 | 256 |
| N | 9 | 512 |
| Hiden | 10 | 1024 |
| HunterSupport | 11 | 2048 |
| NBoost | 12 | 4096 |
| NetCafe | 26 | 67108864 |
| HLRenewing | 27 | 134217728 |
| EXRenewing | 28 | 268435456 |

Add values together for combinations. Example: HunterLife + Extra + NetCafe = 4 + 8 + 67108864 = **67108876**.

#### Manual Database Override

```sql
-- Enable core gameplay courses (Trial + HunterLife + Extra + Premium + Assist)
UPDATE users SET rights = 334 WHERE username = 'player';

-- Enable only HunterLife + Extra + NetCafe
UPDATE users SET rights = 67108876 WHERE username = 'player';
```

#### Example Configs

**Open server (all courses available):**
```json
"DefaultCourses": [1, 2, 3, 4, 5, 6, 7, 8, 9, 23, 24]
```

**Restricted server (core only):**
```json
"DefaultCourses": [1, 2, 23, 24]
```

**Free-to-play (NetCafe only):**
```json
"DefaultCourses": [9, 23, 24]
```

# Database Configuration Schemas
The following can be edited via the database for gameplay customization.
- Distributions
- Certain Shop Types (5/6/7/8/10)
- Event Quests
- Frontier Point Items
- Festival
- Gacha
- Net-Cafe
- [Scenarios](https://github.com/mrsasy89/Erupe/wiki/Scenarios)

There are examples of these in `server/migrations/seed/`, loaded automatically on a fresh database (see [Database Setup](Database-Setup#default-content-seeds)).

# Production Security Checklist

Before running a public-facing server, verify:

- [ ] `DebugOptions.CleanDB: false` — prevents wiping the database on startup
- [ ] `DebugOptions.DisableTokenCheck: false` — enables session token validation
- [ ] `DebugOptions.MaxLauncherHR: false` — prevents bypassing HR requirements
- [ ] `AutoCreateAccount: false` — require manual account creation (or use the API)
- [ ] Strong `Database.Password` set (not blank, not default)
- [ ] `Rights` command disabled in `Commands` array (or restricted to operators)
- [ ] `KeyQuest` command disabled (unless intentional)
- [ ] `Ban` command disabled for non-operators
- [ ] Firewall configured for only necessary ports (see [Server Operations](Server-Operations#firewall-configuration))
- [ ] Database not exposed publicly (bind PostgreSQL to `localhost` or internal network)
- [ ] `Logging` enabled for monitoring and incident response
- [ ] `SaveDumps.Enabled: true` for character backup recovery