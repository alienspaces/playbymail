package adventure_game_character

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/repository"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/repositor"
	adventure_game_record "gitlab.com/alienspaces/playbymail/internal/record/adventure_game"
)

const TableName = adventure_game_record.TableAdventureGameCharacter

// NewRepository matches the RepositoryConstructor signature
func NewRepository(l logger.Logger, tx pgx.Tx) (repositor.Repositor, error) {
	return repository.NewGeneric[adventure_game_record.AdventureGameCharacter](repository.NewArgs{
		Tx:        tx,
		TableName: TableName,
		Record:    adventure_game_record.AdventureGameCharacter{},
	})
}
