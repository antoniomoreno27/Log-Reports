package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/antoniomoreno27/logs-LS/internal/domain"
)

const (
	KibanaUrl string = "https://furylogs-assorted-17.furycloud.io/elasticsearch/fury-listing-sort-api-*/_search?rest_total_hits_as_int=true&ignore_unavailable=true&ignore_throttled=true&preference=1679926399819&timeout=30000ms"
)

type KibanaAPIClient struct {
	Client *http.Client
	Cookie string
}

func New(cookie string) *KibanaAPIClient {
	return &KibanaAPIClient{
		Client: &http.Client{},
		Cookie: cookie,
	}
}

func buildQueryBody(query string, startDate, finishDate time.Time) (bufferPostBody *bytes.Buffer, err error) {
	postBody, err := json.Marshal(domain.PostBody{
		Size: 3000,
		Sort: []domain.Sort{
			{
				Timestamp: domain.Timestamp{
					Order:        "desc",
					UnmappedType: "boolean",
				},
			},
		},
		DocValueFields: []domain.DocValueFields{
			{
				Field:  "timestamp",
				Format: "date_time",
			},
		},
		Query: domain.Query{
			Bool: domain.Bool{
				Must: domain.Must{
					QueryString: domain.QueryString{
						AnalizeWildcard: true,
						Query:           "*",
						TimeZone:        "America/Bogota",
					},
				},
				Filter: []domain.Match{
					{
						MatchPhrase: &domain.MatchPhrase{
							Message: query,
						},
					},
					{
						Range: &domain.Range{
							RangeTimestamp: domain.RangeTimestamp{
								GTE:    startDate,
								LTE:    finishDate,
								Format: "strict_date_optional_time",
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		return
	}
	bufferPostBody = bytes.NewBuffer(postBody)
	return
}
func (kbn *KibanaAPIClient) Get(query string, startDate, finishDate time.Time) (logReport domain.LogReport, err error) {
	var (
		body         *bytes.Buffer
		request      *http.Request
		response     *http.Response
		responseBody []byte
	)
	body, err = buildQueryBody(query, startDate, finishDate)
	if err != nil {
		return
	}
	request, err = http.NewRequest(http.MethodPost, KibanaUrl, body)
	if err != nil {
		return
	}
	request.Header.Set("cookie", kbn.Cookie)
	request.Header.Set("kbn-version", "7.6.1")
	response, err = kbn.Client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	responseBody, err = io.ReadAll(response.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(responseBody, &logReport)
	if err != nil {
		return
	}
	return
}
