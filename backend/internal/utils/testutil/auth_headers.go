package testutil

import (
	"gitlab.com/alienspaces/playbymail/internal/harness"
)

// AuthHeaderByAccountRef returns Authorization header with session token for account by reference
func AuthHeaderByAccountRef(data harness.Data, ref string) map[string]string {
	token, err := data.GetAccountSessionTokenByAccountRef(ref)
	if err != nil {
		// Return empty header if token not found (tests should handle this)
		return map[string]string{}
	}
	return map[string]string{
		"Authorization": "Bearer " + token,
	}
}

// AuthHeaderByAccountID returns Authorization header with session token for account by ID
func AuthHeaderByAccountID(data harness.Data, accountID string) map[string]string {
	token, err := data.GetAccountSessionToken(accountID)
	if err != nil {
		// Return empty header if token not found (tests should handle this)
		return map[string]string{}
	}
	return map[string]string{
		"Authorization": "Bearer " + token,
	}
}

// AuthHeaderStandard returns Authorization header with standard test user token
// Uses StandardAccountRef to get the token from harness data
func AuthHeaderStandard(data harness.Data) map[string]string {
	return AuthHeaderByAccountRef(data, harness.StandardAccountRef)
}

// AuthHeaderProPlayer returns Authorization header with pro player test user token
// Uses ProPlayerAccountRef to get the token from harness data
func AuthHeaderProPlayer(data harness.Data) map[string]string {
	return AuthHeaderByAccountRef(data, harness.ProPlayerAccountRef)
}

// AuthHeaderProDesigner returns Authorization header with pro designer test user token
// Uses ProDesignerAccountRef to get the token from harness data
func AuthHeaderProDesigner(data harness.Data) map[string]string {
	return AuthHeaderByAccountRef(data, harness.ProDesignerAccountRef)
}

// AuthHeaderProManager returns Authorization header with pro manager test user token
// Uses ProManagerAccountRef to get the token from harness data
func AuthHeaderProManager(data harness.Data) map[string]string {
	return AuthHeaderByAccountRef(data, harness.ProManagerAccountRef)
}

// AuthHeaderWithToken returns Authorization header with custom token
func AuthHeaderWithToken(token string) map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + token,
	}
}
