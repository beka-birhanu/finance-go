package timeservice

import (
	"time"

	timeservice "github.com/beka-birhanu/finance-go/application/common/interface/time_service"
)

type TimeService struct{}

var _ timeservice.ITimeService = &TimeService{}

func NewTimeService() *TimeService {
	return &TimeService{}
}

func (t *TimeService) NowUTC() time.Time {
	return time.Now().UTC()
}
