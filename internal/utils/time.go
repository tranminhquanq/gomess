package utils

import "time"

// CurrentTimestamp returns the current timestamp in the specified unit (seconds, milliseconds, or nanoseconds).
// If an invalid unit is provided, it defaults to seconds
func CurrentTimestamp(unit string) int64 {
	now := time.Now()
	switch unit {
	case "seconds":
		return now.Unix()
	case "milliseconds":
		return now.UnixMilli()
	case "nanoseconds":
		return now.UnixNano()
	default:
		// Default to seconds if an invalid unit is provided
		return now.Unix()
	}
}
