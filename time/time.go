// Package time provides utility functions for working with time and dates.
// It includes functions for combining/separating date and time strings.
package time

import (
	"fmt"
	"time"
)

// Combine combines a date string and a time string into a single datetime string.
// It expects the date string in "YYYY-MM-DD" format and the time string in "HH:MM:SS" format.
// The function validates the combined string by attempting to parse it against the standard
// "YYYY-MM-DD HH:MM:SS" format. If parsing is successful, it returns the combined datetime string.
// Returns an error if the combined string cannot be parsed into a valid time.Time object.
func Combine(date, tm string) (string, error) {
	dt := fmt.Sprintf("%s %s", date, tm)
	// Validate the combined string by attempting to parse it.
	if _, err := time.Parse(time.DateTime, dt); err != nil {
		return "", fmt.Errorf("failed to parse combined datetime string \"%s\": %w", dt, err)
	}
	return dt, nil
}

// Separate separates a datetime string into its date and time components.
// It expects the datetime string to be in the "YYYY-MM-DD HH:MM:SS" format.
// Returns the date string (formatted as "YYYY-MM-DD"), the time string (formatted as "HH:MM:SS"),
// or an error if the input datetime string cannot be parsed.
func Separate(dt string) (string, string, error) {
	datetime, err := time.Parse(time.DateTime, dt)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse datetime string \"%s\": %w", dt, err)
	}

	return datetime.Format(time.DateOnly), datetime.Format(time.TimeOnly), nil
}
