package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/schema/api/catalog_schema"
)

// CatalogSubscriptionToData maps a manager subscription, its game, and linked instances
// into a single catalog entry with aggregated capacity and delivery flags.
func CatalogSubscriptionToData(
	l logger.Logger,
	sub *game_record.GameSubscription,
	gameRec *game_record.Game,
	instances []*game_record.GameInstance,
	playerCounts map[string]int,
) *catalog_schema.CatalogSubscriptionData {
	var totalCapacity, totalPlayers int
	var deliveryPost, deliveryLocal, deliveryEmail bool

	for _, inst := range instances {
		totalCapacity += inst.RequiredPlayerCount
		totalPlayers += playerCounts[inst.ID]
		deliveryPost = deliveryPost || inst.DeliveryPhysicalPost
		deliveryLocal = deliveryLocal || inst.DeliveryPhysicalLocal
		deliveryEmail = deliveryEmail || inst.DeliveryEmail
	}

	return &catalog_schema.CatalogSubscriptionData{
		GameSubscriptionID:    sub.ID,
		GameName:              gameRec.Name,
		GameDescription:       gameRec.Description,
		GameType:              gameRec.GameType,
		TurnDurationHours:     gameRec.TurnDurationHours,
		TotalCapacity:         totalCapacity,
		TotalPlayers:          totalPlayers,
		DeliveryPhysicalPost:  deliveryPost,
		DeliveryPhysicalLocal: deliveryLocal,
		DeliveryEmail:         deliveryEmail,
	}
}
