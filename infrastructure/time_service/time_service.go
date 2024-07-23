package timeservice

import "time"

type TimeService struct{}

func NewTimeService() *TimeService {
	return &TimeService{}
}
func (t *TimeService) NowUTC() time.Time {
	return time.Now().UTC()
}
