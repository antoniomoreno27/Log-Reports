package domain

import "regexp"

type Matcher struct {
	name   string
	regexp *regexp.Regexp
}
