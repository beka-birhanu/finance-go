/*
Package timeservice provides an interface for time-related operations.

It includes the `ITimeService` interface for obtaining the current UTC time.
*/
package timeservice

import "time"

// ITimeService defines methods for obtaining time information.
//
// Methods:
// - NowUTC() time.Time: Returns the current time in UTC.
type ITimeService interface {
	// NowUTC returns the current time in UTC.
	NowUTC() time.Time
}
