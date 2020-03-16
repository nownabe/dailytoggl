package dailytoggl

import "time"

type togglClientMock struct{}

var yesterday time.Time

func (c *togglClientMock) getDayTotal(date time.Time) (int64, error) {
	if yesterday.Format("2006-01-02") == date.Format("2006-01-02") {
		return 100, nil
	}

	return 200, nil
}
