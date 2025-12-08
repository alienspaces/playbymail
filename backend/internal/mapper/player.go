package mapper

import "gitlab.com/alienspaces/playbymail/schema/api/player_schema"

func MapVerifyGameSubscriptionTokenResponse(token string) *player_schema.VerifyGameSubscriptionTokenResponse {
	return &player_schema.VerifyGameSubscriptionTokenResponse{SessionToken: token}
}
