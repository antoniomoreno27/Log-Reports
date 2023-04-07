package main

import (
	"time"

	"github.com/antoniomoreno27/logs-LS/internal/controllers"
	"github.com/antoniomoreno27/logs-LS/internal/domain"
	"github.com/antoniomoreno27/logs-LS/internal/services/logger"
	"github.com/antoniomoreno27/logs-LS/internal/services/report"
)

func main() {
	const (
		cookie string = `my cookie`
	)
	logger.Warnf("Configurating Report")
	timeWindowSample, err := time.ParseDuration("2h")
	if err != nil {
		logger.Panic("bat time window sample definition")
	}
	rs, err := report.New(
		cookie,
		domain.Report{
			Name:             "parse_error.csv",
			Path:             "./reports",
			Query:            "Error retrieving iterator next data",
			TimeWindowSample: timeWindowSample,
			StartDate:        time.Now().AddDate(0, 0, -9),
			FinishDate:       time.Now(),
		},
	)
	if err != nil {
		panic(err)
	}
	reportController := controllers.Report{
		ReportService: rs,
	}
	logger.Warnf("Creating report")
	reportController.Create()
}
