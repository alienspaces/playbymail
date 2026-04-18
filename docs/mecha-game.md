# Mecha Game Type

Reference for the `mecha` game type — designer configuration and player rules.

For platform-level configuration shared across all game types, see [shared-game-configuration.md](shared-game-configuration.md).

---

## Game Parameter

| Parameter | Default | Description |
|---|---|---|
| Squad size | 4 | Number of mechs in a player's squad |

---

## Designer Configuration

### How the Game Fits Together

```
Game
  ├── Chassis (mech body blueprints)
  ├── Weapons (weapon definitions)
  ├── Sectors (battlefield map areas)
  │     └── Sector links (adjacency between sectors)
  ├── Computer opponents (AI behaviour profiles)
  └── Squads (groups of mechs)
        └── Squad mechs (specific mechs assigned to a squad)
```

---

### Chassis

Blueprint stats for a mech body type. All mechs in the game are based on a chassis definition.

| Field | Required | Description |
|---|---|---|
| Name | yes | Display name (e.g. "Raven Light", "Enforcer") |
| Description | no | Narrative description |
| Chassis class | yes | Weight class — determines general role and capability |
| Armor points | yes | Maximum armor; absorbed before structure takes damage; 1–1000 |
| Structure points | yes | Maximum structure; the mech is destroyed when this reaches 0; 1–1000 |
| Heat capacity | yes | Maximum heat the mech can accumulate before shutting down; 1–200 |
| Speed | yes | Sector hops allowed per turn; 1–10 (long-range weapons reach 2 hops, so values above ~5 are only useful for rapid repositioning) |
| Small slots | yes | Number of small hardpoints on the chassis; 0–10 |
| Medium slots | yes | Number of medium hardpoints on the chassis; 0–10 |
| Large slots | yes | Number of large hardpoints on the chassis; 0–10 |

A chassis must have at least one slot in total. See **Loadout slots** below for how slots are consumed.

**Chassis class values:**

| Class | Description | Default slots (small / medium / large) |
|---|---|---|
| `light` | Fast, lightly armoured | 2 / 1 / 0 |
| `medium` | Balanced performance | 2 / 2 / 1 |
| `heavy` | Slower, heavily armoured | 2 / 2 / 2 |
| `assault` | Maximum armour and firepower, minimal speed | 2 / 3 / 3 |

The class defaults are applied whenever a new chassis is created without explicit slot values (designer UI, API, and test harness all share the same defaults). Designers are free to override them.

**Loadout slots:**

Every mounted item (today: weapons; future: equipment) consumes one slot on the chassis. Items are placed with **upward spillover**:

- A **large** item can only use a large slot.
- A **medium** item prefers a medium slot, spilling into a large slot if mediums are full.
- A **small** item prefers a small slot, spilling into medium, then large.

This keeps the slot model forgiving (generous medium/large counts accommodate small items naturally) while still preventing a light chassis from carrying a large weapon. The fit check runs when a squad mech is created or updated, and again on every in-game weapon swap — swaps that would overflow the chassis are refused and reported back to the player.

**Requirement:** at least one chassis must exist before a run can be created.

---

### Weapon

Weapon definitions used in mech loadouts and refit orders.

| Field | Required | Description |
|---|---|---|
| Name | yes | Display name (e.g. "AC/10", "Medium Laser", "LRM-20") |
| Description | no | Narrative description |
| Damage | yes | Damage dealt per hit; 1–20 |
| Heat cost | yes | Heat added to the mech each time the weapon fires, even on a miss; 0–20 |
| Range band | yes | Effective engagement distance (see range band values below) |
| Mount size | yes | Mount size category (`small`, `medium`, `large`). Must fit an available slot on the chassis — see the **Loadout slots** subsection under **Chassis** for the upward-spillover rule |

**Range band values** — determines whether a weapon can fire at a given distance:

| Range band | Same sector | Adjacent sector (1 hop) | 2 sectors away |
|---|---|---|---|
| `short` (brawl) | Can fire | Cannot fire | Cannot fire |
| `medium` (versatile) | Can fire | Can fire | Cannot fire |
| `long` (standoff) | Cannot fire | Can fire | Can fire |

Long-range weapons are standoff weapons — they cannot fire into the same sector but can engage targets up to 2 hops away.

**Mount size values:**

| Size | Description |
|---|---|
| `small` | Small mount |
| `medium` | Medium mount |
| `large` | Large mount |

---

### Sector

Map areas that make up the battlefield. The map is a graph of sectors connected by sector links. Mechs occupy one sector at a time and move across connected sectors each turn.

| Field | Required | Description |
|---|---|---|
| Name | yes | Display name (e.g. "Alpha Depot", "Urban Centre", "Ridge Line") |
| Description | no | Narrative description |
| Terrain type | yes | Terrain classification (see terrain type values below); defaults to `open` |
| Elevation | no | Relative height; −10 to 10; used by the AI for tactical positioning (higher elevation is preferred by defensive opponents); default 0 |
| Cover modifier | no | Added directly to attacker hit chance for mechs in this sector; −50 to 50 (step 5 in the designer UI); negative values make mechs harder to hit; default 0 |
| Starting sector | no | If enabled, this is a depot sector — squads spawn here and management sheets are issued when mechs are present |

**Terrain type values:**

| Value | Description |
|---|---|
| `open` | Open ground |
| `urban` | Urban environment with buildings |
| `forest` | Forested area |
| `rough` | Rough terrain |
| `water` | Water terrain |

**Requirement:** at least one sector must exist and at least one must be marked as a starting sector before a run can be created.

---

### Sector Link

Directed connections between sectors that define which moves are legal. A link from sector A to sector B does not automatically create a link from B to A — both must be added separately for two-way movement.

| Field | Description |
|---|---|
| From sector | The origin sector |
| To sector | The destination sector; must differ from the origin |

Cover is now a property of the destination sector, not the link. See the **Sector** section above.

---

### Squad

A design-time squad template. There are two squad types.

| Field | Description |
|---|---|
| Name | Display name |
| Description | Description |
| Type | `starter` or `opponent` (see below) |

**Squad types:**

- **Starter** (`starter`) — the loadout cloned for every player who joins a run. Its mechs are copied into player-specific squad instances at game start. At most one starter squad is allowed per game.
- **Opponent** (`opponent`) — a template randomly assigned to a computer opponent when a run starts. If there are more opponents than templates the templates are reused. No player ever owns an opponent squad directly.

Player-owned squads only exist as runtime **squad instances** — they are never stored in the design-time squad table.

**Requirement:** a starter squad with at least one mech must exist before a run can be created.

---

### Squad Mech

A specific mech assigned to a squad — combining a chassis with a callsign and weapon loadout.

| Field | Description |
|---|---|
| Squad | The squad this mech belongs to |
| Chassis | The chassis blueprint for this mech |
| Callsign | Unique name within the squad (e.g. "Alpha-1", "Shadow Fox") |
| Weapon loadout | The weapons fitted to this mech; each weapon entry specifies the weapon and the slot it occupies (e.g. left arm, right torso, centre torso) |

---

### Computer Opponent

AI behaviour profiles for computer-controlled squads. Each profile controls how aggressively the AI plays and how tactically sophisticated its decisions are.

Computer opponents are managed through the designer studio under **Computer Opponents**.

| Field | Description |
|---|---|
| Name | Display name |
| Description | Description |
| Aggression | How aggressively the AI plays; 1 = purely defensive, 10 = all-out assault |
| IQ | Tactical sophistication; 1 = predictable or random decisions, 10 = expert use of terrain and positioning |

---

## Turn Sheets

Each turn a squad receives a set of turn sheets. Sheets are presented to the player in a specific order, and processed by the game engine in a different order.

| Sheet | Processing order | Presentation order | Notes |
|---|---|---|---|
| Join game | — | — | Sent when a player first joins; handled separately from regular turn processing |
| Squad management | 1st | 2nd | Processed first so repairs and refits are applied before movement |
| Orders | 2nd | 1st | Shown first as the primary strategic action; management is secondary |

### Turn Sheet Background Images

When uploading a background image for a mecha game, select the sheet type the image should apply to.

| Sheet | Description |
|---|---|
| `mecha_join_game` | Join game sheet — required; sent when a player first joins |
| `mecha_orders` | Movement and attack orders sheet |
| `mecha_squad_management` | Repair and refit sheet |

---

## Player Rules

### Orders Sheet

Players submit movement and attack orders for each mech in their squad.

Players submit one order per mech: an optional destination sector to move to, and an optional target mech to attack.

**Movement rules:**
- Mechs may move up to a number of sector hops equal to their chassis **Speed** — the sheet shows all sectors reachable within that many hops
- Destroyed mechs receive no movement orders
- Mechs currently refitting (undergoing repairs or weapon swaps from the previous turn's management sheet) are excluded from movement and combat
- The server validates that the chosen destination is reachable within the mech's speed; invalid moves are ignored

**Attack rules:**
- Attack declarations are collected from all squads and resolved simultaneously after all movement is applied
- Any non-destroyed enemy mech in the run is a valid attack target
- Targets must be within weapon range after movement (see range bands in the Designer Configuration section)

---

### Squad Management Sheet

Players submit repair and refit orders for mechs that are at a depot sector (a starting sector). This sheet is only issued to player-controlled squads; AI squads do not receive management sheets.

A management sheet is issued when at least one mech in the squad is in a depot sector.

Players submit per-mech orders to repair structure and/or swap weapons. Each weapon swap specifies the slot and the replacement weapon.

**Order rules:**
- Orders only apply to mechs that are currently in a depot sector
- Mechs that are already refitting from a previous turn cannot receive new management orders
- **Structure repair** restores the mech to full structure; costs supply points (at least 1); the mech enters a refitting state for the following turn
- **Weapon swaps** install a new weapon into the specified slot; each swap costs 1 supply point; the mech enters a refitting state for the following turn
- Supply points are deducted from the squad's total (cannot go below 0)
- **Refitting effect:** a mech that receives any management order is excluded from movement and combat that turn

---

### Combat Resolution

Combat is resolved simultaneously — all attack orders from all squads are collected first, then resolved together.

All attacks use the positions and hit points from before any combat damage is applied.

**Range:**
- Mechs can engage targets up to 2 sector hops away if they have long-range weapons
- Targets more than 2 hops away cannot be hit by any weapon

**Weapon firing and range:**

| Distance | Short (brawl) | Medium (versatile) | Long (standoff) |
|---|---|---|---|
| Same sector | Fire | Fire | Cannot fire |
| Adjacent (1 hop) | Cannot fire | Fire | Fire |
| 2 sectors away | Cannot fire | Cannot fire | Fire |

**Heat:**
- Each weapon fired adds heat to the mech — even on a miss
- All weapons in the mech's loadout fire together in a single attack

**Hit chance:**
- Base hit chance is 50%, plus 5% per point of pilot skill, modified by the target sector's cover modifier, capped between 0% and 95%
- Formula: `hit_chance = clamp(50 + pilot_skill × 5 + cover_modifier, 0, 95)`
- A negative cover modifier (heavy cover in the target sector) reduces hit chance; a positive modifier makes targets easier to hit

**Damage:**
- Armour absorbs damage first
- Any damage exceeding current armour carries over into structure
- The mech is destroyed when structure reaches 0

**Overheat:**
- After combat damage is applied, accumulated heat is added to the mech
- If heat exceeds the chassis heat capacity the mech shuts down (unless already destroyed)

**Mech status:**

| Status | Description |
|---|---|
| Operational | Fully functional |
| Damaged | Structure reduced but above 0 |
| Destroyed | Structure reached 0; permanently out of action |
| Shutdown | Overheated; cannot act |

---

### End-of-Turn Lifecycle

After combat is resolved, the engine applies the following in order:

1. **Heat dissipation** — heat accumulated during combat is reduced for all mechs
2. **Auto armor repair** — operational mechs in depot sectors receive partial armor restoration
3. **Supply point accrual** — squads receive supply points each turn (used for management orders)
4. **Pilot XP and skill advancement** (see below)

---

### Pilot Progression

Pilots earn experience points (XP) through combat each turn:

| Event | XP earned |
|---|---|
| Participating in an attack | 1 XP |
| Destroying an enemy mech | 2 XP |

XP accumulates across turns. When a pilot's total XP crosses a threshold, their **pilot skill** increases by 1. Higher pilot skill directly improves hit chance (see **Hit chance** above).

**Pilot skill thresholds:**

| Skill level | Total XP required |
|---|---|
| 0 | 0 (starting skill) |
| 1 | 3 |
| 2 | 8 |
| 3 | 15 |
| 4 | 24 |
| 5 | 35 |
| 6 | 48 |
| 7 | 63 |
| 8 | 80 |
| 9 | 99 |

All pilots start at skill level 0.

---

### Computer Opponent AI

The AI makes decisions for all computer-controlled squads each turn, after player orders have been submitted.

**Movement behaviour:**
- Destroyed or shutdown mechs receive no orders
- The AI uses BFS pathfinding and moves up to the mech's chassis **Speed** each turn
- High-aggression opponents (7 or above) advance toward the nearest enemy; tactically skilled (high IQ) opponents prefer routes through high-elevation or high-cover sectors
- Low-aggression opponents (3 or below) fall back toward high-elevation, high-cover positions
- Mid-aggression opponents hold position or move to the best available defensive sector

**Targeting behaviour:**
- High-aggression opponents prefer to finish off weakened mechs (lowest structure)
- Low-aggression opponents prefer to deter the strongest threats (highest structure)
- The AI can engage targets up to 2 hops away (matching long-range weapon capability)
- If no enemies are in range, no attack is issued

When configured, the AI can use an advanced reasoning mode to make more sophisticated decisions based on the full battlefield situation.
