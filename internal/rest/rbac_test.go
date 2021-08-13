package rest_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"rbac/internal/rest"
	"rbac/internal/rest/resttesting"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/gorilla/mux"
)

func TestTasks_Post(t *testing.T) {
	// XXX: Test "serviceArgs"

	t.Parallel()

	type output struct {
		expectedStatus int
		expected       interface{}
		target         interface{}
	}

	tests := []struct {
		name   string
		setup  func(*resttesting.FakeRBACService)
		input  []byte
		output output
	}{
		{
			"OK: 201",
			func(s *resttesting.FakeRBACService) {
				s.CreateAccountReturns(nil)
			},
			func() []byte {
				b, _ := json.Marshal(&rest.RegisterRequest{
					Username:          "test",
					Password:          "test",
					ProfilePicure:     "test",
					ProfileBackground: "test",
					Firstname:         "test",
					Lastname:          "test",
					Mobile:            "09123456789",
					Email:             "test@test.com",
				})

				return b
			}(),
			output{
				http.StatusCreated,
				&rest.AccountResponse{
					Message: "Created Successfully",
				},
				&rest.AccountResponse{},
			},
		},
		// {
		// 	"ERR: 400",
		// 	func(*resttesting.FakeAccountService) {},
		// 	[]byte(`{"invalid":"json`),
		// 	output{
		// 		http.StatusBadRequest,
		// 		&rest.ErrorResponse{
		// 			Error: "invalid request",
		// 		},
		// 		&rest.ErrorResponse{},
		// 	},
		// },
		// {
		// 	"ERR: 500",
		// 	func(s *resttesting.FakeAccountService) {
		// 		s.CreateReturns(internal.Task{},
		// 			errors.New("service error"))
		// 	},
		// 	[]byte(`{}`),
		// 	output{
		// 		http.StatusInternalServerError,
		// 		&rest.ErrorResponse{
		// 			Error: "create failed",
		// 		},
		// 		&rest.ErrorResponse{},
		// 	},
		// },
	}

	//-

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			router := mux.NewRouter()
			svc := &resttesting.FakeRBACService{}
			tt.setup(svc)

			rest.NewRBACHandler(svc).Register(router)

			//-

			res := doRequest(router,
				httptest.NewRequest(http.MethodPost, "/accounts/register", bytes.NewReader(tt.input)))

			//-

			assertResponse(t, res, test{tt.output.expected, tt.output.target})

			if tt.output.expectedStatus != res.StatusCode {
				t.Fatalf("expected code %d, actual %d", tt.output.expectedStatus, res.StatusCode)
			}
		})
	}
}

type test struct {
	expected interface{}
	target   interface{}
}

func doRequest(router *mux.Router, req *http.Request) *http.Response {
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	return rr.Result()
}

func assertResponse(t *testing.T, res *http.Response, test test) {
	t.Helper()

	if err := json.NewDecoder(res.Body).Decode(test.target); err != nil {
		t.Fatalf("couldn't decode %s", err)
	}
	defer res.Body.Close()

	if !cmp.Equal(test.expected, test.target) {
		t.Fatalf("expected results don't match: %s", cmp.Diff(test.expected, test.target))
	}
}
