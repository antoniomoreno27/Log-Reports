package domain

import "time"

type PostBody struct {
	Highlight      Highlight        `json:"highlight,omitempty"`
	Size           int              `json:"size,omitempty"`
	Sort           []Sort           `json:"sort,omitempty"`
	DocValueFields []DocValueFields `json:"docvalue_fields,omitempty"`
	Query          Query            `json:"query,omitempty"`
}

type Highlight struct {
	PreTags           []string       `json:"pre_tags,omitempty"`
	PostTags          []string       `json:"post_tags,omitempty"`
	Fields            map[string]any `json:"fields,omitempty,omitempty"`
	RequireFieldMatch bool           `json:"require_field_match,omitempty"`
	FragmentSize      int            `json:"fragment_size,omitempty"`
}

type Sort struct {
	Timestamp Timestamp `json:"timestamp,omitempty"`
}
type DocValueFields struct {
	Field  string `json:"field,omitempty"`
	Format string `json:"format,omitempty"`
}

type Timestamp struct {
	Order        string `json:"order,omitempty"`
	UnmappedType string `json:"unmapped_type,omitempty"`
}
type Query struct {
	Bool Bool `json:"bool,omitempty"`
}

type Bool struct {
	Must   Must    `json:"must,omitempty"`
	Filter []Match `json:"filter,omitempty"`
}

type Match struct {
	MatchAll    any          `json:"match_all,omitempty"`
	MatchPhrase *MatchPhrase `json:"match_phrase,omitempty"`
	Range       *Range       `json:"range,omitempty"`
}

type Must struct {
	QueryString QueryString `json:"query_string,omitempty"`
}

type QueryString struct {
	AnalizeWildcard bool   `json:"analyze_wildcard,omitempty"`
	Query           string `json:"query,omitempty"`
	TimeZone        string `json:"time_zone,omitempty"`
}
type MatchPhrase struct {
	TagContainer *TagContainer `json:"tags.container,omitempty"`
	Message      string        `json:"message,omitempty"`
}

type Range struct {
	RangeTimestamp RangeTimestamp `json:"timestamp,omitempty"`
}
type RangeTimestamp struct {
	GTE    time.Time `json:"gte,omitempty"`
	LTE    time.Time `json:"lte,omitempty"`
	Format string    `json:"format,omitempty"`
}
type TagContainer struct {
	Query string `json:"query,omitempty"`
}
