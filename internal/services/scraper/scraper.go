package scraper

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/antoniomoreno27/logs-LS/internal/services/logger"
)

type matcher struct {
	Name   string
	Regexp *regexp.Regexp
}

func NewRegExpMatcher(name string, regExp *regexp.Regexp) Matcher {
	return &matcher{
		Name:   name,
		Regexp: regExp,
	}
}

func (m matcher) Match(target string) (match string, err error) {
	match = m.Regexp.FindString(target)
	if match == "" {
		err = fmt.Errorf("coincidence of %s regexp in %s not found", m.Name, target)
		return
	}
	//cleaning the reg exp match
	index := strings.Index(match, ":")
	if index > -1 {
		match = match[index:]
		if err = checkMatch(match, target); err != nil {
			return "", err
		}
	}
	match = match[1:]
	if err = checkMatch(match, target); err != nil {
		return "", err
	}
	match = match[:len(match)-1]
	return
}
func checkMatch(match, regexpName string) error {
	if len(match) == 0 {
		return fmt.Errorf("bad match %s at %s", match, regexpName)
	}
	return nil
}

func (m matcher) String() (name string) {
	return m.Name
}

type scraper struct {
	matchers []Matcher
}

func NewScraper(matchers ...Matcher) Scraper {
	return &scraper{
		matchers: matchers,
	}
}

func (s *scraper) Scrape(target string) (data []string, err error) {
	data = make([]string, 0, len(s.matchers))
	for _, matcher := range s.matchers {
		match, matchErr := matcher.Match(target)
		if matchErr != nil {
			err = errors.Join(err, fmt.Errorf("while matching column %s >> %v", matcher.String(), matchErr))
			data = append(data, "")
			continue
		}
		data = append(data, match)
	}
	if err != nil {
		logger.Errorf("errors while doing scrape: %v", err)
	}
	return data, nil
}
