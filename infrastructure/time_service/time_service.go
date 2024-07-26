package timeservice

import (
	"time"

	timeservice "github.com/beka-birhanu/finance-go/application/common/interface/time_service"
)

// Service provides the current UTC time.
type Service struct{}

// Ensure Service implements timeservice.IService.
var _ timeservice.IService = &Service{}

// New creates a new instance of the Service.
func New() *Service {
	return &Service{}
}

// NowUTC returns the current time in UTC.
func (t *Service) NowUTC() time.Time {
	return time.Now().UTC()
}

