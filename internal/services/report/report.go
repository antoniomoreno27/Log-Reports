package report

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"sync"

	"github.com/antoniomoreno27/logs-LS/internal/domain"
	"github.com/antoniomoreno27/logs-LS/internal/repositories/api"
	"github.com/antoniomoreno27/logs-LS/internal/services/generator"
	"github.com/antoniomoreno27/logs-LS/internal/services/logger"
	"github.com/antoniomoreno27/logs-LS/internal/services/scraper"
)

type ReportService struct {
	wg *sync.WaitGroup
	domain.Report
	LogAPIClient api.LogAPIClient
	Generators   []generator.Generator
	Scraper      scraper.Scraper
}

const (
	resourcePos = iota + 1
	attributePos

	fullDataKey = "full"
)

var (
	ErrNotEnoughColumns error = errors.New("not enough columns to build a key in row")
)

func New(cookie string, report domain.Report) (*ReportService, error) {
	var wg sync.WaitGroup
	matchers := []scraper.Matcher{
		scraper.NewRegExpMatcher("entity_id", regexp.MustCompile("/(M[A-Z]{2})?[0-9]+/")),
		scraper.NewRegExpMatcher("entity_site", regexp.MustCompile(`site:(M[A-Z]{2})\]|/(M[A-Z]{2})[0-9]`)),
		scraper.NewRegExpMatcher("error_message", regexp.MustCompile(`ERROR:(.*?-|.*?\:\s|.*?\")`)),
		scraper.NewRegExpMatcher("dump_attribute", regexp.MustCompile(`attribute:(.*?)\]`)),
		scraper.NewRegExpMatcher("dump_resource", regexp.MustCompile(`resource:(.*?)\]`)),
	}
	rs := &ReportService{
		Report:       report,
		LogAPIClient: api.NewLogAPIClient(cookie, false),
		Scraper:      scraper.NewScraper(matchers...),
		wg:           &wg,
	}
	data, dumpNames, err := rs.loadData()
	if err != nil {
		return nil, err
	}
	dataMap := sortOut(data, dumpNames...)
	for key, dataCh := range dataMap {
		wg.Add(1)
		rs.Generators = append(rs.Generators, generator.NewCSVGenerator(key, rs.Path, ";", &wg, dataCh))
	}
	return rs, nil
}

func (rs *ReportService) Create() error {
	logger.Warnf("Creating report")
	var err error
	for _, gen := range rs.Generators {
		go func(gen generator.Generator, err *error) {
			*err = errors.Join(*err, gen.Generate())
		}(gen, &err)
	}
	if err != nil {
		logger.Errorf("Error while creating a report >> %v", err)
	}
	rs.wg.Wait()
	logger.Warnf("Report created")
	return err
}

func (rs *ReportService) loadData() ([][]string, []string, error) {
	var (
		numBatches                   = int(math.Ceil(rs.FinishDate.Sub(rs.StartDate).Seconds() / rs.TimeWindowSample.Seconds()))
		hits            []domain.Hit = make([]domain.Hit, 0, 3000*numBatches)
		dumpNames                    = make([]string, 0, len(rs.Generators))
		dumpNameMatcher              = make(map[string]bool, 10)
		data            [][]string
		err             error
		// aux vars
		i       = -1
		counter = 0
	)
	logger.Warnf("Starting to load data")
	actualStart := rs.StartDate
	for i < 0 {
		counter++
		logger.Warnf("Loading batch [%d/%d]", counter, numBatches)
		lr, err := rs.LogAPIClient.Get(rs.Query, actualStart, actualStart.Add(rs.TimeWindowSample))
		if err != nil {
			logger.Errorf("error while getting data from kibana >> %v", err)
			actualStart = actualStart.Add(rs.TimeWindowSample)
			i = actualStart.Compare(rs.FinishDate)
			continue
		}
		newStart := actualStart.Add(rs.TimeWindowSample)
		if *lr.Log.Total == 0 {
			logger.Errorf("no data found between >> %s - %s", actualStart.String(), newStart.String())
		}
		actualStart = newStart
		i = actualStart.Compare(rs.FinishDate)
		hits = append(hits, lr.Log.Hits...)
	}
	data = make([][]string, 0, len(hits))
	for _, hit := range hits {
		matchs, matchErr := rs.Scraper.Scrape(*hit.Source.Message)
		if err != nil {
			errors.Join(err, matchErr)
			continue
		}
		key, err := buildDumpKey(matchs)
		if err != nil {
			return nil, nil, err
		}
		if _, ok := dumpNameMatcher[key]; !ok {
			dumpNames = append(dumpNames, key)
		}
		data = append(data, matchs)
	}
	return data, dumpNames, err
}
func sortOut(data [][]string, dumpNames ...string) map[string]chan []string {
	var (
		chanMap = make(map[string]chan []string, 10)
	)
	chanMap[fullDataKey] = make(chan []string)
	for _, dumpNames := range dumpNames {
		chanMap[dumpNames] = make(chan []string)
	}
	go func() {
		defer func() {
			for _, ch := range chanMap {
				close(ch)
			}
		}()
		for _, row := range data {
			key, err := buildDumpKey(row)
			if err != nil {
				if errors.Is(err, ErrNotEnoughColumns) {
					break
				}
				logger.Errorf("unable to build a key >> %v", err)
			}
			chanMap[key] <- row
			chanMap[fullDataKey] <- row
		}
	}()

	return chanMap
}

func buildDumpKey(row []string) (string, error) {
	if len(row) < 3 {
		return "", ErrNotEnoughColumns
	}
	if attr, resource := row[len(row)-attributePos], row[len(row)-resourcePos]; attr != "" && resource != "" {
		return attr + resource, nil
	}
	return "", fmt.Errorf("emty attribute or resource column for row %v", row)
}
