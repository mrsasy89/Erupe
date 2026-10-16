# Quest Type Client Mapping

This page records static findings for the supplied ZZ client. No client binary is redistributed. Runtime confirmation is limited to type 9 (Promotional); type 54 still requires an in-game test on the target server.

## Verified category corrections

| Finding | Evidence |
|---|---|
| Other Quests contains Gear Acquisition | Menu descriptor 88 has child 98 in its byte list at preferred VA `0x11939528`. |
| Gear Acquisition accepts server type 54 | Descriptor 98 at preferred VA `0x11939c30` points to `0x11a96d70`: uint32 values `[2, 54, 1, 56, 0]`. Mode 2 means an event quest_type filter; mode 1 means an internal selector. |
| G-rank Gear Acquisition also accepts type 54 | Descriptor 245 at `0x1193a560` points to `[2, 54, 1, 119, 0]` at `0x11a9703c`. |
| G Deco is a separate internal list | Descriptor 99 at `0x11939c40` uses `[1, 57, 0]` at `0x11a96d84`; 57 is not a quest_type. |
| Promotional accepts type 9 | Descriptor 60 at `0x119399d0` uses `[2, 9, 0]` at `0x11a96a28`. Special Quests descriptor 3 includes child 60. |

Names are read from `mhfpac.bin` after ECD decryption and JPK type-4 decompression: header pointer `0x108c` → first nested pointer list → menu index. Menu indices are distinct from both event types and internal selectors.

## Code trace

- `0x1035b5d0`: dispatches the filter list; mode 1 calls `0x10362570`, mode 2 calls `0x10362110`.
- `0x103621ba`–`0x103621c5`: compares the requested event type against the received event record's category byte.
- `0x103622b1`–`0x103622ce`: type-54 HR/GR filtering uses counter state and quest-body byte `+0x97`, bit `0x08`.
- `0x10362570`: dispatches internal quest-list selectors. Do not copy these selector numbers into SQL as event types.
- Zenith's dedicated selector calls a routine at `0x10364d30`, which checks event type `0x2b` (43).

Addresses are preferred virtual addresses for this exact DLL, not universal offsets for every build.

## Input identity

| File | SHA-256 |
|---|---|
| `mhfo.dll` | `5596017aa16f01c06b92b64cd42d2af9a0858e9f5fedce2772bdb35eb985214f` |
| `mhfdat.bin` | `cb4c66eaaf91031ce5d67b14e1a28bfa73308493b6d19d2d024e4d206cffc476` |
| `mhfinf.bin` | `28a22e7a861376be0351239de983de79808a748c782a537c09ae0f449b6a030b` |
| `mhfpac.bin` | `2d52f80584ec7d641914af46d62492abdf6b59cb3149d4f0f86dce0c13c0e07c` |

The ECD CRC checks passed for all three data files. JPK output lengths matched their declared sizes. The supplied `mhfinf.bin` already contains definitions for 64 G Deco IDs: 21574–21599 and 21601–21638. Their presence does not by itself prove runtime visibility or make server-added records redundant.

## Scope

Menu descriptors can be historical, hidden or restricted to a particular counter. A static event filter proves the association, not that every player can access that menu. Values retained only from upstream are explicitly marked in [Enumerations](Enumerations#quest-types). Types without a traced association have not been renamed based on guesses.

The analysis used GNU objdump, REA 6.5.0 for recorded binary inspection, and ECD/JPK decoding based on the public Erupe implementation. References: [Erupe quest handler](https://github.com/Houmgaor/Erupe/blob/v9.6.1/server/channelserver/handlers_quest.go), [ECD decoder](https://github.com/Houmgaor/Erupe/blob/v9.6.1/common/decryption/ecd.go), [JPK decoder](https://github.com/Houmgaor/Erupe/blob/v9.6.1/common/decryption/jpk.go), and [REA](https://github.com/morluto/rea).
