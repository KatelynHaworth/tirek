package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/miekg/dns"
)

var (
	cpus		*int		= flag.Int("cpus", runtime.NumCPU(), "Sets the number of CPUs to be used during test")

	workers		*int		= flag.Int("workers", 5, "Sets the number of worker threads to use")

	rate		*int		= flag.Int("rate", 50, "Sets the number of requests to attempt per second")

	duration	*time.Duration	= flag.Duration("duration", 1 * time.Minute, "Sets how long the attack should last")

	target		*string		= flag.String("target", "127.0.0.1:53", "Sets the DNS server to target")

	list		[]string
)

func main() {
	flag.Parse()
	getRandomQuestion() // Populate list before the workers start

	runtime.GOMAXPROCS(*cpus)
	workers := make([]*Attacker, *workers)
	for i := range workers {
		workers[i] = NewAttacker(*rate)
	}

	for _, worker := range workers {
		go worker.Start()
	}

	<- time.After(*duration)

	for _, worker := range workers {
		worker.exit <- true
	}

	totalResults := []time.Duration{}
	totalErrors := 0
	for _, worker := range workers {
		totalResults = append(totalResults, worker.results...)
		totalErrors += worker.errors
	}

	totalTime := time.Duration(0)
	for _, result := range totalResults {
		totalTime += result
	}

	fmt.Printf(
		"Total Requests: %d (Errors %d) | Average Response Time: %s",
		len(totalResults),
		totalErrors,
		totalTime / time.Duration(len(totalResults)),
	)
}

type Attacker struct {
	rate		int
	results		[]time.Duration
	errors		int
	exit		chan bool
}

func NewAttacker(rate int) *Attacker {
	attacker := new(Attacker)
	attacker.rate = rate
	attacker.exit = make(chan bool, 1)
	attacker.errors = 0
	return attacker
}

func (self *Attacker) Start() {
	ticker := time.NewTicker(time.Duration(time.Second / time.Duration(50)))

	attack:
	for {
		select {
		case tm := <- ticker.C:
			_, err := dns.Exchange(getRandomQuestion(), *target)
			if err != nil {
				self.errors++
			}
			self.results = append(self.results, time.Since(tm))

		case <- self.exit:
			ticker.Stop()
			break attack
		}
	}
}

func getRandomQuestion() *dns.Msg {
	if list == nil {
		list = []string{}
		resp, _ := http.Get("https://raw.githubusercontent.com/opendns/public-domain-lists/master/opendns-random-domains.txt")
		reader := bufio.NewReader(resp.Body)
		for {
			line, _, err := reader.ReadLine()
			if err != nil {
				break
			}
			list = append(list, string(line))
		}
		resp.Body.Close()
	}

	msg := new(dns.Msg)
	msg.SetQuestion(list[rand.Intn(len(list))] + ".", dns.TypeA)
	return msg
}