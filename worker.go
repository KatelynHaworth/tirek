package main

import (
	"time"
	"github.com/miekg/dns"
	"math/rand"
)

// Worker is a single DNS lookup worker
// tasked with making x amount of DNS requests
// per second for a random domain against the
// provided name server
type Worker struct {
	Stats      *Statistics
	Rate       int
	NameServer string
	Domains    []string

	ticker *time.Ticker
	exit   chan struct{}
}

// Start prepares the worker instance and starts
// the internal go routine loop that will make
// requests
func (worker *Worker) Start() {
	worker.exit = make(chan struct{}, 0)
	worker.ticker = time.NewTicker(time.Second / time.Duration(worker.Rate))
	go worker.loop()
}

// getRandomDomain will lookup a random domain
// from the list of provided domain names
func (worker *Worker) getRandomDomain() string {
	return worker.Domains[rand.Intn(len(worker.Domains))]
}

// loop will create and send a DNS lookup
// every time the ticker ticks and will log
// the result against the statistics struct
func (worker *Worker) loop() {
workerLoop:
	for {
		select {
		case start := <- worker.ticker.C:
			req := new(dns.Msg).SetQuestion(worker.getRandomDomain(), dns.ClassINET)

			resp, err := dns.Exchange(req, worker.NameServer)
			switch {
			case err != nil:
				fallthrough

			case resp.Rcode != dns.RcodeSuccess:
				worker.Stats.IncrementFailedRequests()

			default:
				worker.Stats.IncrementSuccessfulRequests(time.Since(start))
			}

		case <- worker.exit:
			break workerLoop
		}
	}
}

// Stop will inform the internal go
// routine to stop looping and to return
func (worker *Worker) Stop() {
	worker.exit <- struct{}{}
	worker.ticker.Stop()
}