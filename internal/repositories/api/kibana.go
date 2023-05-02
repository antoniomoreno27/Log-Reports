package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/antoniomoreno27/logs-LS/internal/domain"
	"github.com/antoniomoreno27/logs-LS/internal/services/logger"
)

const (
	KibanaUrl string = "https://furylogs-assorted-17.furycloud.io/elasticsearch/fury-listing-sort-api-*/_search?rest_total_hits_as_int=true&ignore_unavailable=true&ignore_throttled=true&preference=1679926399819&timeout=30000ms"
)

type kibanaAPIClient struct {
	Client *http.Client
	Cookie string
}

func NewLogAPIClient(cookie string, mock bool) LogAPIClient {
	if !mock {
		return &kibanaAPIClient{
			Client: &http.Client{},
			Cookie: cookie,
		}
	}
	return &kibanaAPIClientMock{
		GetMock: func(query string, startDate, finishDate time.Time) (logReport *domain.LogReport, err error) {
			logReport = &domain.LogReport{}
			if file, err := os.ReadFile("./internal/repositories/api/kibana_mock.json"); err == nil {
				err = json.Unmarshal(file, logReport)
				if err != nil {
					return nil, err
				}
				return logReport, err
			} else {
				logger.Errorf("Couldn`t load mock file: %v", err)
			}
			return
		},
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
				Must: []domain.Must{
					{
						QueryString: domain.QueryString{
							AnalizeWildcard: true,
							Query:           "*",
							TimeZone:        "America/Bogota",
						},
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
func (kbn *kibanaAPIClient) Get(query string, startDate, finishDate time.Time) (logReport *domain.LogReport, err error) {
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
	err = checkResponse(response)
	if err != nil {
		return
	}
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
func checkResponse(response *http.Response) error {
	if response == nil {
		return fmt.Errorf("could not retrieve data")
	}
	if response.StatusCode == http.StatusNotFound {
		return fmt.Errorf("data not found")
	}
	if !(response.StatusCode >= http.StatusOK && response.StatusCode <= http.StatusIMUsed) {
		return fmt.Errorf("invalid response from api (status = %d '%s')", response.StatusCode, response.Status)
	}
	return nil
}
