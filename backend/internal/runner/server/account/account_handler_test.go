package account_test

import (
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/server"
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

	testCases := []testCase{
		{
			TestCase: testutil.TestCase{
				Name: "authenticated user when get many accounts then returns expected accounts",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[account.GetManyAccounts]
				},
				RequestHeaders: func(d harness.Data) map[string]string {
					return testutil.AuthHeaderStandard(d)
				},
				RequestQueryParams: func(d harness.Data) map[string]any {
					return map[string]any{
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
				Name: "authenticated user when get one account with valid account ID then returns expected account",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[account.GetOneAccount]
				},
				RequestHeaders: func(d harness.Data) map[string]string {
					return testutil.AuthHeaderStandard(d)
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					accountRec, err := d.GetAccountUserRecByRef(harness.StandardAccountRef)
					require.NoError(t, err, "GetAccountUserRecByRef returns without error")
					params := map[string]string{
						":account_id": accountRec.AccountID,
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
			testFunc := func(method string, body any) {
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
					require.NotEmpty(t, d.Status, "Account Status is not empty")
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
				Name: "public access when create account with valid properties then returns created account",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[account.CreateOneAccount]
				},
				RequestBody: func(d harness.Data) any {
					name := gofakeit.Company()
					return account_schema.AccountRequest{
						Name: &name,
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusCreated,
			},
			expectResponse: func(d harness.Data, req account_schema.AccountRequest) account_schema.AccountResponse {
				return account_schema.AccountResponse{
					Data: &account_schema.AccountResponseData{
						Name: *req.Name,
					},
				}
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "authenticated user when update account with valid properties then returns updated account",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[account.UpdateOneAccount]
				},
				RequestHeaders: func(d harness.Data) map[string]string {
					return testutil.AuthHeaderStandard(d)
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					accountRec, err := d.GetAccountUserRecByRef(harness.StandardAccountRef)
					require.NoError(t, err, "GetAccountUserRecByRef returns without error")
					params := map[string]string{
						":account_id": accountRec.AccountID,
					}
					return params
				},
				RequestBody: func(d harness.Data) any {
					name := gofakeit.Company()
					return account_schema.AccountRequest{
						Name: &name,
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
			expectResponse: func(d harness.Data, req account_schema.AccountRequest) account_schema.AccountResponse {
				return account_schema.AccountResponse{
					Data: &account_schema.AccountResponseData{
						Name: *req.Name,
					},
				}
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "authenticated user when delete account with valid account ID then returns no content",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[account.DeleteOneAccount]
				},
				RequestHeaders: func(d harness.Data) map[string]string {
					return testutil.AuthHeaderStandard(d)
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					accountRec, err := d.GetAccountUserRecByRef(harness.StandardAccountRef)
					require.NoError(t, err, "GetAccountUserRecByRef returns without error")
					params := map[string]string{
						":account_id": accountRec.AccountID,
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
			testFunc := func(method string, body any) {
				if testCase.TestResponseCode() == http.StatusNoContent {
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
				require.NotEmpty(t, aResp.Name, "Account Name is not empty")
				xResp := testCase.expectResponse(
					th.Data,
					testCase.TestRequestBody(th.Data).(account_schema.AccountRequest),
				).Data
				require.Equal(t, xResp.Name, aResp.Name, "Account Name matches expected")
				require.False(t, aResp.CreatedAt.IsZero(), "Account CreatedAt is not zero")
			}

			testutil.RunTestCase(t, th, &testCase, testFunc)
		})
	}
}

func Test_getMeHandler(t *testing.T) {
	t.Parallel()

	th := testutil.NewTestHarness(t)
	require.NotNil(t, th, "newTestHarness returns without error")

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	testCaseResponseDecoder := testutil.TestCaseResponseDecoderGeneric[account_schema.AccountUserResponse]

	accountUserRec, err := th.Data.GetAccountUserRecByRef(harness.StandardAccountRef)
	require.NoError(t, err, "GetAccountUserRecByRef returns without error")

	testCases := []testutil.TestCase{
		{
			Name: "authenticated user when get me then returns authenticated account user",
			HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
				return rnr.GetHandlerConfig()[account.GetMe]
			},
			RequestHeaders: func(d harness.Data) map[string]string {
				return testutil.AuthHeaderStandard(d)
			},
			ResponseDecoder: testCaseResponseDecoder,
			ResponseCode:    http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		t.Logf("Running test >%s<", testCase.Name)

		t.Run(testCase.Name, func(t *testing.T) {
			testFunc := func(method string, body any) {
				require.NotNil(t, body, "Response body is not nil")
				resp, ok := body.(account_schema.AccountUserResponse)
				require.True(t, ok, "Response body is of type account_schema.AccountUserResponse")
				aResp := resp.Data
				require.NotNil(t, aResp, "AccountUserResponseData is not nil")
				require.NotEmpty(t, aResp.ID, "AccountUser ID is not empty")
				require.NotEmpty(t, aResp.AccountID, "AccountUser AccountID is not empty")
				require.NotEmpty(t, aResp.Email, "AccountUser Email is not empty")
				require.Equal(t, accountUserRec.ID, aResp.ID, "AccountUser ID matches authenticated user")
				require.Equal(t, accountUserRec.Email, aResp.Email, "AccountUser Email matches authenticated user")
				require.False(t, aResp.CreatedAt.IsZero(), "AccountUser CreatedAt is not zero")
			}

			testutil.RunTestCase(t, th, &testCase, testFunc)
		})
	}
}
