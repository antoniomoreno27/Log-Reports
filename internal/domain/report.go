package domain

import "time"

type Report struct {
	NeedsIndividualReport bool
	Name                  string
	Path                  string
	Query                 string
	TimeWindowSample      time.Duration
	StartDate             time.Time
	FinishDate            time.Time
}
