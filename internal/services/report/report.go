package report

import (
	"bufio"
	"errors"
	"os"
	"path"

	"github.com/antoniomoreno27/logs-LS/internal/domain"
	"github.com/antoniomoreno27/logs-LS/internal/repositories/api"
	"github.com/antoniomoreno27/logs-LS/internal/services/generator"
	"github.com/antoniomoreno27/logs-LS/internal/services/scraper"
)

type ReportService struct {
	domain.Report

	KibanaAPIClient api.KibanaAPIClient
	Generator       generator.Generator
	Scraper         scraper.Scraper
}

func New(cookie string, report domain.Report) (*ReportService, error) {
	rs := &ReportService{
		Report:          report,
		KibanaAPIClient: *api.New(cookie),
		Scraper:         scraper.NewParserErrorScrapper(),
	}
	data, err := rs.loadData()
	if err != nil {
		return nil, err
	}
	rs.Generator = generator.NewCSVGenerator(";", data)
	return rs, nil
}

func (rs *ReportService) Create() error {
	path := path.Join(rs.Path, rs.Name)
	os.Remove(path)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	buffer := bufio.NewWriterSize(file, 4096)
	err = rs.Generator.Generate()
	if err != nil {
		return err
	}
	buffer.WriteString(rs.Generator.String())
	err = buffer.Flush()
	if err != nil {
		return err
	}
	return nil
}

func (rs *ReportService) loadData() (chan []string, error) {
	var (
		hits []domain.Hit  = make([]domain.Hit, 0, 3000*(rs.StartDate.Second()-rs.FinishDate.Second())/int(rs.TimeWindowSample))
		data chan []string = make(chan []string, 1)
		err  error
	)
	actualStart := rs.StartDate
	i := -1
	for i < 0 {
		lr, err := rs.KibanaAPIClient.Get(rs.Query, actualStart, actualStart.Add(rs.TimeWindowSample))
		actualStart = actualStart.Add(rs.TimeWindowSample)
		i = actualStart.Compare(rs.FinishDate)
		if err != nil {
			continue
		}
		if len(lr.Log.Hits) == 0 {
			continue
		}
		hits = append(hits, lr.Log.Hits...)
	}
	go func() {
		defer close(data)
		for _, hit := range hits {
			matchs, matchErr := rs.Scraper.Match(*hit.Source.Message)
			if err != nil {
				errors.Join(err, matchErr)
				continue
			}
			data <- matchs
		}
	}()
	return data, err
}
