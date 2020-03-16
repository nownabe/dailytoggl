package dailytoggl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const togglURL = "https://toggl.com/reports/api/v2/details"

// TogglClient defines toggl client's behaviors.
type TogglClient interface {
	getDayTotal(time.Time) (int64, error)
}

type togglClient struct {
	*http.Client
	token       string
	projectID   string
	workspaceID string
}

type togglDetailedResponse struct {
	TotalGrand int64 `json:"total_grand"`
}

func newToggl(conf *config) TogglClient {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	return &togglClient{
		Client:      client,
		token:       conf.TogglAPIToken,
		projectID:   conf.TogglProjectID,
		workspaceID: conf.TogglWorkspaceID,
	}
}

// getDayTotal returns milliseconds
func (c *togglClient) getDayTotal(date time.Time) (int64, error) {
	v := url.Values{}
	v.Set("user_agent", "dailytoggl")
	v.Set("workspace_id", c.workspaceID)
	v.Set("project_ids", c.projectID)
	v.Set("since", date.Format("2006-01-02"))
	v.Set("until", date.Format("2006-01-02"))
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", togglURL, v.Encode()), nil)
	if err != nil {
		return 0, err
	}
	req.SetBasicAuth(c.token, "api_token")

	resp, err := c.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	rawBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	body := togglDetailedResponse{}
	if err := json.Unmarshal(rawBody, &body); err != nil {
		return 0, err
	}

	return body.TotalGrand, nil
}
