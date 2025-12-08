package player_schema

// RequestGameSubscriptionTokenRequest maps to the /request-game-subscription-token request schema
type RequestGameSubscriptionTokenRequest struct {
	Email string `json:"email"`
}

type RequestGameSubscriptionTokenResponse struct {
	Message string `json:"message"`
}

// VerifyGameSubscriptionTokenRequest maps to the /verify-game-subscription-token request schema
type VerifyGameSubscriptionTokenRequest struct {
	Email          string `json:"email"`
	TurnSheetToken string `json:"turn_sheet_token"`
}

type VerifyGameSubscriptionTokenResponse struct {
	SessionToken string `json:"session_token"`
}


