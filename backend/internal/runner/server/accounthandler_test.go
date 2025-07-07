package runner

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/schema"
)

func Test_getAccountHandler(t *testing.T) {
	t.Parallel()

	th := newTestHarness(t)
	require.NotNil(t, th, "newTestHarness returns without error")

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	type testCase struct {
		TestCase
		collectionRequest     bool
		collectionRecordCount int
	}

	testCaseCollectionResponseDecoder := testCaseResponseDecoderGeneric[schema.AccountCollectionResponse]
	testCaseResponseDecoder := testCaseResponseDecoderGeneric[schema.AccountResponse]

	testCases := []testCase{
		{
			TestCase: TestCase{
				Name: "API key with open access \\ get many accounts \\ returns expected accounts",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[getManyAccounts]
				},
				RequestQueryParams: func(d harness.Data) map[string]any {
					accountRec, err := d.GetAccountRecByRef(harness.AccountOneRef)
					require.NoError(t, err, "GetAccountRecByRef returns without error")
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
			TestCase: TestCase{
				Name: "API key with open access \\ get many accounts with id param \\ returns expected account",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[getManyAccounts]
				},
				RequestQueryParams: func(d harness.Data) map[string]any {
					accountRec, err := d.GetAccountRecByRef(harness.AccountOneRef)
					require.NoError(t, err, "GetAccountRecByRef returns without error")
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
			TestCase: TestCase{
				Name: "API key with open access \\ get one account with valid account ID \\ returns expected account",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[getOneAccount]
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

		t.Logf("Running test >%s<", testCase.Name)

		t.Run(testCase.Name, func(t *testing.T) {
			testFunc := func(method string, body interface{}) {
				if testCase.TestResponseCode() != http.StatusOK {
					return
				}

				require.NotNil(t, body, "Response body is not nil")

				responses := []*schema.AccountResponseData{}
				if testCase.collectionRequest {
					responses = body.(schema.AccountCollectionResponse)
				} else {
					responses = append(responses, body.(schema.AccountResponse).AccountResponseData)
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

			RunTestCase(t, th, &testCase, testFunc)
		})
	}
}

func Test_createUpdateDeleteAccountHandler(t *testing.T) {
	t.Parallel()

	th := newTestHarness(t)
	require.NotNil(t, th, "newTestHarness returns without error")

	_, err := th.Setup()
	require.NoError(t, err, "Test data setup returns without error")
	defer func() {
		err = th.Teardown()
		require.NoError(t, err, "Test data teardown returns without error")
	}()

	type testCase struct {
		TestCase
		expectResponse func(d harness.Data, req schema.AccountRequest) schema.AccountResponse
	}

	testCaseResponseDecoder := testCaseResponseDecoderGeneric[schema.AccountResponse]

	testCases := []testCase{
		{
			TestCase: TestCase{
				Name: "API key with open access \\ create account with valid properties \\ returns created account",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[createAccount]
				},
				RequestBody: func(d harness.Data) interface{} {
					return schema.AccountRequest{
						Email: "test@example.com",
						Name:  "Test User",
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusCreated,
			},
			expectResponse: func(d harness.Data, req schema.AccountRequest) schema.AccountResponse {
				return schema.AccountResponse{
					AccountResponseData: &schema.AccountResponseData{
						Email: req.Email,
						Name:  req.Name,
					},
				}
			},
		},
		{
			TestCase: TestCase{
				Name: "API key with open access \\ update account with valid properties \\ returns updated account",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[updateAccount]
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					accountRec, err := d.GetAccountRecByRef(harness.AccountOneRef)
					require.NoError(t, err, "GetAccountRecByRef returns without error")
					params := map[string]string{
						":account_id": accountRec.ID,
					}
					return params
				},
				RequestBody: func(d harness.Data) interface{} {
					return schema.AccountRequest{
						Email: "updated@example.com",
						Name:  "Updated User",
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
			expectResponse: func(d harness.Data, req schema.AccountRequest) schema.AccountResponse {
				return schema.AccountResponse{
					AccountResponseData: &schema.AccountResponseData{
						Email: req.Email,
						Name:  req.Name,
					},
				}
			},
		},
		{
			TestCase: TestCase{
				Name: "API key with open access \\ delete account with valid account ID \\ returns no content",
				HandlerConfig: func(rnr *Runner) server.HandlerConfig {
					return rnr.HandlerConfig[deleteAccount]
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
				resp, ok := body.(schema.AccountResponse)
				require.True(t, ok, "Response body is of type schema.AccountResponse")
				aResp := resp.AccountResponseData
				require.NotNil(t, aResp, "AccountResponseData is not nil")
				t.Logf("AccountResponseData: %#v", aResp)
				require.NotEmpty(t, aResp.ID, "Account ID is not empty")
				require.NotEmpty(t, aResp.Email, "Account Email is not empty")
				xResp := testCase.expectResponse(
					th.Data,
					testCase.TestRequestBody(th.Data).(schema.AccountRequest),
				).AccountResponseData
				require.Equal(t, xResp.Email, aResp.Email, "Account Email matches expected")
				require.Equal(t, xResp.Name, aResp.Name, "Account Name equals expected")
				require.False(t, aResp.CreatedAt.IsZero(), "Account CreatedAt is not zero")
				// UpdatedAt is allowed to be nil, so do not assert on it
			}

			RunTestCase(t, th, &testCase, testFunc)
		})
	}
}
