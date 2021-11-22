// The parkrunaggregator program is for aggregating and interrogating Parkrun data.
// It stores data locally to scrape an event's result only once and reduce site load.
// Written as an example go project, it is not officially supported and not for re-use.
package main

import (
	"fmt"
)

// TODO: REVIEW: Create runner objects?
// TODO: Add test code!

// main takes a Parkrun location and retrieves any outstanding data.
// Data is written to a CSV file per event.
func main() {
	// Take in a location name:
	fmt.Println("Enter the full Parkrun location name. e.g. 'bathskyline': ")
	var location string
	fmt.Scanln(&location)
	fmt.Printf("Location is %v...\n", location)
	var eventPs []*event
	// Load or fetch all possible events for the location:
	for eN, tmpLimiter := 1, 0; tmpLimiter < 5; eN++ {
		eP, err := loadEventCSV(location, eN)
		if err != nil {
			tmpLimiter += 1 // TODO: REVIEW: Remove limiter?
			eP, err := getEvent(location, eN)
			if err != nil {
				return
			}
			eP.writeCSV()
		}
		eventPs = append(eventPs, eP)
	}
	fmt.Println(len(eventPs))
}
