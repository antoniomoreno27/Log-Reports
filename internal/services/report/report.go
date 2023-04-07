package report

import (
	"bufio"
	"errors"
	"math"
	"os"
	"path"
	"regexp"

	"github.com/antoniomoreno27/logs-LS/internal/domain"
	"github.com/antoniomoreno27/logs-LS/internal/repositories/api"
	"github.com/antoniomoreno27/logs-LS/internal/services/generator"
	"github.com/antoniomoreno27/logs-LS/internal/services/logger"
	"github.com/antoniomoreno27/logs-LS/internal/services/scraper"
)

type ReportService struct {
	domain.Report

	KibanaAPIClient api.KibanaAPIClient
	Generator       generator.Generator
	Scraper         scraper.Scraper
}

func New(cookie string, report domain.Report) (*ReportService, error) {
	matchers := []scraper.Matcher{
		scraper.NewRegExpMatcher("entity_id", regexp.MustCompile("/(M[A-Z]{2})?[0-9]+/")),
		scraper.NewRegExpMatcher("entity_site", regexp.MustCompile(`site:(M[A-Z]{2})\]`)),
		scraper.NewRegExpMatcher("error_message", regexp.MustCompile("ERROR:(.*?)-")),
		scraper.NewRegExpMatcher("dump_attribute", regexp.MustCompile(`attribute:(.*?)\]`)),
		scraper.NewRegExpMatcher("dump_resource", regexp.MustCompile(`resource:(.*?)\]`)),
	}
	rs := &ReportService{
		Report:          report,
		KibanaAPIClient: *api.New(cookie),
		Scraper:         scraper.NewScraper(matchers...),
	}
	data, err := rs.loadData()
	if err != nil {
		return nil, err
	}
	rs.Generator = generator.NewCSVGenerator(";", data)
	return rs, nil
}

func (rs *ReportService) Create() error {
	logger.Warnf("Creating file")
	path := path.Join(rs.Path, rs.Name)
	os.Remove(path)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	buffer := bufio.NewWriterSize(file, 4096)
	logger.Warnf("Writing data")
	err = rs.Generator.Generate()
	if err != nil {
		return err
	}
	buffer.WriteString(rs.Generator.String())
	err = buffer.Flush()
	if err != nil {
		return err
	}
	logger.Warnf("Finished writing")
	return nil
}

func (rs *ReportService) loadData() (chan []string, error) {
	var (
		numBatches               = int(math.Ceil(rs.FinishDate.Sub(rs.StartDate).Seconds() / rs.TimeWindowSample.Seconds()))
		hits       []domain.Hit  = make([]domain.Hit, 0, 3000*numBatches)
		data       chan []string = make(chan []string, 1)
		err        error

		i       = -1
		counter = 0
	)
	logger.Warnf("Starting to load data")
	actualStart := rs.StartDate
	for i < 0 {
		counter++
		logger.Warnf("Loading batch [%d/%d]", counter, numBatches)
		lr, err := rs.KibanaAPIClient.Get(rs.Query, actualStart, actualStart.Add(rs.TimeWindowSample))
		if err != nil {
			logger.Errorf("error while getting data from kibana >> %v", err)
		}
		newStart := actualStart.Add(rs.TimeWindowSample)
		if len(lr.Log.Hits) == 0 {
			logger.Errorf("no data found between >> %s - %s", actualStart.String(), newStart.String())
		}
		actualStart = newStart
		i = actualStart.Compare(rs.FinishDate)
		hits = append(hits, lr.Log.Hits...)
	}
	go func() {
		defer close(data)
		for _, hit := range hits {
			matchs, matchErr := rs.Scraper.Scrape(*hit.Source.Message)
			if err != nil {
				errors.Join(err, matchErr)
				continue
			}
			data <- matchs
		}
	}()
	return data, err
}
