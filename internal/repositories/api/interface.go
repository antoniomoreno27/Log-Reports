package api

import (
	"time"

	"github.com/antoniomoreno27/logs-LS/internal/domain"
)

type LogAPIClient interface {
	Get(query string, startDate, finishDate time.Time) (logReport *domain.LogReport, err error)
}
