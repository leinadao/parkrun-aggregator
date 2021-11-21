package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

type event struct {
	location string
	number   int
	date     string // Format dd-mm-yyy.
	results  []result
}

func (e event) filename() string {
	return fmt.Sprintf("%v_%v_%v.csv", e.location, e.number, e.date)
}

func (e event) writeCSV() {
	data := make([][]string, 0, resultsCapacity+1)
	data = append(data, []string{"id", "name", "ageGroup", "club", "clubId", "gender", "position", "runs", "ageGrade", "achievement", "time", "previousPB"})
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
			fmt.Sprintf("%f", r.ageGrade), // TODO: Result here needs rounding off.
			r.achievement,
			r.time,
			r.previousPB,
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
