// Package threatminer logic
package threatminer

import (
	"context"
	"fmt"

	jsoniter "github.com/json-iterator/go"

	"github.com/ShangRui-hash/subfinder/v2/pkg/subscraping"
)

type response struct {
	StatusCode    string   `json:"status_code"`
	StatusMessage string   `json:"status_message"`
	Results       []string `json:"results"`
}

// Source is the passive scraping agent
type Source struct{}

// Run function returns all subdomains found with the service
func (s *Source) Run(ctx context.Context, domain string, session *subscraping.Session) <-chan subscraping.Result {
	results := make(chan subscraping.Result)

	go func() {
		defer close(results)

		resp, err := session.SimpleGet(ctx, fmt.Sprintf("https://api.threatminer.org/v2/domain.php?q=%s&rt=5", domain))
		if err != nil {
			results <- subscraping.Result{Source: s.Name(), Type: subscraping.Error, Error: err}
			session.DiscardHTTPResponse(resp)
			return
		}

		defer resp.Body.Close()

		var data response
		err = jsoniter.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			results <- subscraping.Result{Source: s.Name(), Type: subscraping.Error, Error: err}
			return
		}

		for _, subdomain := range data.Results {
			results <- subscraping.Result{Source: s.Name(), Type: subscraping.Subdomain, Value: subdomain}
		}
	}()

	return results
}

// Name returns the name of the source
func (s *Source) Name() string {
	return "threatminer"
}
