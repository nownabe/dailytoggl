package dailytoggl

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const pixelaBaseURL = "https://pixe.la/v1/users"

// PixelaClient defines pixela client's behaviors.
type PixelaClient interface {
	getGraph() ([]byte, error)
	graphURL() string
	update(time.Time, int64) error
}

type pixelaClient struct {
	*http.Client

	token    string
	username string
	graphID  string
}

func newPixelaClient(conf *config) PixelaClient {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	return &pixelaClient{
		Client:   client,
		token:    conf.PixelaToken,
		username: conf.PixelaUsername,
		graphID:  conf.PixelaGraphID,
	}
}

func (c *pixelaClient) getGraph() ([]byte, error) {
	fmt.Println(c.graphURL())
	req, err := http.NewRequest("GET", c.graphURL(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("pixela error: status = %d, message = %s",
			resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func (c *pixelaClient) graphURL() string {
	return fmt.Sprintf("%s/%s/graphs/%s", pixelaBaseURL, c.username, c.graphID)
}

func (c *pixelaClient) update(date time.Time, val int64) error {
	reqBody := fmt.Sprintf(`{"quantity":"%d"}`, val/1000)
	req, err := http.NewRequest("PUT",
		c.graphURL()+"/"+date.Format("20060102"),
		strings.NewReader(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("X-USER-TOKEN", c.token)

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("pixela error: %v", err)
		}
		return fmt.Errorf("pixela error: status = %d, message = %s",
			resp.StatusCode, string(respBody))
	}

	return nil
}
