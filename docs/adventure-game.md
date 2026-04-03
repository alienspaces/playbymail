# Adventure Game Type

Reference for the `adventure` game type — designer configuration and player rules.

For platform-level configuration shared across all game types, see [shared-game-configuration.md](shared-game-configuration.md).

---

## Game Parameter

| Parameter | Default | Description |
|---|---|---|
| Character lives | 3 | Number of lives a character has before being permanently removed from the run |

---

## Designer Configuration

### How the World Fits Together

```
Game
  ├── Locations (the places characters can be)
  │     ├── Item placements (items initially found here)
  │     ├── Creature placements (creatures initially here)
  │     └── Location objects (interactive scenery)
  │             ├── Object states (e.g. locked, open, broken)
  │             └── Object effects (what happens when players act on the object)
  ├── Location links (paths between locations)
  │     └── Link requirements (conditions to see or use a path)
  ├── Items (things characters can carry)
  │     └── Item effects (what happens when players act on an item)
  └── Creatures (monsters and NPCs)
```

---

### Location

Named places in the adventure world. The world is built by connecting locations with links.

| Field | Description |
|---|---|
| Name | Display name shown on turn sheets |
| Description | Narrative description shown to players when they are at this location |
| Starting location | If enabled, new characters spawn here and dead characters respawn here |

**Requirement:** at least one location must exist and at least one must be marked as the starting location before a run can be created.

A unique background image can be uploaded for each location to display on that location's turn sheets.

---

### Location Link

Directed connections between locations (exits / paths). A link from A to B does not automatically create a link from B to A — both directions must be added separately for two-way movement.

| Field | Description |
|---|---|
| From location | The origin location |
| To location | The destination location |
| Name | Name of the exit shown to players (e.g. "North Door", "Rope Bridge") |
| Description | Narrative description of the link when it is visible and unlocked |
| Locked description | Description shown when the link is locked or blocked |
| Traversal description | Narrative shown when a player successfully moves through this link |

---

### Location Link Requirement

Conditions that control whether a link is visible to players or can be traversed. Multiple requirements on the same link and purpose all apply — every condition must be met.

| Field | Description |
|---|---|
| Purpose | `traverse` — gates movement through the link; `visible` — gates whether the link appears on the sheet at all |
| Item | Item required (used with item-based conditions) |
| Creature | Creature required (used with creature-based conditions) |
| Condition | The condition that must be met (see below) |
| Quantity | Number of instances required to satisfy the condition |

**Item-based conditions:**

| Condition | Meaning |
|---|---|
| `in_inventory` | Player has the required number of this item in their inventory (unused) |
| `equipped` | Player has the required number of this item equipped |

**Creature-based conditions:**

| Condition | Meaning |
|---|---|
| `dead_at_location` | At least the required number of this creature are dead at the current location |
| `none_alive_at_location` | No living instances of this creature exist at the current location |
| `none_alive_in_game` | No living instances of this creature exist anywhere in the run |

**Effect on the turn sheet:** if a visibility requirement is not met, the link is hidden entirely. If a traversal requirement is not met, the link is shown but marked as locked, displaying the locked description.

---

### Location Object

Interactive scenery at a location. Objects can have discrete states and trigger effects when players perform actions on them.

| Field | Description |
|---|---|
| Location | The location this object belongs to |
| Name | Name shown to players (e.g. "Locked Chest", "Ancient Lever") |
| Description | Narrative description |
| Initial state | The state the object starts in for each new run; leave unset for a stateless object |
| Hidden | If enabled, the object is not visible until revealed by an effect |

---

### Location Object State

Discrete named states for a location object (e.g. "intact", "open", "broken", "activated"). States form a state machine — effects transition objects from one state to another.

| Field | Description |
|---|---|
| Name | State name (e.g. "locked", "unlocked", "open") |
| Description | Description shown when the object is in this state |
| Sort order | Display ordering for the states |

**Validation:** the game validator checks for unreachable states, dead-end states, and objects whose initial state does not match any defined state.

---

### Location Object Effect

Rules that define what happens when a player performs an action on an object. Multiple effects for the same action on the same object state all fire together.

| Field | Description |
|---|---|
| Action type | The player action that triggers this effect (see action types below) |
| Required object state | The object must be in this state for the effect to fire; leave unset to fire in any state |
| Required item | The player must have this item in their inventory for the effect to fire |
| Result description | Narrative text shown to the player when the effect fires |
| Effect type | The mechanical outcome (see effect types below) |
| Target state | For state-change effects — the state to transition to |
| Target item | For item effects — the item to give, remove, or place |
| Target link | For link effects — the link to open or close |
| Target creature | For summon effects — the creature to spawn |
| Target object | For object effects — the object to reveal, hide, change, or remove |
| Target location | For teleport and item-placement effects — the destination |
| Value (min) | Minimum value for damage and heal effects |
| Value (max) | Maximum value for damage and heal effects |
| Repeatable | Whether the effect can fire more than once per run |

**Action types:**

| Action | Description |
|---|---|
| `inspect` | Look at the object |
| `touch` | Touch the object |
| `open` | Open the object |
| `close` | Close the object |
| `lock` | Lock the object |
| `unlock` | Unlock the object |
| `search` | Search the object |
| `break` | Break the object |
| `push` | Push the object |
| `pull` | Pull the object |
| `move` | Move the object |
| `burn` | Burn the object |
| `read` | Read the object |
| `take` | Take from the object |
| `listen` | Listen to the object |
| `insert` | Insert something into the object |
| `pour` | Pour something on or into the object |
| `disarm` | Disarm the object (e.g. a trap) |
| `climb` | Climb the object |
| `use` | Use the object generically |

**Effect types:**

| Effect | Description |
|---|---|
| `info` | Narrative only — no mechanical change |
| `nothing` | No effect |
| `change_state` | Change this object to the target state |
| `change_object_state` | Change another object to the target state |
| `give_item` | Add the target item to the player's inventory |
| `remove_item` | Remove the target item from the player's inventory |
| `place_item` | Place the target item at the target location |
| `open_link` | Remove traversal restrictions from the target link |
| `close_link` | Add a traversal restriction to the target link |
| `reveal_object` | Make the target object visible |
| `hide_object` | Hide the target object |
| `damage` | Deal damage to the player (random between value min and max) |
| `heal` | Restore the player's health (random between value min and max; cannot exceed 100) |
| `summon_creature` | Spawn the target creature at the player's current location |
| `teleport` | Move the player to the target location |
| `remove_object` | Permanently remove this object from the run |

---

### Item

Item definitions. Items can be carried in inventory, equipped, and used. Items marked as starting items are given to every new character when a run starts.

| Field | Description |
|---|---|
| Name | Display name |
| Description | Narrative description |
| Can be equipped | Whether the item can be equipped to an equipment slot |
| Item category | Free-text category for grouping (e.g. "weapon", "potion", "key") |
| Equipment slot | The slot the item occupies when equipped; required if the item can be equipped |
| Starting item | If enabled, every new character starts the run with one of these |

**Equipment slot values:**

| Slot | Description |
|---|---|
| `weapon` | Weapon slot — determines combat damage |
| `armor` | Armour slot — determines damage absorbed |
| `clothing` | Clothing slot |
| `jewelry` | Jewelry slot |

---

### Item Effect

Rules for what happens when a player performs an action on an item. Passive effects (`weapon_damage`, `armor_defense`) apply while the item is equipped; all other effects trigger when the specified action is performed.

| Field | Description |
|---|---|
| Action type | The action that triggers this effect |
| Required item | Player must also have this item in inventory for the effect to fire |
| Required location | Player must be at this specific location for the effect to fire |
| Result description | Narrative text shown when the effect fires |
| Effect type | The mechanical outcome (see effect types below) |
| Target item | For item effects — the item to give or remove |
| Target link | For link effects — the link to open or close |
| Target creature | For summon effects — the creature to spawn |
| Target location | For teleport effects — the destination |
| Value (min) | Minimum value for damage, heal, and passive stat effects |
| Value (max) | Maximum value for damage and heal effects |
| Repeatable | Whether the effect can fire more than once |

**Action types:**

| Action | Description |
|---|---|
| `use` | Player activates the item |
| `equip` | Player equips the item |
| `unequip` | Player removes the item from a slot |
| `inspect` | Player inspects the item |
| `drop` | Player drops the item |
| `pickup` | Player picks up the item |

**Effect types:**

| Effect | Trigger | Description |
|---|---|---|
| `info` | Active | Narrative only |
| `nothing` | Active | No effect |
| `weapon_damage` | Passive (while equipped) | Sets weapon damage; value (min) is the damage dealt per attack |
| `armor_defense` | Passive (while equipped) | Sets armor defence; value (min) is damage absorbed; stacks across all equipped armour |
| `damage_target` | Active | Deal damage to a target creature (value min–max) |
| `damage_wielder` | Active | Deal damage to the player using the item |
| `heal_target` | Active | Restore health to a target (cap 100) |
| `heal_wielder` | Active | Restore the player's own health |
| `teleport` | Active | Move the player to the target location |
| `open_link` | Active | Remove traversal restrictions from the target link |
| `close_link` | Active | Add a traversal restriction to the target link |
| `give_item` | Active | Add the target item to the player's inventory |
| `remove_item` | Active | Remove the target item from the player's inventory |
| `summon_creature` | Active | Spawn the target creature at the player's current location |

---

### Item Placement

Defines where items are initially found in the world when a run starts.

| Field | Description |
|---|---|
| Item | The item to place |
| Location | The location to place it at |
| Initial count | How many of this item are placed at this location |

---

### Creature

Creature definitions — monsters and NPCs. Includes combat stats, behaviour, and how the creature is handled after death.

| Field | Description |
|---|---|
| Name | Display name |
| Description | Narrative description |
| Attack damage | Base damage dealt per attack |
| Defense | Damage absorbed from incoming attacks |
| Disposition | Behavioural disposition (see below) |
| Max health | Maximum health points |
| Attack method | How the creature attacks — used in narrative descriptions |
| Attack description | Narrative describing the creature's attack |
| Body decay turns | How many turns after death the creature's corpse remains visible on encounter sheets |
| Respawn turns | How many turns after death before a new instance of this creature spawns at its placement location (0 = no respawn) |

**Disposition values:**

| Value | Combat behaviour |
|---|---|
| `aggressive` | Attacks the player when encountered; retaliates without being provoked; deals damage when the player flees |
| `inquisitive` | Does not attack first; becomes provoked if attacked |
| `indifferent` | Does not attack; does not retaliate even if attacked |

**Attack method values** (narrative only — no mechanical effect):

`claws`, `bite`, `sting`, `weapon`, `spell`, `slam`, `touch`, `breath`, `gaze`

A portrait image can be uploaded for each creature.

---

### Creature Placement

Defines where creatures initially appear in the world when a run starts.

| Field | Description |
|---|---|
| Creature | The creature to place |
| Location | The location to place it at |
| Initial count | How many of this creature are placed at this location |

---

## Turn Sheets

Each turn a character receives a set of turn sheets to fill out. Sheets are presented to the player in a specific order, and processed by the game engine in a different order.

| Sheet | Processing order | Presentation order | Notes |
|---|---|---|---|
| Join game | — | — | Sent when a player first joins; handled separately from regular turn processing |
| Inventory management | 1st | 2nd | Processed first; taking inventory actions forfeits combat that turn |
| Creature encounter | 2nd | 1st | Shown first so players see what they are facing before deciding on items |
| Location choice | 3rd | 3rd | Movement is processed last so the flee penalty uses the final creature state |
| Combat | — | — | Reserved — not yet available |
| Puzzle | — | — | Reserved — not yet available |

### Turn Sheet Background Images

When uploading a background image for an adventure game, select the sheet type the image should apply to.

| Sheet | Description |
|---|---|
| `adventure_game_join_game` | Join game sheet — required; sent when a player first joins |
| `adventure_game_location_choice` | Movement and object interaction sheet |
| `adventure_game_inventory_management` | Item management sheet |
| `adventure_game_monster` | Creature encounter sheet |
| `adventure_game_combat` | _(reserved — not yet available)_ |
| `adventure_game_puzzle` | _(reserved — not yet available)_ |

---

## Player Rules

### Character Stats

| Stat | Value | Description |
|---|---|---|
| Starting health | 100 | Health assigned when a character joins a run |
| Respawn health | 50 | Health restored after a character dies |
| Unarmed attack damage | 5 | Damage dealt when no weapon is equipped |
| Maximum health | 100 | Health cannot exceed this value from healing |

---

### Location Choice Sheet

Players choose where to move or interact with an object at their current location. Moving and interacting with an object are mutually exclusive — only one can be submitted per turn.

**Movement rules:**
- The player selects one of the exits shown on their sheet
- Available exits are pre-filtered by link visibility and traversal requirements
- The character moves to the chosen location
- If a traversal description is set on the link, it is shown in the turn narrative

**Flee penalty:**
- If the character moves away from a location where aggressive creatures are alive, each aggressive creature makes a free attack
- Flee damage = creature attack damage minus the character's armour defence, minimum 1
- Indifferent and inquisitive creatures do not deal flee damage

**Object interaction:**
- Players can choose to act on a visible object at their location instead of moving
- Available actions are determined by the object's current state and the effects defined on that state
- If an effect requires a specific item, the player must have that item in their inventory (unused)
- All matching effects for the chosen action fire at the same time

---

### Creature Encounter Sheet

Players submit combat actions for creatures present at their location. This sheet is omitted if no creatures (alive or recently dead) are present.

Players submit up to 3 combat actions per turn. Each action is either do nothing or attack, with a choice of target creature.

**Sheet display:**
- Alive creatures are always listed
- Dead creatures remain visible for the number of turns set by the creature's body decay turns setting
- If only corpses are present the sheet is read-only
- Up to 3 combat actions are available when the sheet is interactive

**Combat forfeiture:** if the player picks up, drops, equips, or unequips items on their inventory sheet that turn, all combat is skipped and the encounter sheet explains why.

**Attack resolution** (processed in order per action):
1. Player attack damage = weapon damage from equipped weapon minus the creature's defence, minimum 1
2. Non-aggressive creatures are provoked on the first hit
3. If the creature's health reaches 0: the creature is killed; any items it was carrying drop to the player's current location
4. If the creature survives: it retaliates if aggressive or if it was provoked this encounter
   - Retaliation damage = creature attack damage minus character armour defence, minimum 1

**Character death:**
- If character health reaches 0 the character is moved to the starting location
- Health is restored to the respawn health (50)
- A narrative event is generated

---

### Inventory Management Sheet

Players manage their carried items — picking up items from the floor, dropping items, equipping and unequipping gear, and using consumables.

Players select items to pick up, drop, equip, unequip, or use.

**Processing order within a single turn:** unequip → drop → pick up → equip → use

**Key rules:**
- **No ground items when aggressive creatures are present:** if alive aggressive creatures are at the location, ground items are not shown and pickup is not offered
- **Auto-pickup on equip:** if a player equips an item that is on the ground at their current location, it is automatically picked up first
- **Using items:** a consumable can only be used if it has uses remaining; uses are decremented on each use and the item is marked as exhausted when all uses are spent
