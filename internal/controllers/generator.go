package controllers

import (
	"github.com/antoniomoreno27/logs-LS/internal/services/logger"
	"github.com/antoniomoreno27/logs-LS/internal/services/report"
)

type Report struct {
	ReportService *report.ReportService
}

func (r *Report) Create() {
	err := r.ReportService.Create()
	if err != nil {
		logger.Errorf("%v", err)
	}
}
