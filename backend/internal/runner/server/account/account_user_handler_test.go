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

func Test_getAccountUserHandler(t *testing.T) {
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

	testCaseCollectionResponseDecoder := testutil.TestCaseResponseDecoderGeneric[account_schema.AccountUserCollectionResponse]
	testCaseResponseDecoder := testutil.TestCaseResponseDecoderGeneric[account_schema.AccountUserResponse]

	// accountUserRec IS the AccountUser record returned by GetAccountUserRecByRef
	accountUserRec, err := th.Data.GetAccountUserRecByRef(harness.StandardAccountRef)
	require.NoError(t, err, "GetAccountUserRecByRef returns without error")

	testCases := []testCase{
		{
			TestCase: testutil.TestCase{
				Name: "authenticated user when get many account users then returns expected account users",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[account.GetManyAccountUsers]
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
				Name: "authenticated user when get many account users by account then returns expected account users",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()["get-many-account-users-by-account"]
				},
				RequestHeaders: func(d harness.Data) map[string]string {
					return testutil.AuthHeaderStandard(d)
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":account_id": accountUserRec.AccountID,
					}
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
				Name: "authenticated user when get one account user with valid IDs then returns expected account user",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[account.GetOneAccountUser]
				},
				RequestHeaders: func(d harness.Data) map[string]string {
					return testutil.AuthHeaderStandard(d)
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":account_id":      accountUserRec.AccountID,
						":account_user_id": accountUserRec.ID,
					}
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

				var responses []*account_schema.AccountUserResponseData
				if testCase.collectionRequest {
					responses = body.(account_schema.AccountUserCollectionResponse).Data
				} else {
					responses = append(responses, body.(account_schema.AccountUserResponse).Data)
				}

				if testCase.collectionRequest {
					require.True(t, len(responses) >= testCase.collectionRecordCount, "Response record count length >= expected")
				}

				if testCase.collectionRequest && testCase.collectionRecordCount == 0 {
					require.Empty(t, responses, "Response body should be empty")
				} else {
					require.NotEmpty(t, responses, "Response body is not empty")
				}

				for _, d := range responses {
					require.NotEmpty(t, d.ID, "Account user ID is not empty")
					require.NotEmpty(t, d.AccountID, "Account user AccountID is not empty")
					require.NotEmpty(t, d.Email, "Account user Email is not empty")
					require.NotEmpty(t, d.Status, "Account user Status is not empty")
				}
			}

			testutil.RunTestCase(t, th, &testCase, testFunc)
		})
	}
}

func Test_createUpdateDeleteAccountUserHandler(t *testing.T) {
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
		expectResponse func(d harness.Data, req account_schema.AccountUserRequest) account_schema.AccountUserResponse
	}

	testCaseResponseDecoder := testutil.TestCaseResponseDecoderGeneric[account_schema.AccountUserResponse]

	// accountUserRec IS the AccountUser; .AccountID is the parent account_record.Account ID
	accountUserRec, err := th.Data.GetAccountUserRecByRef(harness.StandardAccountRef)
	require.NoError(t, err, "GetAccountUserRecByRef returns without error")

	testCases := []testCase{
		{
			TestCase: testutil.TestCase{
				Name: "authenticated user when create account user with valid properties then returns created account user",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[account.CreateOneAccountUser]
				},
				RequestHeaders: func(d harness.Data) map[string]string {
					return testutil.AuthHeaderStandard(d)
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":account_id": accountUserRec.AccountID,
					}
				},
				RequestBody: func(d harness.Data) any {
					e := gofakeit.Email()
					s := "pending_approval"
					return account_schema.AccountUserRequest{
						Email:  &e,
						Status: &s,
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusCreated,
			},
			expectResponse: func(d harness.Data, req account_schema.AccountUserRequest) account_schema.AccountUserResponse {
				return account_schema.AccountUserResponse{
					Data: &account_schema.AccountUserResponseData{
						AccountID: accountUserRec.AccountID,
						Email:     *req.Email,
						Status:    *req.Status,
					},
				}
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "authenticated user when update account user with valid properties then returns updated account user",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[account.UpdateOneAccountUser]
				},
				RequestHeaders: func(d harness.Data) map[string]string {
					return testutil.AuthHeaderStandard(d)
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":account_id":      accountUserRec.AccountID,
						":account_user_id": accountUserRec.ID,
					}
				},
				RequestBody: func(d harness.Data) any {
					s := "disabled"
					return account_schema.AccountUserRequest{
						Status: &s,
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
			expectResponse: func(d harness.Data, req account_schema.AccountUserRequest) account_schema.AccountUserResponse {
				return account_schema.AccountUserResponse{
					Data: &account_schema.AccountUserResponseData{
						AccountID: accountUserRec.AccountID,
						Email:     accountUserRec.Email,
						Status:    *req.Status,
					},
				}
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "authenticated user when delete account user with valid IDs then returns no content",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[account.DeleteOneAccountUser]
				},
				RequestHeaders: func(d harness.Data) map[string]string {
					return testutil.AuthHeaderStandard(d)
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":account_id":      accountUserRec.AccountID,
						":account_user_id": accountUserRec.ID,
					}
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
				resp, ok := body.(account_schema.AccountUserResponse)
				require.True(t, ok, "Response body is of type account_schema.AccountUserResponse")
				auResp := resp.Data
				require.NotNil(t, auResp, "AccountUserResponseData is not nil")
				t.Logf("AccountUserResponseData: %#v", auResp)
				require.NotEmpty(t, auResp.ID, "Account user ID is not empty")
				require.NotEmpty(t, auResp.AccountID, "Account user AccountID is not empty")
				require.NotEmpty(t, auResp.Email, "Account user Email is not empty")
				require.NotEmpty(t, auResp.Status, "Account user Status is not empty")

				if testCase.expectResponse != nil {
					xResp := testCase.expectResponse(
						th.Data,
						testCase.TestRequestBody(th.Data).(account_schema.AccountUserRequest),
					).Data
					require.Equal(t, xResp.AccountID, auResp.AccountID, "Account user AccountID matches expected")

					if xResp.Email != "" {
						require.Equal(t, xResp.Email, auResp.Email, "Account user Email equals expected")
					}
					if xResp.Status != "" {
						require.Equal(t, xResp.Status, auResp.Status, "Account user Status equals expected")
					}
				}

				require.False(t, auResp.CreatedAt.IsZero(), "Account user CreatedAt is not zero")
			}

			testutil.RunTestCase(t, th, &testCase, testFunc)
		})
	}
}
