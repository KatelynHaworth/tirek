package main

import (
	"sync"
	"time"
	"fmt"
)

// Statistics stores information about
// a DNS load test in a synchronise way
// while providing a detailed string output
type Statistics struct {
	startTime          time.Time
	totalRequests      int64
	successfulRequests int64
	failedRequests     int64
	resolutionTimes    []time.Duration
	lock               sync.RWMutex
}

// Start prepares the statistics struct
// by setting the start time
func (stats *Statistics) Start() {
	stats.startTime = time.Now()
}

// IncrementSuccessfulRequests synchronously increments
// the count of total requests and successful requests,
// it will also add the time to resolve the request
func (stats *Statistics) IncrementSuccessfulRequests(d time.Duration) {
	defer stats.lock.Unlock()
	stats.lock.Lock()

	stats.totalRequests += 1
	stats.successfulRequests += 1
	stats.resolutionTimes = append(stats.resolutionTimes, d)
}

// IncrementSuccessfulRequests synchronously increments
// the count of total requests and failed requests
func (stats *Statistics) IncrementFailedRequests() {
	defer stats.lock.Unlock()
	stats.lock.Lock()

	stats.totalRequests += 1
	stats.failedRequests += 1
}

// String will lock the statistics for reading
// and will produce a string output of all the
// statistics in an easy to read format
func (stats *Statistics) String() string {
	defer stats.lock.RUnlock()
	stats.lock.RLock()

	sinceStart := time.Since(stats.startTime)
	totalTime := time.Duration(0)
	for i := range stats.resolutionTimes {
		totalTime += stats.resolutionTimes[i]
	}

	return fmt.Sprintf(
		"Run Time: %ds | Total Requests: %d (Errors %d / %.02f%%) | Average Response Time: %s | Requests Per Second: %.00f",
		int(sinceStart.Seconds()),
		stats.totalRequests,
		stats.failedRequests,
		(float64(stats.failedRequests) / float64(stats.totalRequests)) * 100,
		totalTime / time.Duration(stats.successfulRequests),
		float64(stats.successfulRequests) / float64(sinceStart / time.Second),
	)
}