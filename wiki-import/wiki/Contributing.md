# Contributing

This guide covers the development workflows for working on the Erupe codebase.

## Project Directory Guide

```
Erupe/
├── main.go                  # Entry point — starts all servers, handles shutdown
├── config.example.json      # Template config (copy to config.json)
│
├── config/                  # Config loading (Viper-backed)
│   └── config.go            #   Main Config struct, game version modes (S1–ZZ)
│
├── network/                 # Protocol layer
│   ├── packetid.go          #   Packet ID enum (~200+ opcodes, `go generate` for stringer)
│   ├── crypt_conn.go        #   Blowfish-encrypted TCP connection
│   ├── crypt_packet.go      #   Packet framing and checksum
│   ├── mhfpacket/           #   ~400 packet structs (one file per message type)
│   │   └── msg_mhf_*.go     #     Parse/Build for each packet
│   ├── binpacket/           #   Binary relay packets (chat, mail, targeted)
│   ├── crypto/              #   Low-level Blowfish cipher
│   └── clientctx/           #   Per-connection context
│
├── server/
│   ├── channelserver/       # Gameplay server (by far the largest)
│   │   ├── handlers_table.go    # PacketID → handler dispatch (~200 entries)
│   │   ├── handlers_quest.go    # Quest system
│   │   ├── handlers_guild.go    # Guild operations (14 tables)
│   │   ├── handlers_stage.go    # Multiplayer rooms/lobbies
│   │   ├── handlers_cast_binary.go  # Real-time state relay (position, animation)
│   │   ├── handlers_mail.go     # In-game mail
│   │   ├── handlers_shop.go     # Shops and gacha
│   │   ├── handlers_*.go        # ... ~50 handler files by game system
│   │   ├── sys_session.go       # Per-connection state (character, stage, send queue)
│   │   ├── sys_stage.go         # Stage lifecycle, player sync
│   │   ├── sys_semaphore.go     # Distributed locks (Raviente, guild ops)
│   │   ├── sys_channel_server.go    # Server lifecycle, Raviente shared state
│   │   └── compression/         # Save data compression
│   ├── signserver/          # Authentication (TCP 53312)
│   ├── entranceserver/      # World list & character select (TCP 53310)
│   ├── api/                 # REST API (port 8080) — launcher, screenshots
│   └── discordbot/          # Discord bot — slash commands, chat relay
│
├── common/                  # Shared utility libraries
│   ├── byteframe/           #   Binary read/write (big-endian)
│   ├── bfutil/              #   Blowfish helpers
│   ├── decryption/          #   Game file decryption
│   ├── gametime/            #   MHF time ↔ real time conversion
│   ├── mhfcourse/           #   Course/subscription definitions
│   ├── mhfitem/             #   Item ID lookups
│   ├── mhfmon/              #   Monster ID lookups
│   ├── mhfcid/              #   Character ID handling
│   ├── token/               #   Session token generation
│   ├── pascalstring/        #   Length-prefixed string parsing
│   ├── stringstack/         #   Stack for stage movement history
│   └── stringsupport/       #   String utilities
│
├── server/migrations/       # PostgreSQL schema, embedded in the binary and applied on startup
│   ├── sql/                 #   0001_init.sql baseline + numbered migrations
│   └── seed/                #   Default content loaded once on a fresh database
│
├── game-data/               # Quest/scenario/road data (formerly bin/; fetched by binsync)
│   ├── quests/              #   Quest files (.bin or .json)
│   ├── scenarios/           #   Cutscene/scenario files
│   └── events/              #   Event quest archives
│
├── docker/                  # Docker Compose setup (PostgreSQL, pgAdmin, server)
├── tools/loganalyzer/       # Log analysis utility
├── vendor/                  # Vendored Go dependencies
└── .github/                 # CI workflows (test, race, coverage, build)
```

### Where to look first

- **Implementing a feature?** Start with the `handlers_*.go` file for that game system, then check the packet structs in `network/mhfpacket/`.
- **Fixing a bug?** Enable packet logging (see [Debugging](#debugging)) and trace from `handlers_table.go` to the relevant handler.
- **Adding a DB table?** Add `server/migrations/sql/NNNN_description.sql` with the next free number; it is applied automatically. See [Database Setup](Database-Setup#adding-a-migration-developers).
- **Understanding binary formats?** `common/byteframe/` is the serialization layer used everywhere.

## Adding a Packet Handler

This is the most common type of contribution. A new packet handler touches 4 files and optionally a 5th for code generation.

### 1. Define the packet struct

Create `network/mhfpacket/msg_mhf_your_packet.go`. Every packet implements the `MHFPacket` interface (Parse, Build, Opcode):

```go
package mhfpacket

import (
    "errors"
    "erupe-ce/common/byteframe"
    "erupe-ce/network"
    "erupe-ce/network/clientctx"
)

type MsgMhfYourPacket struct {
    AckHandle uint32
    SomeField uint16
}

func (m *MsgMhfYourPacket) Opcode() network.PacketID {
    return network.MSG_MHF_YOUR_PACKET
}

func (m *MsgMhfYourPacket) Parse(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
    m.AckHandle = bf.ReadUint32()
    m.SomeField = bf.ReadUint16()
    return nil
}

func (m *MsgMhfYourPacket) Build(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
    return errors.New("NOT IMPLEMENTED")
}
```

`Parse` reads fields from the binary frame in the exact order the client sends them. `Build` is the inverse — only implement it if the server needs to construct this packet type (most server→client packets reuse `MsgSysAck`).

The `AckHandle` field is present on most packets. The client sends it and expects the same value back in the ack response.

### 2. Register the packet in the opcode switch

Add a case to `network/mhfpacket/opcode_to_packet.go`:

```go
case network.MSG_MHF_YOUR_PACKET:
    return &MsgMhfYourPacket{}
```

### 3. Register and implement the handler

Add to `server/channelserver/handlers_table.go`:

```go
handlerTable[network.MSG_MHF_YOUR_PACKET] = handleMsgMhfYourPacket
```

Then implement the handler in the appropriate `handlers_*.go` file. The handler signature is always `func(s *Session, p mhfpacket.MHFPacket)`:

```go
func handleMsgMhfYourPacket(s *Session, p mhfpacket.MHFPacket) {
    pkt := p.(*mhfpacket.MsgMhfYourPacket)

    // Do work (DB queries, state changes, etc.)
    result, err := s.server.db.Query("SELECT ...")
    if err != nil {
        doAckSimpleFail(s, pkt.AckHandle, make([]byte, 4))
        return
    }

    bf := byteframe.NewByteFrame()
    bf.WriteUint32(someValue)
    doAckSimpleSucceed(s, pkt.AckHandle, bf.Data())
}
```

### 4. Add the packet ID constant

If this is a genuinely new packet (not an existing `reserve` slot), add it to `network/packetid.go` and regenerate the stringer:

```bash
cd Erupe/network && go generate
```

This regenerates `packetid_string.go` so debug logs show the packet name instead of a number.

### ACK Response Types

Every handler that receives an `AckHandle` must respond with one of these (defined in `handlers.go`):

| Function | When to use |
|----------|-------------|
| `doAckSimpleSucceed(s, ackHandle, data)` | Single operation succeeded |
| `doAckSimpleFail(s, ackHandle, data)` | Single operation failed |
| `doAckBufSucceed(s, ackHandle, data)` | Returning a list/buffer of results |
| `doAckBufFail(s, ackHandle, data)` | List/buffer query failed |
| `doAckEarthSucceed(s, ackHandle, frames)` | Multi-world responses (rare) |
| `stubEnumerateNoResults(s, ackHandle)` | Placeholder for unimplemented enumerate packets |

The difference between Simple and Buf is the `IsBufferResponse` flag in the ack — the client uses it to determine how to parse the response data.

### Real Example

`msg_mhf_create_guild.go` → `handlers_guild.go:618` is a clean end-to-end example:

1. Packet struct has `AckHandle` and `Name` fields
2. Parse reads them from the binary frame (Shift-JIS → UTF-8 conversion)
3. Handler calls `CreateGuild(s, pkt.Name)` for DB work
4. On error: `doAckSimpleFail` with an error code
5. On success: `doAckSimpleSucceed` with the new guild ID

## What Handlers Have Access To

### Session (`s *Session`)

The session represents one connected player. Key fields available in handlers:

| Field | Type | Purpose |
|-------|------|---------|
| `s.charID` | `uint32` | Character ID |
| `s.stage` | `*Stage` | Current multiplayer stage/room |
| `s.server` | `*Server` | Channel server instance |
| `s.logger` | `*zap.Logger` | Session-scoped logger |
| `s.courses` | `[]mhfcourse.Course` | Active subscription courses |
| `s.semaphore` | `*Semaphore` | Current semaphore lock ownership |
| `s.Name` | `string` | Character name |
| `s.token` | `string` | Login token |
| `s.mailList` | `[]int` | Mail ID mapping |

Key methods: `s.QueueSendMHF(packet)`, `s.isOp()`.

### Server (`s.server`)

The server holds shared state across all sessions on the channel:

| Field | Type | Purpose |
|-------|------|---------|
| `s.server.db` | `*sqlx.DB` | PostgreSQL connection |
| `s.server.erupeConfig` | `*Config` | Server configuration |
| `s.server.stages` | `map[string]*Stage` | All active stages |
| `s.server.Channels` | `[]*Server` | All channel server instances |
| `s.server.raviente` | `*Raviente` | Raviente raid state |

Key methods: `s.server.BroadcastMHF(packet, ignored)`, `s.server.FindSessionByCharID(id)`.

### ByteFrame (`byteframe.NewByteFrame()`)

Binary data serialization. Big-endian by default.

```go
// Writing a response
bf := byteframe.NewByteFrame()
bf.WriteUint32(value)
bf.WriteUint16(count)
bf.WriteBytes(rawData)
data := bf.Data()

// Reading from a packet
bf.ReadUint32()
bf.ReadNullTerminatedBytes()
bf.Seek(10, io.SeekStart)
```

## Debugging

### Packet Logging

Enable in `config.json` under `DebugOptions`:

```json
{
    "LogInboundMessages": true,
    "LogOutboundMessages": true,
    "LogMessageData": true,
    "MaxHexdumpLength": 512
}
```

This prints every packet with its opcode name, direction, timing, and hex dump. High-frequency packets (PING, TIME, NOP, POSITION_OBJECT, EXTEND_THRESHOLD, END) are filtered out automatically.

Example output:
```
[PlayerName] -> [Server]
Opcode: (Dec: 24832 Hex: 0x6100 Name: MSG_MHF_CREATE_GUILD)
Data [24 bytes]:
00000000  61 00 00 00 00 01 00 12  54 65 73 74 00 00 00 00  |a.......Test....|
```

### Proxy

Set `ProxyPort` in debug options to redirect all channel server traffic through a proxy on that port, useful for external packet inspection tools.

### Quest Debugging

Set `QuestTools: true` to enable extra quest-related logging (file loading, backporting, caching).

## Testing

```bash
go test -v ./...                               # All tests
go test -race ./... -timeout=10m               # Race detection (mandatory)
go test -v ./server/channelserver/...          # One package
go test -run TestGuildRank ./server/channelserver/...  # One test
go test -coverprofile=coverage.out ./...       # Coverage
```

Race detection is critical — Erupe uses goroutines extensively with shared mutable state (stages, semaphores, Raviente). The CI pipeline runs `-race` on every push.

### Integration Tests

Tests prefixed with `IntegrationTest_` or using `testing.Short()` are skipped in short mode. A separate test database is available via `docker/docker-compose.test.yml` (PostgreSQL on port 5433, tmpfs-backed for speed):

```bash
docker compose -f docker/docker-compose.test.yml up -d
```

### Test Patterns

- **Unit tests**: Table-driven with subtests (`t.Run`), test business logic in isolation
- **Handler tests**: Create mock sessions, call handler directly, verify ack response
- **Integration tests**: Full session lifecycle with simulated connections

## CI

GitHub Actions runs on push to `main`, `develop`, `fix-*`, `feature-*` branches and on PRs:

1. `go test -v ./...` (10 minute timeout)
2. `go test -race ./...` (10 minute timeout)
3. Coverage report → Codecov
4. Build Linux and Windows binaries

Linting uses `golangci-lint`.
