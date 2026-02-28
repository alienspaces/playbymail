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

func Test_getAccountUserContactHandler(t *testing.T) {
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

	testCaseCollectionResponseDecoder := testutil.TestCaseResponseDecoderGeneric[account_schema.AccountContactCollectionResponse]
	testCaseResponseDecoder := testutil.TestCaseResponseDecoderGeneric[account_schema.AccountContactResponse]

	// Setup: get an account for reference
	accountRec, err := th.Data.GetAccountUserRecByRef(harness.StandardAccountRef)
	require.NoError(t, err, "GetAccountUserRecByRef returns without error")

	// Get account contact (created by harness for all accounts)
	accountUserContactRec, err := th.Data.GetAccountUserContactRecByAccountUserID(accountRec.ID)
	require.NoError(t, err, "GetAccountUserContactRecByAccountUserID returns without error")

	testCases := []testCase{
		{
			TestCase: testutil.TestCase{
				Name: "authenticated user when get many account contacts by user path then returns expected account contacts",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()["get-many-account-user-contacts-by-user"]
				},
				RequestHeaders: func(d harness.Data) map[string]string {
					return testutil.AuthHeaderStandard(d)
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":account_id":      accountRec.AccountID,
						":account_user_id": accountRec.ID,
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
				Name: "authenticated user when get many account contacts without user path then returns only own contacts",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[account.GetManyAccountUserContacts]
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
				Name: "authenticated user when get one account contact with valid account contact ID then returns expected account contact",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[account.GetOneAccountUserContact]
				},
				RequestHeaders: func(d harness.Data) map[string]string {
					return testutil.AuthHeaderStandard(d)
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					params := map[string]string{
						":account_id":              accountRec.AccountID,
						":account_user_id":         accountRec.ID,
						":account_user_contact_id": accountUserContactRec.ID,
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

				var responses []*account_schema.AccountContactResponseData
				if testCase.collectionRequest {
					responses = body.(account_schema.AccountContactCollectionResponse).Data
				} else {
					responses = append(responses, body.(account_schema.AccountContactResponse).Data)
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
					require.NotEmpty(t, d.ID, "Account contact ID is not empty")
					require.NotEmpty(t, d.AccountUserID, "Account contact AccountUserID is not empty")
					require.NotEmpty(t, d.Name, "Account contact Name is not empty")
					require.NotEmpty(t, d.PostalAddressLine1, "Account contact PostalAddressLine1 is not empty")
					require.NotEmpty(t, d.StateProvince, "Account contact StateProvince is not empty")
					require.NotEmpty(t, d.Country, "Account contact Country is not empty")
					require.NotEmpty(t, d.PostalCode, "Account contact PostalCode is not empty")
				}
			}

			testutil.RunTestCase(t, th, &testCase, testFunc)
		})
	}
}

func Test_createUpdateDeleteAccountUserContactHandler(t *testing.T) {
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
		expectResponse func(d harness.Data, req account_schema.AccountContactRequest) account_schema.AccountContactResponse
	}

	testCaseResponseDecoder := testutil.TestCaseResponseDecoderGeneric[account_schema.AccountContactResponse]

	// Setup: get an account for reference
	accountRec, err := th.Data.GetAccountUserRecByRef(harness.StandardAccountRef)
	require.NoError(t, err, "GetAccountUserRecByRef returns without error")

	// Get account contact (created by harness for all accounts)
	accountUserContactRec, err := th.Data.GetAccountUserContactRecByAccountUserID(accountRec.ID)
	require.NoError(t, err, "GetAccountUserContactRecByAccountUserID returns without error")

	testCases := []testCase{
		{
			TestCase: testutil.TestCase{
				Name: "authenticated user when create account contact with valid properties then returns created account contact",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					cfg := rnr.GetHandlerConfig()[account.CreateOneAccountUserContact]
					t.Logf("DEBUG CONFIG PATH: %s", cfg.Path)
					return cfg
				},
				RequestHeaders: func(d harness.Data) map[string]string {
					return testutil.AuthHeaderStandard(d)
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					return map[string]string{
						":account_id":      accountRec.AccountID,
						":account_user_id": accountRec.ID,
					}
				},
				RequestBody: func(d harness.Data) any {
					return account_schema.AccountContactRequest{
						Name:               gofakeit.Name(),
						PostalAddressLine1: gofakeit.Address().Address,
						PostalAddressLine2: gofakeit.Address().Street,
						StateProvince:      gofakeit.Address().State,
						Country:            gofakeit.Address().Country,
						PostalCode:         gofakeit.Address().Zip,
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusCreated,
			},
			expectResponse: func(d harness.Data, req account_schema.AccountContactRequest) account_schema.AccountContactResponse {
				return account_schema.AccountContactResponse{
					Data: &account_schema.AccountContactResponseData{
						AccountUserID:      accountRec.ID,
						Name:               req.Name,
						PostalAddressLine1: req.PostalAddressLine1,
						PostalAddressLine2: req.PostalAddressLine2,
						StateProvince:      req.StateProvince,
						Country:            req.Country,
						PostalCode:         req.PostalCode,
					},
				}
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "authenticated user when update account contact with valid properties then returns updated account contact",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[account.UpdateOneAccountUserContact]
				},
				RequestHeaders: func(d harness.Data) map[string]string {
					return testutil.AuthHeaderStandard(d)
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					params := map[string]string{
						":account_id":              accountRec.AccountID,
						":account_user_id":         accountRec.ID,
						":account_user_contact_id": accountUserContactRec.ID,
					}
					return params
				},
				RequestBody: func(d harness.Data) any {
					return account_schema.AccountContactRequest{
						Name:               gofakeit.Name(),
						PostalAddressLine1: gofakeit.Address().Address,
						PostalAddressLine2: gofakeit.Address().Street,
						StateProvince:      gofakeit.Address().State,
						Country:            gofakeit.Address().Country,
						PostalCode:         gofakeit.Address().Zip,
					}
				},
				ResponseDecoder: testCaseResponseDecoder,
				ResponseCode:    http.StatusOK,
			},
			expectResponse: func(d harness.Data, req account_schema.AccountContactRequest) account_schema.AccountContactResponse {
				return account_schema.AccountContactResponse{
					Data: &account_schema.AccountContactResponseData{
						AccountUserID:      accountRec.ID,
						Name:               req.Name,
						PostalAddressLine1: req.PostalAddressLine1,
						PostalAddressLine2: req.PostalAddressLine2,
						StateProvince:      req.StateProvince,
						Country:            req.Country,
						PostalCode:         req.PostalCode,
					},
				}
			},
		},
		{
			TestCase: testutil.TestCase{
				Name: "authenticated user when delete account contact with valid account contact ID then returns no content",
				HandlerConfig: func(rnr testutil.TestRunnerer) server.HandlerConfig {
					return rnr.GetHandlerConfig()[account.DeleteOneAccountUserContact]
				},
				RequestHeaders: func(d harness.Data) map[string]string {
					return testutil.AuthHeaderStandard(d)
				},
				RequestPathParams: func(d harness.Data) map[string]string {
					params := map[string]string{
						":account_id":              accountRec.AccountID,
						":account_user_id":         accountRec.ID,
						":account_user_contact_id": accountUserContactRec.ID,
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
					// No content expected
					return
				}

				require.NotNil(t, body, "Response body is not nil")
				t.Logf("Actual response body: %#v", body)
				resp, ok := body.(account_schema.AccountContactResponse)
				require.True(t, ok, "Response body is of type account_schema.AccountContactResponse")
				acResp := resp.Data
				require.NotNil(t, acResp, "AccountContactResponseData is not nil")
				t.Logf("AccountContactResponseData: %#v", acResp)
				require.NotEmpty(t, acResp.ID, "Account contact ID is not empty")
				require.NotEmpty(t, acResp.AccountUserID, "Account contact AccountUserID is not empty")
				require.NotEmpty(t, acResp.Name, "Account contact Name is not empty")
				require.NotEmpty(t, acResp.PostalAddressLine1, "Account contact PostalAddressLine1 is not empty")
				require.NotEmpty(t, acResp.StateProvince, "Account contact StateProvince is not empty")
				require.NotEmpty(t, acResp.Country, "Account contact Country is not empty")
				require.NotEmpty(t, acResp.PostalCode, "Account contact PostalCode is not empty")

				if testCase.expectResponse != nil {
					xResp := testCase.expectResponse(
						th.Data,
						testCase.TestRequestBody(th.Data).(account_schema.AccountContactRequest),
					).Data
					require.Equal(t, xResp.AccountUserID, acResp.AccountUserID, "Account contact AccountUserID matches expected")
					require.Equal(t, xResp.Name, acResp.Name, "Account contact Name equals expected")
					require.Equal(t, xResp.PostalAddressLine1, acResp.PostalAddressLine1, "Account contact PostalAddressLine1 equals expected")
					require.Equal(t, xResp.PostalAddressLine2, acResp.PostalAddressLine2, "Account contact PostalAddressLine2 equals expected")
					require.Equal(t, xResp.StateProvince, acResp.StateProvince, "Account contact StateProvince equals expected")
					require.Equal(t, xResp.Country, acResp.Country, "Account contact Country equals expected")
					require.Equal(t, xResp.PostalCode, acResp.PostalCode, "Account contact PostalCode equals expected")
				}

				require.False(t, acResp.CreatedAt.IsZero(), "Account contact CreatedAt is not zero")
			}

			testutil.RunTestCase(t, th, &testCase, testFunc)
		})
	}
}
