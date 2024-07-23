package timeservice

import "time"

type ITimeService interface {
	NowUTC() time.Time
}
