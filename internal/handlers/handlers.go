package handlers

import (
	"fmt"
	"time"
)

func DateParse(start, end string) (time.Time, time.Time, error) {
	startTime, err := time.Parse("mm-YYYY", start)
	if err != nil {
		return time.Now(), time.Now(), fmt.Errorf("invalid start_month format: %w", err)
	}
	endTime, err := time.Parse("mm-YYYY", end)
	if err != nil {
		return time.Now(), time.Now(), fmt.Errorf("invalid end_month format: %w", err)
	}
	return startTime, endTime, nil
}
