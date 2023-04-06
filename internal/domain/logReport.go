package domain

type LogReport struct {
	Took int `json:"took"`
	Log  Log `json:"hits"`
}
type Log struct {
	Hits []Hit `json:"hits"`
}
type Hit struct {
	Source *Source `json:"_source"`
}
type Source struct {
	Message *string `json:"message"`
}
