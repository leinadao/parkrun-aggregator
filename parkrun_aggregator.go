package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

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
	}
}

func main() {
	// Take in a location name:
	fmt.Println("Enter the full Parkrun location name. e.g. 'bathskyline': ")
	var location string
	fmt.Scanln(&location)

	fmt.Printf("Location is %v...\n", location)
	eventNum := 298
	// TODO: Check which events are missing and iterate through (initial cap of 5?)
	testEvent := getEvent(location, eventNum)
	fmt.Println(testEvent)
	fmt.Println(testEvent.filename())
	testEvent.writeCSV()
}
