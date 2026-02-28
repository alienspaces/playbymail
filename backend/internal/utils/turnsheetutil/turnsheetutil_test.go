package turnsheetutil

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateCode_Subtypes(t *testing.T) {
	t.Run("JoinGameTurnSheetCodeData", func(t *testing.T) {
		data := JoinGameTurnSheetCodeData{
			TurnSheetCodeData:  TurnSheetCodeData{CodeType: TurnSheetCodeTypeJoiningGame},
			GameSubscriptionID: "sub-456",
		}

		code, err := generateCode(&data)
		require.NoError(t, err)

		// Decode and check if fields are present
		decodedBytes, err := base64.RawURLEncoding.DecodeString(code)
		require.NoError(t, err)

		var decodedMap map[string]interface{}
		err = json.Unmarshal(decodedBytes, &decodedMap)
		require.NoError(t, err)

		assert.Equal(t, string(TurnSheetCodeTypeJoiningGame), decodedMap["code_type"])
		assert.Equal(t, "sub-456", decodedMap["game_subscription_id"], "game_subscription_id should be present")
	})

	t.Run("PlayGameTurnSheetCodeData", func(t *testing.T) {
		data := PlayGameTurnSheetCodeData{
			TurnSheetCodeData: TurnSheetCodeData{CodeType: TurnSheetCodeTypePlayingGame},
			GameTurnSheetID:   "sheet-789",
		}

		code, err := generateCode(&data)
		require.NoError(t, err)

		// Decode and check if fields are present
		decodedBytes, err := base64.RawURLEncoding.DecodeString(code)
		require.NoError(t, err)

		var decodedMap map[string]interface{}
		err = json.Unmarshal(decodedBytes, &decodedMap)
		require.NoError(t, err)

		assert.Equal(t, string(TurnSheetCodeTypePlayingGame), decodedMap["code_type"])
		assert.Equal(t, "sheet-789", decodedMap["game_turn_sheet_id"], "game_turn_sheet_id should be present")
	})
}

func TestParseTurnSheetCodeTypeFromCode(t *testing.T) {
	t.Run("Success Join", func(t *testing.T) {
		data := TurnSheetCodeData{CodeType: TurnSheetCodeTypeJoiningGame}
		code, err := generateCode(&data)
		require.NoError(t, err)

		codeType, err := ParseTurnSheetCodeTypeFromCode(code)
		require.NoError(t, err)
		assert.Equal(t, TurnSheetCodeTypeJoiningGame, codeType)
	})

	t.Run("Success Play", func(t *testing.T) {
		data := TurnSheetCodeData{CodeType: TurnSheetCodeTypePlayingGame}
		code, err := generateCode(&data)
		require.NoError(t, err)

		codeType, err := ParseTurnSheetCodeTypeFromCode(code)
		require.NoError(t, err)
		assert.Equal(t, TurnSheetCodeTypePlayingGame, codeType)
	})

	t.Run("Invalid Base64", func(t *testing.T) {
		_, err := ParseTurnSheetCodeTypeFromCode("invalid-base64-%%%")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to decode turn sheet code")
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		// "inv" base64 decodes to junk bytes
		code := base64.RawURLEncoding.EncodeToString([]byte("{invalid-json}"))
		_, err := ParseTurnSheetCodeTypeFromCode(code)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unmarshal turn sheet code")
	})
}

func TestParseJoinGameTurnSheetCodeData(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		data := JoinGameTurnSheetCodeData{
			TurnSheetCodeData:  TurnSheetCodeData{CodeType: TurnSheetCodeTypeJoiningGame},
			GameSubscriptionID: "sub-456",
		}
		// Manually generate code to simulate real usage
		jsonData, _ := json.Marshal(data)
		code := base64.RawURLEncoding.EncodeToString(jsonData)

		result, err := ParseJoinGameTurnSheetCodeData(code)
		require.NoError(t, err)
		assert.Equal(t, "sub-456", result.GameSubscriptionID)
		assert.Equal(t, TurnSheetCodeTypeJoiningGame, result.CodeType)
	})

	t.Run("Incorrect Code Type", func(t *testing.T) {
		data := PlayGameTurnSheetCodeData{
			TurnSheetCodeData: TurnSheetCodeData{CodeType: TurnSheetCodeTypePlayingGame},
			GameTurnSheetID:   "sheet-789",
		}
		jsonData, _ := json.Marshal(data)
		code := base64.RawURLEncoding.EncodeToString(jsonData)

		_, err := ParseJoinGameTurnSheetCodeData(code)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "incorrect turn sheet code type")
	})

	t.Run("Invalid Base64", func(t *testing.T) {
		_, err := ParseJoinGameTurnSheetCodeData("invalid-base64")
		assert.Error(t, err)
	})
}

func TestParsePlayGameTurnSheetCodeData(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		data := PlayGameTurnSheetCodeData{
			TurnSheetCodeData: TurnSheetCodeData{CodeType: TurnSheetCodeTypePlayingGame},
			GameTurnSheetID:   "sheet-789",
		}
		jsonData, _ := json.Marshal(data)
		code := base64.RawURLEncoding.EncodeToString(jsonData)

		result, err := ParsePlayGameTurnSheetCodeData(code)
		require.NoError(t, err)
		assert.Equal(t, "sheet-789", result.GameTurnSheetID)
		assert.Equal(t, TurnSheetCodeTypePlayingGame, result.CodeType)
	})

	t.Run("Incorrect Code Type", func(t *testing.T) {
		data := JoinGameTurnSheetCodeData{
			TurnSheetCodeData:  TurnSheetCodeData{CodeType: TurnSheetCodeTypeJoiningGame},
			GameSubscriptionID: "sub-456",
		}
		jsonData, _ := json.Marshal(data)
		code := base64.RawURLEncoding.EncodeToString(jsonData)

		_, err := ParsePlayGameTurnSheetCodeData(code)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "incorrect turn sheet code type")
	})

	t.Run("Invalid Base64", func(t *testing.T) {
		_, err := ParsePlayGameTurnSheetCodeData("invalid-base64")
		assert.Error(t, err)
	})
}
