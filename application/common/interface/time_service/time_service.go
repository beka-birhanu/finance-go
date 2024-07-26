/*
Package timeservice provides an interface for time-related operations.

It includes the `IService` interface for obtaining the current UTC time.
*/
package itimeservice

import "time"

// IService defines methods for obtaining time information.
//
// Methods:
// - NowUTC() time.Time: Returns the current time in UTC.
type IService interface {
	// NowUTC returns the current time in UTC.
	NowUTC() time.Time
}
