package adventure_game_record

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/record"
)

const TableAdventureGameCreature = "adventure_game_creature"

const (
	FieldAdventureGameCreatureID                = "id"
	FieldAdventureGameCreatureGameID            = "game_id"
	FieldAdventureGameCreatureName              = "name"
	FieldAdventureGameCreatureDescription       = "description"
	FieldAdventureGameCreatureAttackDamage      = "attack_damage"
	FieldAdventureGameCreatureDefense           = "defense"
	FieldAdventureGameCreatureDisposition       = "disposition"
	FieldAdventureGameCreatureAttackMethod      = "attack_method"
	FieldAdventureGameCreatureAttackDescription = "attack_description"
	FieldAdventureGameCreatureMaxHealth         = "max_health"
	FieldAdventureGameCreatureBodyDecayTurns    = "body_decay_turns"
	FieldAdventureGameCreatureRespawnTurns      = "respawn_turns"
)

const (
	AdventureGameCreatureDispositionAggressive  = "aggressive"
	AdventureGameCreatureDispositionInquisitive = "inquisitive"
	AdventureGameCreatureDispositionIndifferent = "indifferent"
)

const (
	AdventureGameCreatureAttackMethodClaws   = "claws"
	AdventureGameCreatureAttackMethodBite    = "bite"
	AdventureGameCreatureAttackMethodSting   = "sting"
	AdventureGameCreatureAttackMethodWeapon  = "weapon"
	AdventureGameCreatureAttackMethodSpell   = "spell"
	AdventureGameCreatureAttackMethodSlam    = "slam"
	AdventureGameCreatureAttackMethodTouch   = "touch"
	AdventureGameCreatureAttackMethodBreath  = "breath"
	AdventureGameCreatureAttackMethodGaze    = "gaze"
)

type AdventureGameCreature struct {
	record.Record
	GameID            string `db:"game_id"`
	Name              string `db:"name"`
	Description       string `db:"description"`
	AttackDamage      int    `db:"attack_damage"`
	Defense           int    `db:"defense"`
	Disposition       string `db:"disposition"`
	MaxHealth         int    `db:"max_health"`
	AttackMethod      string `db:"attack_method"`
	AttackDescription string `db:"attack_description"`
	BodyDecayTurns    int    `db:"body_decay_turns"`
	RespawnTurns      int    `db:"respawn_turns"`
}

func (r *AdventureGameCreature) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldAdventureGameCreatureGameID] = r.GameID
	args[FieldAdventureGameCreatureName] = r.Name
	args[FieldAdventureGameCreatureDescription] = r.Description
	args[FieldAdventureGameCreatureAttackDamage] = r.AttackDamage
	args[FieldAdventureGameCreatureDefense] = r.Defense
	args[FieldAdventureGameCreatureDisposition] = r.Disposition
	args[FieldAdventureGameCreatureMaxHealth] = r.MaxHealth
	args[FieldAdventureGameCreatureAttackMethod] = r.AttackMethod
	args[FieldAdventureGameCreatureAttackDescription] = r.AttackDescription
	args[FieldAdventureGameCreatureBodyDecayTurns] = r.BodyDecayTurns
	args[FieldAdventureGameCreatureRespawnTurns] = r.RespawnTurns
	return args
}
