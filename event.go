package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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
func (e *event) filename() string {
	return fmt.Sprintf("%v_%v_%v.csv", e.location, e.number, e.date)
}

// writeCSV writes the event's data to a CSV file.
// The filename method is used to determine the filename used.
// Any existing file will be overwritten.
func (e *event) writeCSV() {
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

// parseEventFilename reads the event metadata from a saved CSV filename.
func parseEventFilename(filename string) (location string, number int, date string) {
	// TODO: Handle errors.
	filename = strings.ReplaceAll(filename, ".csv", "")
	parts := strings.Split(filename, "_")
	number, _ = strconv.Atoi(parts[1])
	return parts[0], number, parts[2]
}

// newEventFromFilename loads the event data from the
// given filename and returns a pointer to the new instance.
func newEventFromFilename(filename string) *event {
	// TODO: Error handling add / upgrade.
	location, eventNum, date := parseEventFilename(filename)
	newEvent := event{
		location: location,
		number:   eventNum,
		date:     date,
	}
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal("Unable to read input file "+filename, err)
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filename, err)
	}
	for i, r := range data {
		// TODO: REVIEW: Could use reflection to handle mis-ordered columns.
		if i == 0 {
			// TODO: Add checking for header line. -> def as const? --> gen from struct values??
			// TODO: REVIEW: Could check for column names and ignore odd columns.
			continue
		}
		id, _ := strconv.Atoi(r[0])
		clubId, _ := strconv.Atoi(r[4])
		position, _ := strconv.Atoi(r[6])
		runs, _ := strconv.Atoi(r[7])
		ageGrade, _ := strconv.ParseFloat(r[8], 32)
		newEvent.results = append(newEvent.results, result{
			id:          id,
			name:        r[1],
			ageGroup:    r[2],
			club:        r[3],
			clubId:      clubId,
			gender:      r[5],
			position:    position,
			runs:        runs,
			ageGrade:    float32(ageGrade),
			achievement: r[9],
			time:        r[10],
			currentPB:   r[11],
		})
	}
	return &newEvent
}

// TODO: Add no. first timer's method
// TODO: Add no. PBs (inc 1sts and 1sts as 2nd return?) method
// TODO: Add total no. method.
// TODO: Add no. clubs, most rep club methods?
// TODO: Add no. M/F methods.
// TODO: Add fastest time, fastest M/F, slowest time, slowest M/F methods?
// TODO: Add method to report on age group %'s?
// TODO: Add no. unknown method?
