# Distribution Types
Distributions are available from the Guide near the entrance of Mezeporta.
These control which menu at the Guide your Distribution will show up. 

| ID | Name              |
|----|-------------------|
| 0  | Bought            |
| 1  | Event             |
| 2  | Compensation      |
| 4  | Promotion         |
| 6  | Subscription      |
| 7  | Event Item        |
| 8  | Promotion Item    |
| 9  | Subscription Item |

# Item Types 
These define the item type you would receive from a Distribution or Event Tent Distribution.

| ID | Type                 | Notes |
|----|----------------------|-------|
| 0  | Legs                 | Equipment slot 0 |
| 1  | Head                 | Equipment slot 2 |
| 2  | Chest                | Equipment slot 3 |
| 3  | Arms                 | Equipment slot 4 |
| 4  | Waist                | Equipment slot 5 |
| 5  | Melee                | Equipment slot 6 |
| 6  | Ranged               | Equipment slot 7 |
| 7  | Item                 | Uses item_id + quantity |
| 8  | Furniture            | Gallery box item |
| 9  | Campaign Quest       | Quest unlock, no client reward |
| 10 | Zenny                | |
| 11 | Unused               | Client no-op (G1 through ZZ) |
| 12 | Festival Points      | Auto-handled by festa system |
| 13 | Unused               | Client no-op (G1 through ZZ) |
| 14 | Tore Points          | |
| 15 | Poogie Outfits       | Uses item_id for outfit |
| 16 | Restyle Points       | |
| 17 | N Points             | Server: addPointNetcafe |
| 18 | Goocoo Outfits       | |
| 19 | Gacha Koban          | Server: gacha_premium |
| 20 | Trial Gacha Koban    | Server: gacha_trial |
| 21 | Frontier Points (FP) | Server: frontier_points |
| 22 | UNKNOWN              | ZZ-only, uses item_id + amount, notify-only |
| 23 | Ryoudan Points (RP)  | Server: saveData.RP |
| 24 | UNKNOWN              | ZZ-only, amount only, notify-only |
| 25 | Bond/Kizuna Points   | |
| 26 | Unused               | Client no-op (G1 through ZZ) |
| 27 | Unused               | Client no-op (G1 through ZZ) |
| 28 | Special Hall         | |
| 29 | Song Note            | |
| 30 | Item Box pages       | |
| 31 | Equipment Box pages  | Obfuscated stamp system |

# Course Types

| ID      | Course                                                |
|---------|-------------------------------------------------------|
| 1       | TRIAL (Trial, TL)                                     |
| 2       | BASIC (HunterLife, HL)                                |
| 3       | EXTRA_A (Extra, ExtraA, EX)                           |
| 4       | EXTRA_B (ExtraB)                                      |
| 5       | MOBILE (Mobile)                                       |
| 6       | PREMIUM (Premium)                                     |
| 7       | PALLONE (Pallone, ExtraC)                             |
| 8       | LEGEND (Assist, Legend, Rasta)                        |
| 9       | XBCAFE (N)                                            |
| 10      | SECRET (Hiden, Secret)                                |
| 11      | ROYAL (HunterSupport, HunterAid, Support, Aid, Royal) |
| 12      | NBOOST (NBoost, NetCafeBoost, Boost)                  |
| 13 - 19 | show up as (unknown)                                  |
| 20      | DEBUG                                                 |
| 21      | COG_LINK_EXPIRED                                      |
| 22      | 360_GOLD                                              |
| 23      | PS3_TROP                                              |
| 24      | COG                                                   |
| 25      | CAFE_SP                                               |
| 26      | OFFICIAL (NetCafe, Cafe, OfficialCafe, Official)      |
| 27      | CARD (HLRenewing, HLR, HLRenewal, HLRenew, CardHL)    |
| 28      | CARD_EX (EXRenewing, EXR, EXRenewal, EXRenew, CardEX) |
| 29      | FREE                                                  |
| 30      | NETCAFE                                               |

# Event Quests

## Quest Types

`event_quests.quest_type` is a server-supplied event category. It is **not** the client menu index or the internal quest-list selector. A menu can combine server event types with quests listed in the client data.

The client labels below were traced in the supplied ZZ client (`mhfo.dll` and `mhfpac.bin`). **Client code** means a static code/data trace, not a successful in-game test. Promotional (9) was also tested in game. Other client builds/localizations may differ.

| Type | Upstream name / alias | Client entry | Evidence |
|---:|---|---|---|
| 5 | Frontier Quests | Frontier Quests | Client code |
| 7 | Hunting Tournament | Hunting Competition | Client code |
| 8 | Festa Winners | Winners Quests | Client code |
| 9 | Promotional (Campaign) | Special Quests > Promotional | Client code; in-game confirmed |
| 10 | Elder Dragon Quests | Elder Dragon Quests | Client code |
| 11 | Special Quests | Special Quests | Client code |
| 12 | Frontier Quests (additional client filter) | Frontier Quests | Client code |
| 13 | VS | Quest Orders | Client code |
| 14 | Goushu | Gou Quests | Client code |
| 15 | Goushu (Premium) | Gou Quests | Client code |
| 16 | Ravi (Normal) | Great Slaying Quests | Client code |
| 18 | HR Event | Event Quests (HR) | Client code |
| 19 | Certified NetCafe Benefits | Certified NetCafe Benefits | Client code |
| 21 | Hunting Moves | Hunting Skill Quests | Client code |
| 22 | Ravi (Violent) | Not independently mapped in this investigation | Retained from upstream; not reclassified |
| 24 | Series | Not independently mapped in this investigation | Retained from upstream; not reclassified |
| 25 | Hunting Training | Training Quests | Client code |
| 26 | Daily | Daily Quests (HR) | Client code |
| 27 | G Frontier Quest (Clan Hall Only) | Frontier Quests | Client code |
| 28 | G Event | Event Quests (G) | Client code |
| 29 | G Special Quests | Special Quests | Client code |
| 30 | G Daily | Daily Quests (G) | Client code |
| 31 | G Urgent | G Urgent Quests | Client code |
| 32 | G Frontier Quest | Frontier Quests | Client code |
| 33 | Extreme Conquest | Not independently mapped in this investigation | Retained from upstream; not reclassified |
| 34 | Supremacy Quest | Supremacy Quests | Client code |
| 35 | Burst Origin Quest | ※ Burst & Origin Quests | Client code |
| 36 | Heavenly Conquest | Not independently mapped in this investigation | Retained from upstream; not reclassified |
| 38 | HR Exotic | ※ Exotic Quests | Client code |
| 39 | G Exotic | ※ G Exotic Quests | Client code |
| 40 | Ravi (Berserk) | Not independently mapped in this investigation | Retained from upstream; not reclassified |
| 41 | Returning Clan Quests | Returning Clan Quests | Client code |
| 42 | G Returning Clan Quests | G Returning Clan Quests | Client code |
| 43 | Zenith | Zenith Quests > Zenith Monster Quests (dedicated filter) | Client code |
| 44 | Special Clan Quests | Special Clan Quests | Client code |
| 45 | G Special Clan Quests | G Special Clan Quests | Client code |
| 46 | Interception | 迎撃戦クエスト; Ｇ級迎撃戦クエスト | Client code |
| 47 | Interception Branch | 分岐ルート用クエスト | Client code |
| 48 | Interception Urgent | Ｇ級迎撃戦クエスト | Client code |
| 50 | Ravi (Extreme) | Not independently mapped in this investigation | Retained from upstream; not reclassified |
| 51 | Ravi (Small Berserk) | Not independently mapped in this investigation | Retained from upstream; not reclassified |
| 52 | Superior | Special Quests > Superior Quests (HR filter) | Client code |
| 53 | G Superior | Special Quests > Superior Quests (G filter) | Client code |
| 54 | Gear Acquisition Quests (formerly G Armor) | Other Quests > Gear Acquisition Quests; also the G-rank Gear Acquisition entry | Client code |

### Menus backed by internal client lists

These numbers are **internal selectors**, not values to put into `event_quests.quest_type`. The source of the list can be `mhfinf.bin`, `mhfdat.bin`, or a dedicated client routine.

| Client entry | Server event type accepted by this entry | Internal selector(s), not quest_type |
|---|---|---|
| Gear Acquisition Quests | 54 | 56 |
| Gear Acquisition Quests (G-rank entry) | 54 | 119 |
| G Deco Quests | No direct event-type filter in this descriptor | 57 |
| Paw Coins/Yearly Benefits | No direct event-type filter in this descriptor | 38 |
| HR Ranking Rewards | No direct event-type filter in this descriptor | 39 |
| Weapon Souls & Hiden | No direct event-type filter in this descriptor | 66 |
| Premium Quests | No direct event-type filter in this descriptor | 0, 1 |
| Luxury Quests | No direct event-type filter in this descriptor | 8 |
| Training Quests | No direct event-type filter in this descriptor | 53 |
| Special Quests | No direct event-type filter in this descriptor | 54 |
| Memo Quests | No direct event-type filter in this descriptor | 55 |
| Skill Cuff Quests | No direct event-type filter in this descriptor | 58 |
| Long Term Quests | No direct event-type filter in this descriptor | 59 |
| Hunting Skill Quests | No direct event-type filter in this descriptor | 60 |

Type **54**, not internal selector **57**, is the server event filter accepted by **Gear Acquisition Quests**. The separate **G Deco Quests** entry uses internal selector 57. HR/GR, course, time and other eligibility checks still apply.

Erupe v9.6.1 applies campaign-unlock checks to type 9 in `makeEventQuest`; a visible Promotional entry does not guarantee every type-9 quest is unlocked on every server. See [client mapping evidence](Quest-Type-Client-Mapping) for the exact scope and addresses.

## Quest Flags
|   Season/Time    | Flag  |
|------------------|-------|
| Automatic        | 0     |
| Breeding/Day     | 9     |
| Warm/Day         | 10    |
| Cold/Day         | 12    |
| Breeding/Night   | 17    |
| Warm/Night       | 18    |
| Cold/Night       | 32    |

## Mark
|  Type       | Bitfield |
|-------------|----------|
| None        |    0     |
| Recommended |    1     |
| New         |    2     |
| Refine      |    4     |
