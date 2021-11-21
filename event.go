package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

// event is a single Parkrun event including metadata and runner results.
type event struct {
	location string
	number   int
	date     string // Format dd-mm-yyy.
	results  []result
}

// filename returns a CSV filename to use for storing the event's results data.
// The filename also acts as the store of the event metadata.
func (e event) filename() string {
	return fmt.Sprintf("%v_%v_%v.csv", e.location, e.number, e.date)
}

// writeCSV writes the event's data to a CSV file.
// The filename method is used to determine the filename used.
// Any existing file will be overwritten.
func (e event) writeCSV() {
	data := make([][]string, 0, resultsCapacity+1)
	data = append(data, []string{"id", "name", "ageGroup", "club", "clubId", "gender", "position", "runs", "ageGrade", "achievement", "time", "currentPB"})
	for _, r := range e.results {
		data = append(data, []string{
			strconv.Itoa(r.id),
			r.name,
			r.ageGroup,
			r.club,
			strconv.Itoa(r.clubId),
			r.gender,
			strconv.Itoa(r.position),
			strconv.Itoa(r.runs),
			fmt.Sprintf("%.2f", r.ageGrade),
			r.achievement,
			r.time,
			r.currentPB,
		})
	}
	// TODO: Handle custom file location?
	file, err := os.Create(e.filename())
	if err != nil {
		log.Fatal("Cannot create file", err)
		// TODO: Update error handling.
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	for _, value := range data {
		err := writer.Write(value)
		if err != nil {
			log.Fatal("Cannot write to file", err)
			// TODO: Update error handling.
		}
	}
}

// TODO: Add no. first timer's method
// TODO: Add no. PBs (inc 1sts and 1sts as 2nd return?) method
// TODO: Add total no. method.
// TODO: Add no. clubs, most rep club methods?
// TODO: Add no. M/F methods.
// TODO: Add fastest time, fastest M/F, slowest time, slowest M/F methods?
// TODO: Add method to report on age group %'s?
// TODO: Add no. unknown method?
