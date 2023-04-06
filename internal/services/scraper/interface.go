package scraper

type Scraper interface {
	Match(input string) (data []string, err error)
}
