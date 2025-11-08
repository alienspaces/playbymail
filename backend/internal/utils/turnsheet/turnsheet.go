package turnsheet

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// TurnSheetCodeType represents the encoding format for a turn sheet code
// "playing" codes include full identifiers encoded as base64 JSON with checksum validation
// "joining" codes are simple identifiers that contain only the game ID (used for join sheets)
type TurnSheetCodeType string

const (
	TurnSheetCodeTypePlayingGame TurnSheetCodeType = "playing"
	TurnSheetCodeTypeJoiningGame TurnSheetCodeType = "joining"
)

// TurnSheetIdentifier contains the unique identifiers for a turn sheet
type TurnSheetIdentifier struct {
	CodeType        TurnSheetCodeType `json:"code_type"`
	GameID          string            `json:"game_id"`
	GameInstanceID  string            `json:"game_instance_id"`
	AccountID       string            `json:"account_id"`
	GameTurnSheetID string            `json:"game_turn_sheet_id"`
	GeneratedAt     time.Time         `json:"generated_at"`
	Checksum        string            `json:"checksum"`
}

// GenerateTurnSheetCode creates a unique, scannable code for a turn sheet
func GenerateTurnSheetCode(gameID, gameInstanceID, accountID, gameTurnSheetID string) (string, error) {
	identifier := TurnSheetIdentifier{
		CodeType:        TurnSheetCodeTypePlayingGame,
		GameID:          gameID,
		GameInstanceID:  gameInstanceID,
		AccountID:       accountID,
		GameTurnSheetID: gameTurnSheetID,
		GeneratedAt:     time.Now(),
	}

	// Create a checksum for integrity verification
	checksum, err := generateChecksum(identifier)
	if err != nil {
		return "", fmt.Errorf("failed to generate checksum: %w", err)
	}
	identifier.Checksum = checksum

	// Encode as JSON
	jsonData, err := json.Marshal(identifier)
	if err != nil {
		return "", fmt.Errorf("failed to marshal identifier: %w", err)
	}

	// Base64 encode for compact representation
	encoded := base64.URLEncoding.EncodeToString(jsonData)
	return encoded, nil
}

// GenerateJoinTurnSheetCode returns a structured code for join-game sheets using the same
// encoding format as playing-game turn sheets but only including the game identifier.
func GenerateJoinTurnSheetCode(gameID string) (string, error) {
	if gameID == "" {
		return "", fmt.Errorf("gameID is required")
	}
	if _, err := uuid.Parse(gameID); err != nil {
		return "", fmt.Errorf("invalid gameID format: %w", err)
	}

	identifier := TurnSheetIdentifier{
		CodeType:        TurnSheetCodeTypeJoiningGame,
		GameID:          gameID,
		GeneratedAt:     time.Now(),
		GameInstanceID:  "",
		AccountID:       "",
		GameTurnSheetID: "",
	}

	checksum, err := generateChecksum(identifier)
	if err != nil {
		return "", fmt.Errorf("failed to generate checksum: %w", err)
	}
	identifier.Checksum = checksum

	jsonData, err := json.Marshal(identifier)
	if err != nil {
		return "", fmt.Errorf("failed to marshal identifier: %w", err)
	}

	return base64.URLEncoding.EncodeToString(jsonData), nil
}

// ParseTurnSheetCode decodes and validates a turn sheet code
func ParseTurnSheetCode(code string) (*TurnSheetIdentifier, error) {
	jsonData, err := base64.URLEncoding.DecodeString(code)
	if err != nil {
		return nil, fmt.Errorf("failed to decode turn sheet code: %w", err)
	}

	var identifier TurnSheetIdentifier
	if err := json.Unmarshal(jsonData, &identifier); err != nil {
		return nil, fmt.Errorf("failed to unmarshal identifier: %w", err)
	}

	if identifier.CodeType != TurnSheetCodeTypePlayingGame && identifier.CodeType != TurnSheetCodeTypeJoiningGame {
		return nil, fmt.Errorf("unsupported turn sheet code type: %s", identifier.CodeType)
	}

	expectedChecksum, err := generateChecksum(identifier)
	if err != nil {
		return nil, fmt.Errorf("failed to generate expected checksum: %w", err)
	}

	if identifier.Checksum != expectedChecksum {
		return nil, fmt.Errorf("checksum verification failed")
	}

	return &identifier, nil
}

// generateChecksum creates an HMAC-SHA256 checksum for the identifier
func generateChecksum(identifier TurnSheetIdentifier) (string, error) {
	// Use a secret key (in production, this should come from config)
	secretKey := "playbymail-turn-sheet-secret-key"

	// Create data for checksum (exclude the checksum field itself)
	data := fmt.Sprintf("%s:%s:%s:%s:%s:%d",
		identifier.CodeType,
		identifier.GameID,
		identifier.GameInstanceID,
		identifier.AccountID,
		identifier.GameTurnSheetID,
		identifier.GeneratedAt.Unix())

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))

	return base64.URLEncoding.EncodeToString(h.Sum(nil)), nil
}

// GenerateQRCodeData creates QR code data for the turn sheet
func GenerateQRCodeData(gameID, gameInstanceID, accountID, gameTurnSheetID string) (string, error) {
	// QR codes work well with shorter data, so we can use a shorter format
	// Format: "PBM:gameID:gameInstanceID:accountID:gameTurnSheetID"
	shortCode := fmt.Sprintf("PBM:%s:%s:%s:%s", gameID, gameInstanceID, accountID, gameTurnSheetID)
	return shortCode, nil
}

// ParseQRCodeData parses QR code data back to identifiers
func ParseQRCodeData(qrData string) (gameID, gameInstanceID, accountID, gameTurnSheetID string, err error) {
	// Check if it's our format
	if len(qrData) < 4 || qrData[:4] != "PBM:" {
		return "", "", "", "", fmt.Errorf("invalid QR code format")
	}

	// Split by colon
	// Note: This is a simplified parser - in production you'd want more robust parsing
	// For now, we'll assume the format is consistent

	// This is a simplified implementation - you'd need to handle the actual parsing
	// based on your specific QR code format
	return "", "", "", "", fmt.Errorf("QR code parsing not fully implemented")
}
