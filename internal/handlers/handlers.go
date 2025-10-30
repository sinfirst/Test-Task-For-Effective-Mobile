package handlers

import (
	"fmt"
	"time"
)

func DateParse(start, end string) (time.Time, time.Time, error) {
	var startTime, endTime time.Time
	startTime, err := time.Parse("01-2006", start)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid start_month format: %w", err)
	}

	if end != "" {
		endTime, err = time.Parse("01-2006", end)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid end_month format: %w", err)
		}
		return startTime, endTime, nil
	}
	return startTime, time.Time{}, nil
}
