package passive

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/projectdiscovery/gologger"
	"github.com/ShangRui-hash/subfinder/v2/pkg/subscraping"
)

// EnumerateSubdomains enumerates all the subdomains for a given domain
func (a *Agent) EnumerateSubdomains(domain string, keys *subscraping.Keys, proxy string, rateLimit, timeout int, maxEnumTime time.Duration, localIP net.IP) chan subscraping.Result {
	results := make(chan subscraping.Result)

	go func() {
		session, err := subscraping.NewSession(domain, keys, proxy, rateLimit, timeout, localIP)
		if err != nil {
			results <- subscraping.Result{Type: subscraping.Error, Error: fmt.Errorf("could not init passive session for %s: %s", domain, err)}
		}

		ctx, cancel := context.WithTimeout(context.Background(), maxEnumTime)

		timeTaken := make(map[string]string)
		timeTakenMutex := &sync.Mutex{}

		wg := &sync.WaitGroup{}
		// Run each source in parallel on the target domain
		for source, runner := range a.sources {
			wg.Add(1)

			now := time.Now()
			go func(source string, runner subscraping.Source) {
				for resp := range runner.Run(ctx, domain, session) {
					results <- resp
				}

				duration := time.Since(now)
				timeTakenMutex.Lock()
				timeTaken[source] = fmt.Sprintf("Source took %s for enumeration\n", duration)
				timeTakenMutex.Unlock()

				wg.Done()
			}(source, runner)
		}
		wg.Wait()

		for source, data := range timeTaken {
			gologger.Verbose().Label(source).Msg(data)
		}

		close(results)
		cancel()
	}()

	return results
}
