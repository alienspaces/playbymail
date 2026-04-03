# Shared Game Configuration

Reference for platform-level configuration that applies to all game types (adventure and mecha).

## Game Settings

These settings define the game as it appears in the catalog and establish the basic rules for how turns are managed.

| Setting | Description |
|---|---|
| Name | Display name shown to players |
| Description | Appears on the join game turn sheet |
| Game type | The type of game — `adventure` or `mecha`; cannot be changed after the game is created |
| Turn duration (hours) | Default length of each turn; can be overridden per run |
| Status | `draft` while the game is being designed; `published` once it is available to players — this transition is one-way |

---

## Game Runs (Instances)

A run is a single playthrough of a published game. Each run can be configured independently.

### Run Settings

| Setting | Description |
|---|---|
| Required player count | Minimum number of players before the run can start |
| Turn duration (hours) | Turn length for this specific run; overrides the game-level setting; 0 means no fixed schedule |
| Process when all submitted | Process the turn immediately once all players have submitted their orders, rather than waiting for the turn deadline |
| Email delivery | Deliver turn sheets to players by email |
| Physical post delivery | Deliver turn sheets by physical post |
| Local physical delivery | Deliver turn sheets by local physical collection |
| Closed testing | Restrict joining to players who have been given the closed-testing key |
| Closed-testing key | The key players must enter to join when closed testing is enabled |
| Closed-testing key expiry | When the closed-testing key expires |

**Delivery:** at least one delivery method must be enabled.

**Closed testing:** if closed testing is enabled, email delivery must also be enabled.

**Turn duration:** can only be changed while the run has not yet started.

### Run Status

| Status | Meaning |
|---|---|
| Created | Run exists but has not started; players can join |
| Started | Game is actively running turns |
| Paused | Turn processing is temporarily suspended |
| Completed | Game has ended normally |
| Cancelled | Game was terminated early |

---

## Game Parameters

Parameters let you tune the rules of a specific game type. Each parameter has a default value that applies unless overridden on a particular run.

### Available Parameters

| Game Type | Parameter | Default | Description |
|---|---|---|---|
| Adventure | Character lives | 3 | Number of lives a character has before being permanently removed from the run |
| Mecha | Lance size | 4 | Number of mechs in a player's lance |

---

## Turn Sheet Background Images

Each turn sheet can have a custom background image. Images are uploaded per sheet type and apply to all sheets of that type generated for the game.

### Image Requirements

| Requirement | Value |
|---|---|
| Accepted formats | WebP, PNG, JPEG |
| Maximum file size | 1 MB |
| Minimum width | 400 px |
| Maximum width | 4,000 px |
| Minimum height | 200 px |
| Maximum height | 6,000 px |
| Recommended dimensions | 2,480 × 3,508 px (A4 at 300 DPI) |

Images that fall within the dimension limits but differ significantly from the A4 aspect ratio will generate a warning; the upload still succeeds.

### Turn Sheet Types

Every game type has a join game sheet — sent to a player when they first join a run — which must have a background image uploaded before the game is published. Each game type also has its own set of play sheets that can each have a unique background.

For the full list of sheet types and background image identifiers for each game type, see:

- [Adventure game — turn sheet background images](adventure-game.md#turn-sheet-background-images)
- [Mecha game — turn sheet background images](mecha-game.md#turn-sheet-background-images)

---

## Game Validation

Before creating a run, a game must pass a readiness check. The designer studio shows any issues that need to be resolved.

Issues are categorised as errors (which block run creation) or warnings (which are advisory only).

### Readiness Requirements

**Adventure — errors (must be resolved):**
- At least one location must exist
- At least one location must be marked as the starting location

**Adventure — warnings (advisory):**
- Location object state graph issues — unreachable states, dead-end states, or objects missing an initial state

**Mecha — errors (must be resolved):**
- At least one sector must exist
- At least one sector must be marked as the starting sector
- At least one chassis must exist
- A player starter lance must exist with at least one mech on it

---

## Designer Studio Navigation

The studio provides different sections depending on the game type.

**All game types:**
- Games list
- Turn sheet backgrounds

**Adventure only:**
Locations → Location Links → Link Requirements → Items → Item Placements → Item Effects → Creatures → Creature Placements → Location Objects → Object Effects

**Mecha only:**
Chassis → Weapons → Sectors → Sector Links → Lances
