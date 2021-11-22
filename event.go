package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

// parseEventFilename reads and returns the event metadata from a saved CSV filename.
func parseEventFilename(filename string) (location string, number int, date string) {
	// TODO: Handle errors.
	filename = strings.ReplaceAll(filename, ".csv", "")
	parts := strings.Split(filename, "_")
	number, _ = strconv.Atoi(parts[1])
	return parts[0], number, parts[2]
}

// newEventFromCSV loads the event data from the
// given filename and returns a pointer to the new instance.
func newEventFromCSV(filename string) *event {
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

// loadEvent tries to load and return an event pointer from
// a CSV file based only on the event location and number given.
func loadEventCSV(location string, eventNum int) (*event, error) {
	tmpE := event{location: location, number: eventNum, date: "*"}
	matches, err := filepath.Glob(tmpE.filename())
	if err != nil {
		return nil, err
	}
	if matches == nil {
		return nil, errors.New("no files found")
	}
	if len(matches) != 1 {
		return nil, errors.New("multiple files found")
	}
	eP := newEventFromCSV(matches[0])
	return eP, nil
}

// getEvent retrieves the available data for the given event number
// at the given Parkrun location.
// It returns a new event instance pointer.
func getEvent(location string, eventNum int) (*event, error) {
	resp, err := http.Get(fmt.Sprintf("https://www.parkrun.org.uk/%v/results/%v/", location, eventNum))
	if err != nil {
		// TODO: Add error handling.
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("not 200 response")
		// TODO: Improve error handling.
		// log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		// TODO: Change error handling?
		log.Fatal(err)
	}
	docRes := doc.Find(".Results-table-row")
	if docRes.Length() == 0 {
		return nil, errors.New("no results found")
	}
	date := strings.ReplaceAll(doc.Find(".Results-header").Find(".format-date").Text(), "/", "-")
	results := make([]result, 0, resultsCapacity)
	docRes.Each(func(i int, s *goquery.Selection) {
		idStr := s.Find(".Results-table-td.Results-table-td--name").Find(".compact").Find("a").AttrOr("href", "")
		id, _ := strconv.Atoi(idStr[strings.LastIndex(idStr, "/")+1:]) // TODO: Handle error.
		clubIdStr := s.Find(".Results-table-club.Results-tablet").Find(".detailed").Find("a").AttrOr("href", "")
		clubId, _ := strconv.Atoi(clubIdStr[strings.LastIndex(clubIdStr, "=")+1:])                  // TODO: Handle error.
		currentPB, _ := s.Find(".Results-table-td.Results-table-td--time").Find(".detailed").Html() // TODO: Handle error.
		currentPB = strings.ReplaceAll(currentPB, "<span class=\"Results-table--normal\">PB", "")
		currentPB = strings.ReplaceAll(currentPB, "<span class=\"Results-table--red\">", "")
		currentPB = strings.ReplaceAll(currentPB, "<span class=\"Results-table--green\">", "") // TODO: Group using regex?
		currentPB = strings.ReplaceAll(currentPB, "</span>", "")
		currentPB = strings.ReplaceAll(currentPB, "&nbsp;", "")
		time := s.Find(".Results-table-td.Results-table-td--time").Find(".compact").Text()
		// Handle text in PB field:
		if currentPB == "New PB!" || currentPB == "First Timer!" {
			currentPB = time
		}
		position, _ := strconv.Atoi(s.AttrOr("data-position", ""))           // TODO: Handle error.
		runs, _ := strconv.Atoi(s.AttrOr("data-runs", ""))                   // TODO: Handle error.
		ageGrade, _ := strconv.ParseFloat(s.AttrOr("data-agegrade", ""), 32) // Still made as float64 but convertable to 32. // TODO: Handle error.
		results = append(results, result{
			id:          id,
			clubId:      clubId,
			currentPB:   currentPB,
			name:        s.AttrOr("data-name", ""),
			ageGroup:    s.AttrOr("data-agegroup", ""),
			club:        s.AttrOr("data-club", ""),
			gender:      s.AttrOr("data-gender", ""),
			position:    position,
			runs:        runs,
			ageGrade:    float32(ageGrade),
			achievement: s.AttrOr("data-achievement", ""),
			time:        time,
		})
	})
	return &event{
		location: location,
		number:   eventNum,
		date:     date,
		results:  results,
	}, nil
}

// TODO: Add no. first timer's method
// TODO: Add no. PBs (inc 1sts and 1sts as 2nd return?) method
// TODO: Add total no. method.
// TODO: Add no. clubs, most rep club methods?
// TODO: Add no. M/F methods.
// TODO: Add fastest time, fastest M/F, slowest time, slowest M/F methods?
// TODO: Add method to report on age group %'s?
// TODO: Add no. unknown method?
