package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const resultsCapacity = 500 // TODO: REVIEW capacity of 500.

type result struct {
	id          int
	name        string
	ageGroup    string
	club        string
	clubId      int
	gender      string
	position    int
	runs        int
	ageGrade    float32
	achievement string
	time        string
	previousPB  string // TODO: Handle 'New PB etc. strings --> Set to time.'
}

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
	data := make([][]string, 0, resultsCapacity+1) // TODO: Try interface to see if conversions can be skipped?
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

// TODO: REVIEW: Create runner objects?

func getEvent(location string, eventNum int) *event {
	resp, err := http.Get(fmt.Sprintf("https://www.parkrun.org.uk/%v/results/%v/", location, eventNum))
	if err != nil {
		// TODO: Add error handling.
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		// TODO: Add error handling.
		// log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		// TODO: Change error handling?
		log.Fatal(err)
	}
	date := strings.ReplaceAll(doc.Find(".Results-header").Find(".format-date").Text(), "/", "-")
	results := make([]result, 0, resultsCapacity)
	doc.Find(".Results-table-row").Each(func(i int, s *goquery.Selection) {
		idStr := s.Find(".Results-table-td.Results-table-td--name").Find(".compact").Find("a").AttrOr("href", "")
		id, _ := strconv.Atoi(idStr[strings.LastIndex(idStr, "/")+1:]) // TODO: Handle error.
		clubIdStr := s.Find(".Results-table-club.Results-tablet").Find(".detailed").Find("a").AttrOr("href", "")
		clubId, _ := strconv.Atoi(clubIdStr[strings.LastIndex(clubIdStr, "=")+1:])                   // TODO: Handle error.
		previousPB, _ := s.Find(".Results-table-td.Results-table-td--time").Find(".detailed").Html() // TODO: Handle error.
		previousPB = strings.ReplaceAll(previousPB, "<span class=\"Results-table--normal\">", "")
		previousPB = strings.ReplaceAll(previousPB, "<span class=\"Results-table--red\">", "")
		previousPB = strings.ReplaceAll(previousPB, "<span class=\"Results-table--green\">", "") // TODO: Group using regex?
		previousPB = strings.ReplaceAll(previousPB, "</span>", "")
		previousPB = strings.ReplaceAll(previousPB, "</span>&nbsp;", " ")
		position, _ := strconv.Atoi(s.AttrOr("data-position", ""))           // TODO: Handle error.
		runs, _ := strconv.Atoi(s.AttrOr("data-runs", ""))                   // TODO: Handle error.
		ageGrade, _ := strconv.ParseFloat(s.AttrOr("data-agegrade", ""), 32) // Still made as float64 but convertable to 32. // TODO: Handle error.
		results = append(results, result{
			id:          id,
			clubId:      clubId,
			previousPB:  previousPB,
			name:        s.AttrOr("data-name", ""),
			ageGroup:    s.AttrOr("data-agegroup", ""),
			club:        s.AttrOr("data-club", ""),
			gender:      s.AttrOr("data-gender", ""),
			position:    position,
			runs:        runs,
			ageGrade:    float32(ageGrade),
			achievement: s.AttrOr("data-achievement", ""),
			time:        s.Find(".Results-table-td.Results-table-td--time").Find(".compact").Text(),
		})
	})
	return &event{
		location: location,
		number:   eventNum,
		date:     date,
		results:  results,
	}
}

func main() {
	// Take in a location name:
	fmt.Println("Enter the full Parkrun location name. e.g. 'bathskyline': ")
	var location string
	fmt.Scanln(&location)

	fmt.Printf("Location is %v...\n", location)
	// TODO: Take in a date?
	// TODO: Work out expected event no. from date?
	eventNum := 298
	// TODO: Check which events are missing and iterate through (initial cap of 5?)
	testEvent := getEvent(location, eventNum)
	fmt.Println(testEvent)
	fmt.Println(testEvent.filename())
	testEvent.writeCSV()
}
