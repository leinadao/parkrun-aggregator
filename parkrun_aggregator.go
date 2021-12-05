// The parkrunaggregator program is for aggregating and interrogating Parkrun data.
// It stores data locally to scrape an event's result only once and reduce site load.
// Written as an example go project, it is not officially supported and not for re-use.
package main

import (
	"fmt"
	"sort"
	"time"
)

// TODO: Add test code!

// findOccurancesPerRunner prints out each runner and their total number of runs,
// in ascending order, from the given slice of event pointers.
func findOccurancesPerRunner(eventPs []*event) {
	numPerRunner := map[runner]int{}
	for _, eP := range eventPs {
		for _, r := range eP.results {
			_, ok := numPerRunner[r.runner]
			if ok {
				numPerRunner[r.runner] += 1 // TODO: REVIEW: Can just have this line without an ok check due to the default 0?
			} else {
				numPerRunner[r.runner] = 1
			}
		}
	}
	runnersPerNum := map[int][]runner{}
	for r, n := range numPerRunner {
		runnersPerNum[n] = append(runnersPerNum[n], r)
	}
	var nums []int
	for n := range runnersPerNum {
		nums = append(nums, n)
	}
	sort.Ints(nums)
	for _, n := range nums {
		for _, r := range runnersPerNum[n] {
			fmt.Printf("%s (%d), %d\n", r.name, r.id, n)
		}
	}
}

// main takes a Parkrun location and retrieves any outstanding data.
// Data is written to a CSV file per event.
func main() {
	// Take in a location name:
	fmt.Println("Enter the full Parkrun location name. e.g. 'bathskyline': ") // TODO: Split off this input as own fn.
	var location string
	fmt.Scanln(&location)
	fmt.Printf("Location is %v...\n", location)
	var eventPs []*event
	incomplete := true
	// Load or fetch all possible events for the location:
	for eN, tmpLimiter := 1, 0; tmpLimiter < 100; eN++ {
		var (
			eP  *event // Needed for use outside if, if set in if.
			err error
		)
		eP, err = loadEventCSV(location, eN)
		if err != nil {
			tmpLimiter += 1              // TODO: REVIEW: Remove limiter?
			time.Sleep(10 * time.Second) // TODO: TEMP?
			eP, err = getEvent(location, eN)
			if err != nil {
				fmt.Println("No more events to fetch.")
				incomplete = false
				break
			}
			fmt.Printf("Fetched event %v: %v...\n", eP.number, eP.date)
			eP.writeCSV()
		} else {
			fmt.Printf("Loaded event %v: %v...\n", eP.number, eP.date)
		}
		eventPs = append(eventPs, eP)
	}
	if !incomplete {
		// TODO: TEMP: print ordered list of most frequent runners:
		findOccurancesPerRunner(eventPs)
	}
}
