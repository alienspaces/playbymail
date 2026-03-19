package account

import (
	"net/http"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
)

// authorizeAccountRead verifies the request is authenticated.
// Used by handlers that self-scope to the authenticated user's own data.
func authorizeAccountRead(l logger.Logger, r *http.Request) (*server.AuthenData, error) {
	authenData := server.GetRequestAuthenData(l, r)
	if authenData == nil || authenData.AccountUser.ID == "" {
		l.Warn("authenticated account is required")
		return nil, coreerror.NewUnauthorizedError()
	}
	return authenData, nil
}

// authorizeAccountMember verifies the authenticated user belongs to the account
// identified by accountID. Used by handlers operating on a specific account.
func authorizeAccountMember(l logger.Logger, r *http.Request, accountID string) (*server.AuthenData, error) {
	authenData, err := authorizeAccountRead(l, r)
	if err != nil {
		return nil, err
	}

	if authenData.AccountUser.AccountID != accountID {
		l.Warn("authenticated account >%s< does not match requested account >%s<",
			authenData.AccountUser.AccountID, accountID)
		return nil, coreerror.NewUnauthorizedError()
	}

	return authenData, nil
}

// authorizeAccountUserSelf verifies the authenticated user belongs to the account
// AND is the specific account user identified by accountUserID.
// Used by handlers operating on a specific account user's resources.
func authorizeAccountUserSelf(l logger.Logger, r *http.Request, accountID, accountUserID string) (*server.AuthenData, error) {
	authenData, err := authorizeAccountMember(l, r, accountID)
	if err != nil {
		return nil, err
	}

	if authenData.AccountUser.ID != accountUserID {
		l.Warn("authenticated user >%s< does not match requested account user >%s<",
			authenData.AccountUser.ID, accountUserID)
		return nil, coreerror.NewUnauthorizedError()
	}

	return authenData, nil
}
