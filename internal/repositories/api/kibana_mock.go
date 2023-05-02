package api

import (
	"time"

	"github.com/antoniomoreno27/logs-LS/internal/domain"
)

type kibanaAPIClientMock struct {
	GetMock func(query string, startDate, finishDate time.Time) (logReport *domain.LogReport, err error)
}

func (kM *kibanaAPIClientMock) Get(query string, startDate, finishDate time.Time) (logReport *domain.LogReport, err error) {
	if kM.GetMock != nil {
		return kM.GetMock(query, startDate, finishDate)
	}
	return
}
