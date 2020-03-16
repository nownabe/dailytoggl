package dailytoggl

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
)

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
			name:     "valid request",
			method:   "POST",
			request:  fmt.Sprintf(`{"auth_token":"%s"}`, conf.AuthToken),
			status:   200,
			response: "ok",
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
