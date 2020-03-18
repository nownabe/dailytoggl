package dailytoggl

import "time"

type pixelaClientMock struct{}

func (c *pixelaClientMock) getGraph() ([]byte, error) {
	return nil, nil
}

func (c *pixelaClientMock) graphURL() string {
	return ""
}

func (c *pixelaClientMock) update(date time.Time, val int64) error {
	return nil
}
