package domain

import "time"

type Report struct {
	Name             string
	Path             string
	Query            string
	TimeWindowSample time.Duration
	StartDate        time.Time
	FinishDate       time.Time
}
