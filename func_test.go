package dailytoggl

import (
	"fmt"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	yesterday = time.Now().Add(-24 * time.Hour).In(loc)
	toggl = &togglClientMock{}
	pixela = &pixelaClientMock{}

	code := m.Run()
	os.Exit(code)
}

func TestDailyToggl(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		request  string
		status   int
		response string
	}{
		{
			name:     "unsupported method",
			method:   "GET",
			request:  ``,
			status:   400,
			response: "unsupported http method\n",
		},
		{
			name:     "invalid json",
			method:   "POST",
			request:  ``,
			status:   400,
			response: "failed to parse request\n",
		},
		{
			name:     "wrong auth_token",
			method:   "POST",
			request:  `{"auth_token":"wrong token"}`,
			status:   401,
			response: "authentication error\n",
		},
		{
			name:     "without date parameter",
			method:   "POST",
			request:  fmt.Sprintf(`{"auth_token":"%s"}`, conf.AuthToken),
			status:   200,
			response: "100",
		},
		{
			name:     "with date parameter",
			method:   "POST",
			request:  fmt.Sprintf(`{"auth_token":"%s","date":"2020-01-01"}`, conf.AuthToken),
			status:   200,
			response: "200",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/", strings.NewReader(tt.request))
			req.Header.Add("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			DailyToggl(rr, req)

			if got := rr.Result().StatusCode; got != tt.status {
				t.Errorf("Response status code = %d, expected %d", got, tt.status)
			}

			if got := rr.Body.String(); got != tt.response {
				t.Errorf("Response body = '%s', expected '%s'", got, tt.response)
			}
		})
	}
}
