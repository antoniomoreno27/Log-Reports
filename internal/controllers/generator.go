package controllers

import (
	"log"

	"github.com/antoniomoreno27/logs-LS/internal/services/report"
)

type Report struct {
	ReportService *report.ReportService
}

func (r *Report) Create() {
	err := r.ReportService.Create()
	if err != nil {
		log.Fatal(err)
	}
}
