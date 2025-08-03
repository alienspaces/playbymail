package account_test

import (
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/core/type/storer"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/runner/server/account"
	"gitlab.com/alienspaces/playbymail/internal/utils/testutil"
	"gitlab.com/alienspaces/playbymail/schema/api/account_schema"
)

func Test_getAccountHandler(t *testing.T) {
	t.Parallel()

	th := testutil.NewTestHarness(t)
	require.NotNil(t, th, "newTestHarness returns without error")

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	type testCase struct {
		testutil.TestCase
		collectionRequest     bool
		collectionRecordCount int
	}

	testCaseCollectionResponseDecoder := testutil.TestCaseResponseDecoderGeneric[account_schema.AccountCollectionResponse]
	testCaseResponseDecoder := testutil.TestCaseResponseDecoderGeneric[account_schema.AccountResponse]

	// Setup: get an account for reference
	accountRec, err := th.Data.GetAccountRecByRef(harness.AccountOneRef)
	require.NoError(t, err, "GetAccountRecByRef returns without error")

	testCases := []testCase{
		{
			TestCase: testutil.TestCase{
				Name: "API key with open access \\ get many accounts \\ returns expected accounts",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[account.GetManyAccounts]
				},
				RequestQueryParams: func(d harness.Data) map[string]any {
					return map[string]any{
						"id":          accountRec.ID,
						"page_size":   10,
						"page_number": 1,
					}
				},
				ResponseDecoder: testCaseCollectionResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
			collectionRequest:     true,
			collectionRecordCount: 1,
		},
		{
			TestCase: testutil.TestCase{
				Name: "API key with open access \\ get one account with valid account ID \\ returns expected account",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[account.GetOneAccount]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					accountRec, err := d.GetAccountRecByRef(harness.AccountOneRef)
					require.NoError(t, err, "GetAccountRecByRef returns without error")
					params := map[string]string{
						":account_id": accountRec.ID,
					}
					return params
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
			collectionRequest: false,
		},
	}

	for _, testCase := range testCases {

		t.Logf("Running test >%s<\n", testCase.Name)

		t.Run(testCase.Name, func(t *testing.T) {
			testFunc := func(method string, body interface{}) {
				if testCase.TestResponseCode() != http.StatusOK {
					return
				}

				require.NotNil(t, body, "Response body is not nil")

				var responses []*account_schema.AccountResponseData
				if testCase.collectionRequest {
					responses = body.(account_schema.AccountCollectionResponse).Data
				} else {
					responses = append(responses, body.(account_schema.AccountResponse).Data)
				}

				if testCase.collectionRequest {
					require.Equal(t, testCase.collectionRecordCount, len(responses), "Response record count length equals expected")
				}

				if testCase.collectionRequest && testCase.collectionRecordCount == 0 {
					require.Empty(t, responses, "Response body should be empty")
				} else {
					require.NotEmpty(t, responses, "Response body is not empty")
				}

				for _, d := range responses {
					require.NotEmpty(t, d.ID, "Account ID is not empty")
					require.NotEmpty(t, d.Email, "Account Email is not empty")
				}
			}

			testutil.RunTestCase(t, th, &testCase, testFunc)
		})
	}
}

func Test_createUpdateDeleteAccountHandler(t *testing.T) {
	t.Parallel()

	th := testutil.NewTestHarness(t)
	require.NotNil(t, th, "newTestHarness returns without error")

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	type testCase struct {
		testutil.TestCase
		expectResponse func(d harness.Data, req account_schema.AccountRequest) account_schema.AccountResponse
	}

	testCaseResponseDecoder := testutil.TestCaseResponseDecoderGeneric[account_schema.AccountResponse]

	testCases := []testCase{
		{
			TestCase: testutil.TestCase{
				Name: "API key with open access \\ create account with valid properties \\ returns created account",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[account.CreateAccount]
				},
				RequestBody: func(d harness.Data) any {
					email := gofakeit.Email()
					return account_schema.AccountRequest{
						Email: &email,
						Name:  gofakeit.Name(),
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusCreated,
			},
			expectResponse: func(d harness.Data, req account_schema.AccountRequest) account_schema.AccountResponse {
				return account_schema.AccountResponse{
					Data: &account_schema.AccountResponseData{
						Email: *req.Email,
						Name:  req.Name,
					},
				}
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "API key with open access \\ update account with valid properties \\ returns updated account",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[account.UpdateAccount]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					accountRec, err := d.GetAccountRecByRef(harness.AccountOneRef)
					require.NoError(t, err, "GetAccountRecByRef returns without error")
					params := map[string]string{
						":account_id": accountRec.ID,
					}
					return params
				},
				RequestBody: func(d harness.Data) any {
					email := harness.UniqueEmail(gofakeit.Email())
					return account_schema.AccountRequest{
						Email: &email,
						Name:  gofakeit.Name(),
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
			expectResponse: func(d harness.Data, req account_schema.AccountRequest) account_schema.AccountResponse {
				// Get the original account to get the original email (email cannot be updated)
				accountRec, err := d.GetAccountRecByRef(harness.AccountOneRef)
				require.NoError(t, err, "GetAccountRecByRef returns without error")

				return account_schema.AccountResponse{
					Data: &account_schema.AccountResponseData{
						Email: accountRec.Email, // Email should remain unchanged
						Name:  req.Name,
					},
				}
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "API key with open access \\ delete account with valid account ID \\ returns no content",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[account.DeleteAccount]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					accountRec, err := d.GetAccountRecByRef(harness.AccountOneRef)
					require.NoError(t, err, "GetAccountRecByRef returns without error")
					params := map[string]string{
						":account_id": accountRec.ID,
					}
					return params
				},
				ResponseCode: http.StatusNoContent,
			},
			expectResponse: nil,
		},
	}

	for _, testCase := range testCases {
		t.Logf("Running test >%s<", testCase.Name)

		t.Run(testCase.Name, func(t *testing.T) {
			testFunc := func(method string, body interface{}) {
				if testCase.TestResponseCode() == http.StatusNoContent {
					// No content expected
					return
				}

				require.NotNil(t, body, "Response body is not nil")
				t.Logf("Actual response body: %#v", body)
				resp, ok := body.(account_schema.AccountResponse)
				require.True(t, ok, "Response body is of type account_schema.AccountResponse")
				aResp := resp.Data
				require.NotNil(t, aResp, "AccountResponseData is not nil")
				t.Logf("AccountResponseData: %#v", aResp)
				require.NotEmpty(t, aResp.ID, "Account ID is not empty")
				require.NotEmpty(t, aResp.Email, "Account Email is not empty")
				xResp := testCase.expectResponse(
					th.Data,
					testCase.TestRequestBody(th.Data).(account_schema.AccountRequest),
				).Data
				require.Equal(t, xResp.Email, aResp.Email, "Account Email matches expected")
				require.Equal(t, xResp.Name, aResp.Name, "Account Name equals expected")
				require.False(t, aResp.CreatedAt.IsZero(), "Account CreatedAt is not zero")
				// UpdatedAt is allowed to be nil, so do not assert on it
			}

			testutil.RunTestCase(t, th, &testCase, testFunc)
		})
	}
}

func Test_accountMeHandler(t *testing.T) {
	t.Parallel()

	th := testutil.NewTestHarness(t)
	require.NotNil(t, th, "newTestHarness returns without error")

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	testCaseResponseDecoder := testutil.TestCaseResponseDecoderGeneric[account_schema.AccountResponse]

	// Setup: get an account for reference
	accountRec, err := th.Data.GetAccountRecByRef(harness.AccountOneRef)
	require.NoError(t, err, "GetAccountRecByRef returns without error")

	testCases := []testutil.TestCase{
		{
			Name: "API key with open access \\ get my account \\ returns authenticated user account",
			NewRunner: func(l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx]) (testutil.TestRunnerer, error) {
				rnr, err := testutil.NewTestRunner(l, s, j)
				if err != nil {
					return nil, err
				}

				// Mock authentication to return the test account
				rnr.AuthenticateRequestFunc = func(l logger.Logger, m domainer.Domainer, r *http.Request, authType server.AuthenticationType) (server.AuthenData, error) {
					return server.AuthenData{
						Type: server.AuthenticatedTypeToken,
						Account: server.AuthenticatedAccount{
							ID:    accountRec.ID,
							Name:  accountRec.Name,
							Email: accountRec.Email,
						},
					}, nil
				}

				return rnr, nil
			},
			HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
				return rnr.GetHandlerConfig()[account.GetMyAccount]
			},
			ResponseDecoder: testCaseResponseDecoder,
			ResponseCode:    http.StatusOK,
		},
		{
			Name: "API key with open access \\ update my account \\ updates authenticated user account name",
			NewRunner: func(l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx]) (testutil.TestRunnerer, error) {
				rnr, err := testutil.NewTestRunner(l, s, j)
				if err != nil {
					return nil, err
				}

				// Mock authentication to return the test account
				rnr.AuthenticateRequestFunc = func(l logger.Logger, m domainer.Domainer, r *http.Request, authType server.AuthenticationType) (server.AuthenData, error) {
					return server.AuthenData{
						Type: server.AuthenticatedTypeToken,
						Account: server.AuthenticatedAccount{
							ID:    accountRec.ID,
							Name:  accountRec.Name,
							Email: accountRec.Email,
						},
					}, nil
				}

				return rnr, nil
			},
			HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
				return rnr.GetHandlerConfig()[account.UpdateMyAccount]
			},
			RequestBody: func(d harness.Data) any {
				return account_schema.AccountRequest{
					Name: "Updated Name",
				}
			},
			ResponseDecoder: testCaseResponseDecoder,
			ResponseCode:    http.StatusOK,
		},
		{
			Name: "API key with open access \\ delete my account \\ deletes authenticated user account",
			NewRunner: func(l logger.Logger, s storer.Storer, j *river.Client[pgx.Tx]) (testutil.TestRunnerer, error) {
				rnr, err := testutil.NewTestRunner(l, s, j)
				if err != nil {
					return nil, err
				}

				// Mock authentication to return the test account
				rnr.AuthenticateRequestFunc = func(l logger.Logger, m domainer.Domainer, r *http.Request, authType server.AuthenticationType) (server.AuthenData, error) {
					return server.AuthenData{
						Type: server.AuthenticatedTypeToken,
						Account: server.AuthenticatedAccount{
							ID:    accountRec.ID,
							Name:  accountRec.Name,
							Email: accountRec.Email,
						},
					}, nil
				}

				return rnr, nil
			},
			HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
				return rnr.GetHandlerConfig()[account.DeleteMyAccount]
			},
			ResponseCode: http.StatusNoContent,
		},
	}

	for _, testCase := range testCases {
		t.Logf("Running test >%s<", testCase.Name)

		t.Run(testCase.Name, func(t *testing.T) {
			testFunc := func(method string, body interface{}) {
				if testCase.TestResponseCode() == http.StatusNoContent {
					// No content expected
					return
				}

				require.NotNil(t, body, "Response body is not nil")
				t.Logf("Actual response body: %#v", body)
				resp, ok := body.(account_schema.AccountResponse)
				require.True(t, ok, "Response body is of type account_schema.AccountResponse")
				aResp := resp.Data
				require.NotNil(t, aResp, "AccountResponseData is not nil")
				t.Logf("AccountResponseData: %#v", aResp)
				require.NotEmpty(t, aResp.ID, "Account ID is not empty")
				require.NotEmpty(t, aResp.Email, "Account Email is not empty")
				require.Equal(t, accountRec.ID, aResp.ID, "Account ID matches authenticated user")
				require.Equal(t, accountRec.Email, aResp.Email, "Account Email matches authenticated user")

				if method == http.MethodPut {
					require.Equal(t, "Updated Name", aResp.Name, "Account name is updated")
				} else {
					require.Equal(t, accountRec.Name, aResp.Name, "Account name matches authenticated user")
				}

				require.False(t, aResp.CreatedAt.IsZero(), "Account CreatedAt is not zero")
			}

			testutil.RunTestCase(t, th, &testCase, testFunc)
		})
	}
}
