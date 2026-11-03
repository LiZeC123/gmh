package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptrace"
	"strings"
	"sync"
	"time"

	"github.com/urfave/cli/v3"
)

// Result holds measurement results for a single request
type Result struct {
	ConnectDur time.Duration
	RequestDur time.Duration
	Err        error
}

func LatencyCommand() *cli.Command {
	return &cli.Command{
		Name:  "latency",
		Usage: "Measure HTTP latency (connection and request duration)",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "url",
				Aliases:  []string{"u"},
				Required: true,
			},
			&cli.Uint8Flag{
				Name:     "rounds",
				Aliases:  []string{"r"},
				Value:    1,
				Required: false,
			},
			&cli.Uint8Flag{
				Name:     "concurrency",
				Aliases:  []string{"c"},
				Value:    1,
				Required: false,
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			url := c.String("url")
			rounds := c.Uint8("rounds")
			concurrency := c.Uint8("concurrency")

			// Create HTTP client, disable Keep-Alive to ensure each request establishes a new connection
			tr := &http.Transport{
				DisableKeepAlives: true,
			}
			client := &http.Client{Transport: tr}

			// Execute per round
			for round := uint8(1); round <= rounds; round++ {
				if concurrency == 1 {
					// Sequential execution for concurrency=1, print each result
					connectDur, requestDur, err := doRequest(url, client)
					if err != nil {
						log.Printf("Round %d request failed: %v", round, err)
						continue
					}
					fmt.Printf("Round %d: connection duration=%v, request duration=%v\n", round, connectDur, requestDur)
				} else {
					// Concurrent requests: fire all at once, wait for completion, then output statistics
					results := make([]Result, concurrency)
					var wg sync.WaitGroup
					wg.Add(int(concurrency))

					for i := range concurrency {
						go func(idx uint8) {
							defer wg.Done()
							cd, rd, err := doRequest(url, client)
							results[idx] = Result{ConnectDur: cd, RequestDur: rd, Err: err}
						}(i)
					}
					wg.Wait()

					// Collect successful results
					var connectDurs []time.Duration
					var requestDurs []time.Duration
					for _, r := range results {
						if r.Err != nil {
							log.Printf("Round %d request error: %v", round, r.Err)
							continue
						}
						connectDurs = append(connectDurs, r.ConnectDur)
						requestDurs = append(requestDurs, r.RequestDur)
					}

					if len(connectDurs) == 0 {
						log.Printf("Round %d has no successful requests", round)
						continue
					}

					// Calculate statistics
					minC, maxC, avgC := stats(connectDurs)
					minR, maxR, avgR := stats(requestDurs)

					fmt.Printf("Round %d (concurrency=%d):\n", round, concurrency)
					fmt.Printf("  Connection duration: min=%v, max=%v, avg=%v\n", minC, maxC, avgC)
					fmt.Printf("  Request duration: min=%v, max=%v, avg=%v\n", minR, maxR, avgR)
				}
			}
			return nil
		},
	}
}

// doRequest performs a single HTTP request and returns connection duration and request duration
func doRequest(url string, client *http.Client) (connectDur, requestDur time.Duration, err error) {
	var gotConnTime time.Time

	// Create ClientTrace to capture connection establishment time
	trace := &httptrace.ClientTrace{
		GotConn: func(info httptrace.GotConnInfo) {
			gotConnTime = time.Now()
		},
	}
	ctx := httptrace.WithClientTrace(context.Background(), trace)

	// Build POST request with a simple body (echo server will echo it back)
	body := strings.NewReader("hello")
	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return 0, 0, err
	}
	req.Header.Set("Content-Type", "text/plain")

	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	// Read response body (discard content, just ensure it is fully received)
	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		return 0, 0, err
	}
	end := time.Now()

	if gotConnTime.IsZero() {
		return 0, 0, fmt.Errorf("GotConn event not triggered (connection may have been reused)")
	}

	total := end.Sub(start)
	connectDur = gotConnTime.Sub(start)
	requestDur = total - connectDur
	return connectDur, requestDur, nil
}

// stats calculates the minimum, maximum and average of a slice of durations
func stats(durs []time.Duration) (min, max, avg time.Duration) {
	if len(durs) == 0 {
		return 0, 0, 0
	}
	min = durs[0]
	max = durs[0]
	sum := time.Duration(0)
	for _, d := range durs {
		if d < min {
			min = d
		}
		if d > max {
			max = d
		}
		sum += d
	}
	avg = sum / time.Duration(len(durs))
	return
}
