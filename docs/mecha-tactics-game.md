# Mecha Tactics Game Type

Reference for the `mecha_tactics` game type — designer configuration and player rules.

Mecha Tactics is an individual mech combat game played on a hex grid. Each player controls one or more mechs directly. There are no lances — mechs are owned and operated individually. Tactical positioning, facing, movement efficiency, and weapon range management are the core of play.

For platform-level configuration shared across all game types, see [shared-game-configuration.md](shared-game-configuration.md).

---

## Game Parameter

| Parameter | Default | Description |
|---|---|---|
| Mech count | 1 | Number of mechs each player controls |

---

## Designer Configuration

### How the Game Fits Together

```
Game
  ├── Chassis (mech body blueprints)
  ├── Weapons (weapon definitions)
  ├── Terrain types (movement and cover properties for hex terrain)
  ├── Hex map (battlefield grid)
  │     └── Hexes (individual map cells with terrain, elevation, and depot flags)
  └── Computer opponents (AI behaviour profiles)
```

---

### Chassis

Blueprint stats for a mech body type. All mechs in the game are based on a chassis definition.

| Field | Description |
|---|---|
| Name | Display name (e.g. "Raven Light", "Enforcer") |
| Description | Narrative description |
| Chassis class | Weight class — determines general role and capability |
| Armor points | Maximum armor; absorbed before structure takes damage; must be greater than 0 |
| Structure points | Maximum structure; the mech is destroyed when this reaches 0; must be greater than 0 |
| Heat capacity | Maximum heat the mech can accumulate before shutting down; must be greater than 0 |
| Speed | Movement point (MP) budget per turn; the number of MP available to spend on hex entry; must be greater than 0 |

**Chassis class values:**

| Class | Description |
|---|---|
| `light` | Fast, lightly armoured — high MP budget, vulnerable to sustained fire |
| `medium` | Balanced performance — moderate MP and armour |
| `heavy` | Slower, heavily armoured — lower MP, absorbs significant damage |
| `assault` | Maximum armour and firepower, minimal speed — lowest MP budget |

**Requirement:** at least one chassis must exist before a run can be created.

---

### Weapon

Weapon definitions used in mech loadouts and refit orders.

| Field | Description |
|---|---|
| Name | Display name (e.g. "AC/10", "Medium Laser", "LRM-20") |
| Description | Narrative description |
| Damage | Damage dealt per hit; must be greater than 0 |
| Heat cost | Heat added to the mech each time the weapon fires; can be 0 |
| Range band | Effective engagement distance measured in hex distance (see range band values below) |
| Mount size | Physical size of the weapon mount required |

**Range band values** — determines whether a weapon can fire at a given hex distance:

| Range band | Distance 0 (same hex) | Distance 1 (adjacent) | Distance 2 |
|---|---|---|---|
| `short` | Can fire | Cannot fire | Cannot fire |
| `medium` | Can fire | Can fire | Cannot fire |
| `long` | Cannot fire | Can fire | Can fire |

Short-range weapons are brawling weapons — effective only when sharing a hex with the target. Medium-range weapons are versatile — effective at distance 0 and 1. Long-range weapons are standoff weapons — they require a minimum distance of 1 hex and can reach 2 hexes away, but cannot fire at targets in the same hex.

**Mount size values:**

| Size | Description |
|---|---|
| `small` | Small mount |
| `medium` | Medium mount |
| `large` | Large mount |

---

### Terrain Type

Designer-defined terrain classifications applied to hexes on the battlefield. Each terrain type has its own movement point cost and optional cover modifier, allowing the game designer to tailor the tactical feel of the map.

| Field | Description |
|---|---|
| Name | Display name (e.g. "Open Ground", "Urban Ruins", "Dense Forest") |
| Description | Narrative description shown on turn sheets |
| Movement point cost | Number of movement points (MP) a mech spends to enter a hex of this terrain; must be greater than 0 |
| Cover modifier | Adjustment to hit chance when attacking a mech in this terrain; negative values make the target harder to hit; can be 0 |

**Requirement:** at least one terrain type must exist before a run can be created.

---

### Hex

Individual cells that make up the battlefield hex grid. The map is defined by placing hexes at axial coordinates (column, row). Adjacency between hexes is determined automatically by their coordinates — there is no need to manually link hexes.

| Field | Description |
|---|---|
| Column | Axial column coordinate |
| Row | Axial row coordinate |
| Name | Optional display name (e.g. "Alpha Depot", "Ridge Line") |
| Description | Optional narrative description |
| Terrain type | The terrain type applied to this hex; determines MP entry cost and cover |
| Elevation | Relative height; used by the AI for tactical positioning and shown on turn sheets |
| Starting hex | If enabled, this hex is a depot — mechs spawn here and repair sheets are issued when a mech is present |

**Adjacency:** every hex has up to 6 neighbours determined by its axial coordinates. Movement and weapon range are measured in hex distance — the minimum number of hex steps between two hexes.

**Requirement:** at least one hex must exist and at least one must be marked as a starting hex before a run can be created.

---

### Computer Opponent

AI behaviour profiles for computer-controlled mechs. Each profile controls how aggressively the AI plays and how tactically sophisticated its decisions are.

Computer opponents are not yet configurable through the designer studio — they are configured when setting up a game.

| Field | Description |
|---|---|
| Name | Display name |
| Description | Description |
| Aggression | How aggressively the AI plays; 1 = purely defensive, 10 = all-out assault |
| IQ | Tactical sophistication; 1 = predictable or random decisions, 10 = expert use of terrain, facing, and positioning |

---

## Turn Sheets

Each turn a player receives one turn sheet per mech they control. Sheets are presented to the player in a specific order, and processed by the game engine in a different order.

| Sheet | Processing order | Presentation order | Notes |
|---|---|---|---|
| Join game | — | — | Sent when a player first joins; handled separately from regular turn processing |
| Mech repair | 1st | 2nd | Processed first so repairs are applied before movement and combat |
| Orders | 2nd | 1st | Shown first as the primary tactical action; repair is secondary |

### Turn Sheet Background Images

When uploading a background image for a mecha tactics game, select the sheet type the image should apply to.

| Sheet | Description |
|---|---|
| `mecha_tactics_join_game` | Join game sheet — required; sent when a player first joins |
| `mecha_tactics_orders` | Movement and attack orders sheet |
| `mecha_tactics_repair` | Repair and refit sheet |

---

## Player Rules

### Facing

Each mech has a facing — one of the 6 faces (sides) of its hex. Facing is chosen freely at the end of every turn as part of the mech's orders, with no movement point cost. A stationary mech may also freely reorient. Facing determines the mech's firing arcs, which affect hit chance in combat.

A hex has 6 faces. When a mech faces through one of those faces, the 6 faces divide into arcs:

```
            ↑ Facing
        ___________
       /     F     \    ← Front arc (3 faces)
      / F         F \
      │     [M]     │  
      \ S         S /   ← Left side (1 face) + Right side (1 face)
       \_____R_____/    ← Rear arc (1 face)
```

- **Front arc** — the 3 faces directly ahead (Front-Left, Front, Front-Right) and their adjacent hexes
- **Left side arc** — the 1 face to the left and its adjacent hex
- **Right side arc** — the 1 face to the right and its adjacent hex
- **Rear arc** — the 1 face directly behind and its adjacent hex

All 6 faces are covered, accounting for all 6 adjacent hexes.

---

### Orders Sheet

Players submit a destination hex, a final facing, and an attack declaration for each mech they control.

**Movement rules:**
- Each mech has a movement point (MP) budget equal to its chassis speed
- Entering a hex costs the MP cost of that hex's terrain type; the budget limits how far the mech can reach this turn
- Facing is chosen freely at the destination — it costs no movement points and is always available as an option
- The orders sheet presents only the hexes the mech can legally reach given its current MP budget; players pick from that set rather than calculating costs themselves
- On the interactive (HTML) sheet, reachable hexes are shown as selectable; on the print (PDF) sheet, only reachable hexes are labelled — unreachable hexes are blank
- Destroyed mechs receive no movement orders
- Mechs currently undergoing repairs (from the previous turn's repair sheet) are excluded from movement and combat

**Attack rules:**
- Each mech may declare one attack per turn — a single target mech
- The attack is resolved after all movement is applied
- The target must be within range of at least one weapon in the mech's loadout after movement is complete
- Attack declarations from all mechs are collected and resolved simultaneously

---

### Mech Repair Sheet

Players submit repair orders for mechs that are at a depot hex (a starting hex). This sheet is only issued to player-controlled mechs; AI mechs do not receive repair sheets.

A repair sheet is issued when the mech is in a depot hex.

Players submit orders to repair structure and/or swap weapons. Each weapon swap specifies the slot and the replacement weapon.

**Order rules:**
- Orders only apply to mechs that are currently in a depot hex
- Mechs that are already undergoing repairs from a previous turn cannot receive new repair orders
- **Structure repair** restores the mech to full structure; costs supply points (at least 1); the mech enters a repair state for the following turn
- **Weapon swaps** install a new weapon into the specified slot; each swap costs 1 supply point; the mech enters a repair state for the following turn
- Supply points are deducted from the player's total (cannot go below 0)
- **Repair effect:** a mech that receives any repair order is excluded from movement and combat that turn

---

### Combat Resolution

Combat is resolved simultaneously — all attack orders from all mechs are collected first, then resolved together.

All attacks use the positions, facings, and hit points from before any combat damage is applied.

**Range:**
- Range between mechs is measured in hex distance (minimum hex steps between the two hexes)
- Each weapon fires only if the target is within that weapon's range band
- Only weapons that can legally fire at the measured distance contribute to the attack

**Weapon firing and range:**

| Hex distance | Short-range weapons | Medium-range weapons | Long-range weapons |
|---|---|---|---|
| 0 (same hex) | Fire | Fire | Cannot fire |
| 1 (adjacent) | Cannot fire | Fire | Fire |
| 2 | Cannot fire | Cannot fire | Fire |
| 3+ | Cannot fire | Cannot fire | Cannot fire |

**Firing arcs:**

The arc is determined by the **defender's** facing — a mech attacked from behind gains a rear-arc penalty regardless of which direction the attacker is facing.

| Arc | Faces and adjacent hexes covered | Hit chance modifier |
|---|---|---|
| Front arc | The 3 faces ahead of the defender (Front-Left, Front, Front-Right) | No modifier (base hit chance) |
| Left side arc | The 1 face to the defender's left | +10% hit chance |
| Right side arc | The 1 face to the defender's right | +10% hit chance |
| Rear arc | The 1 face directly behind the defender | +20% hit chance |

**Heat:**
- Each weapon that fires adds its heat cost to the mech — even on a miss
- All weapons in the mech's loadout that can legally fire at the target's range fire together in a single attack

**Hit chance:**
- Base hit chance is 50%, plus 5% per point of pilot skill, capped at 95%
- Firing arc modifier is applied after the pilot skill bonus, before the cap
- The terrain cover modifier of the **target's** hex is subtracted from the final hit chance

**Damage:**
- Armour absorbs damage first
- Any damage exceeding current armour carries over into structure
- The mech is destroyed when structure reaches 0

**Overheat:**
- After combat damage is applied, accumulated heat is added to the mech
- If heat exceeds the chassis heat capacity the mech shuts down (unless already destroyed)
- A shutdown mech cannot move or attack the following turn

**Mech status:**

| Status | Description |
|---|---|
| Operational | Fully functional |
| Damaged | Structure reduced but above 0 |
| Destroyed | Structure reached 0; permanently out of action |
| Shutdown | Overheated; cannot move or attack next turn |

---

### Computer Opponent AI

The AI makes decisions for all computer-controlled mechs each turn, after player orders have been submitted.

**Movement behaviour:**
- Destroyed or shutdown mechs receive no orders
- High-aggression opponents (7 or above) advance toward the nearest enemy using the lowest-cost hex path that fits within the mech's MP budget; tactically skilled opponents prefer routes through higher-cover terrain and higher-elevation hexes
- Low-aggression opponents (3 or below) fall back toward depot hexes and high-elevation positions with favourable terrain cover
- Mid-aggression opponents hold position, using remaining MP to adjust facing and place the most valuable enemy targets in the rear or side arc

**Targeting behaviour:**
- The AI selects a target within weapon range, prioritising targets in a favourable arc (side or rear)
- High-aggression opponents prefer to finish off weakened mechs (lowest structure)
- Low-aggression opponents prefer to deter the strongest threats (highest structure)
- If no enemies are within range of any weapon, no attack is issued

When configured, the AI can use an advanced reasoning mode to make more sophisticated decisions based on the full battlefield situation, including hex distances, terrain costs, facing, and firing arc positioning.
