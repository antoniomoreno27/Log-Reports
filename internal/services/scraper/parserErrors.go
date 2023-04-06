package scraper

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
)

var (
	ColumnName [5]string = [5]string{
		"id",
		"site",
		"errMessage",
		"attribute",
		"resource",
	}
	parserErrorRegExpMap map[string]*regexp.Regexp = map[string]*regexp.Regexp{
		ColumnName[0]: regexp.MustCompile("/(M[A-Z]{2})?[0-9]+/"),
		ColumnName[1]: regexp.MustCompile(`site:(M[A-Z]{2})\]`),
		ColumnName[2]: regexp.MustCompile("ERROR:(.*?)-"),
		ColumnName[3]: regexp.MustCompile(`attribute:(.*?)\]`),
		ColumnName[4]: regexp.MustCompile(`resource:(.*?)\]`),
	}
)

type parseErrorScraperService struct {
	ColNames  []string
	regExpMap map[string]*regexp.Regexp
}

func NewParserErrorScrapper() Scraper {
	return &parseErrorScraperService{
		regExpMap: parserErrorRegExpMap,
		ColNames:  ColumnName[:],
	}
}

func (PES *parseErrorScraperService) Match(input string) ([]string, error) {
	var (
		data     = make([]string, 0, len(PES.ColNames))
		match    string
		matchErr error
		err      error
	)
	for _, id := range PES.ColNames {
		matcher := PES.regExpMap[id]
		if matcher == nil {
			return nil, fmt.Errorf("regexp not found to %s", id)
		}
		match, matchErr = PES.match(input, matcher)
		if matchErr != nil {
			err = errors.Join(err, fmt.Errorf("while matching column %s >> %v", id, matchErr))
			data = append(data, "")
			continue
		}
		match = clean(match)
		data = append(data, match)
	}
	if err != nil {
		log.Default().Println(err)
	}
	return data, nil
}

func (PES *parseErrorScraperService) match(input string, rExp *regexp.Regexp) (match string, err error) {
	match = rExp.FindString(input)
	if match == "" {
		err = errors.Join(err, errors.New("coincidence not found"))
		return
	}
	return
}

func clean(input string) string {
	index := strings.Index(input, ":")
	if index > -1 {
		input = input[index:]
	}
	input = input[1:]
	input = input[:len(input)-1]
	return input
}
