package turnsheetutil

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type TurnSheetCodeType string

const (
	TurnSheetCodeTypeJoiningGame TurnSheetCodeType = "join"
	TurnSheetCodeTypePlayingGame TurnSheetCodeType = "play"
)

type TurnSheetCodeData struct {
	CodeType TurnSheetCodeType `json:"code_type"`
}

func generateCode(v any) (string, error) {
	json, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(json), nil
}

// A join game turn sheet code only needs to contain the following data to identify the specific
// game the player wishes to join. The account record ID and the game subscription record ID
// belong to the game manager who is running the game.
type JoinGameTurnSheetCodeData struct {
	TurnSheetCodeData
	GameSubscriptionID string `json:"game_subscription_id"` // game_subscription_id
}

// A play game turn sheet code only needs to contain the following data to identify the specific
// game turn sheet the player wishes to play. The account record ID and the game turn sheet record ID
// belong to the player who is playing the game.
type PlayGameTurnSheetCodeData struct {
	TurnSheetCodeData
	GameTurnSheetID string `json:"game_turn_sheet_id"` // game_turn_sheet_id
}

func ParseTurnSheetCodeTypeFromCode(code string) (TurnSheetCodeType, error) {
	// Decode base64
	decoded, err := base64.RawURLEncoding.DecodeString(code)
	if err != nil {
		return "", fmt.Errorf("failed to decode turn sheet code: %w", err)
	}

	// Unmarshal JSON
	var turnSheetCodeData TurnSheetCodeData
	if err := json.Unmarshal(decoded, &turnSheetCodeData); err != nil {
		return "", fmt.Errorf("failed to unmarshal turn sheet code: %w", err)
	}
	return turnSheetCodeData.CodeType, nil
}

// GenerateJoinGameTurnSheetCode generates a join game turn sheet code
func GenerateJoinGameTurnSheetCode(gameSubscriptionID string) (string, error) {
	return generateCode(&JoinGameTurnSheetCodeData{
		TurnSheetCodeData: TurnSheetCodeData{
			CodeType: TurnSheetCodeTypeJoiningGame,
		},
		GameSubscriptionID: gameSubscriptionID,
	})
}

// GeneratePlayGameTurnSheetCode generates a play game turn sheet code
func GeneratePlayGameTurnSheetCode(gameTurnSheetID string) (string, error) {
	return generateCode(&PlayGameTurnSheetCodeData{
		TurnSheetCodeData: TurnSheetCodeData{
			CodeType: TurnSheetCodeTypePlayingGame,
		},
		GameTurnSheetID: gameTurnSheetID,
	})
}

// ParseJoinGameTurnSheetCodeData decodes and validates a join game turn sheet code
func ParseJoinGameTurnSheetCodeData(code string) (*JoinGameTurnSheetCodeData, error) {
	// Decode base64
	decoded, err := base64.RawURLEncoding.DecodeString(code)
	if err != nil {
		return nil, fmt.Errorf("failed to decode turn sheet code: %w", err)
	}

	// Unmarshal JSON
	var joinGameTurnSheetCodeData JoinGameTurnSheetCodeData
	if err := json.Unmarshal(decoded, &joinGameTurnSheetCodeData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal turn sheet code: %w", err)
	}

	// Validate code type
	if joinGameTurnSheetCodeData.CodeType != TurnSheetCodeTypeJoiningGame {
		return nil, fmt.Errorf("incorrect turn sheet code type: %s", joinGameTurnSheetCodeData.CodeType)
	}

	return &joinGameTurnSheetCodeData, nil
}

// ParsePlayGameTurnSheetCodeData decodes and validates a play game turn sheet code
func ParsePlayGameTurnSheetCodeData(code string) (*PlayGameTurnSheetCodeData, error) {
	// Decode base64
	decoded, err := base64.RawURLEncoding.DecodeString(code)
	if err != nil {
		return nil, fmt.Errorf("failed to decode turn sheet code: %w", err)
	}

	// Unmarshal JSON
	var playGameTurnSheetCodeData PlayGameTurnSheetCodeData
	if err := json.Unmarshal(decoded, &playGameTurnSheetCodeData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal turn sheet code: %w", err)
	}

	// Validate code type
	if playGameTurnSheetCodeData.CodeType != TurnSheetCodeTypePlayingGame {
		return nil, fmt.Errorf("incorrect turn sheet code type: %s", playGameTurnSheetCodeData.CodeType)
	}

	return &playGameTurnSheetCodeData, nil
}
