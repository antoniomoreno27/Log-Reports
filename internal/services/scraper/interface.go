package scraper

type Scraper interface {
	Scrape(target string) (data []string, err error)
}

type Matcher interface {
	Match(target string) (match string, err error)
	String() (name string)
}
