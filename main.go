package main

import (
	"bufio"
	"flag"
	"net/http"
	"runtime"
	"time"

	"fmt"
	"log"
	"io"
	"strings"
)

var (
	listSource = flag.String("list", "https://raw.githubusercontent.com/opendns/public-domain-lists/master/opendns-top-domains.txt", "Sets the list of domains to use in requests")
	cpus = flag.Int("cpus", runtime.NumCPU(), "Sets the number of CPUs to be used during test")
	workerCount = flag.Int("workers", 5, "Sets the number of worker threads to use")
	rate = flag.Int("rate", 50, "Sets the number of requests to attempt per second")
	duration = flag.Duration("duration", 1 * time.Minute, "Sets how long the attack should last")
	target = flag.String("target", "127.0.0.1:53", "Sets the DNS server to target")

	domainList []string
)

func main() {
	flag.Parse()
	if err := loadDomainList(); err != nil {
		log.Fatal(err)
	}

	stats := new(Statistics)

	runtime.GOMAXPROCS(*cpus)
	workers := make([]*Worker, *workerCount)
	for i := range workers {
		workers[i] = &Worker{
			Stats:      stats,
			Rate:       *rate,
			NameServer: *target,
			Domains:    domainList,
		}
	}

	fmt.Println("CPU Threads....................:", *cpus)
	fmt.Println("Workers........................:", *workerCount)
	fmt.Println("Rate (per second per worker)...:", *rate)
	fmt.Println("Duration.......................:", *duration)
	fmt.Println("Target Name Server.............:", *target)
	fmt.Println("Available domains..............:", len(domainList))
	fmt.Println()

	stats.Start()
	for _, worker := range workers {
		worker.Start()
	}

	timeUp := time.After(*duration)
	interval := time.NewTicker(time.Second / 16)

	for {
		select {
		case <- interval.C:
			fmt.Print(stats, "                    \r")

		case <- timeUp:
			interval.Stop()
			for _, worker := range workers {
				worker.Stop()
			}

			fmt.Println()
			return
		}
	}
}

// loadDomainList will fetch the list of
// domain names and parse it into a slice
// of FQDNs ready for using in requests
func loadDomainList() error {
	resp, err := http.Get(*listSource)
	if err != nil {
		return fmt.Errorf("retreive domain list: %s", err)
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err == io.EOF {
			break
		}else if err != nil {
			return fmt.Errorf("read line from list: %s", err)
		}

		domainList = append(domainList, strings.TrimSpace(line)+".")
	}

	return nil
}
