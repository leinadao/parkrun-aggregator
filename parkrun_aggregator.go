package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// TODO: Create result objects and event objects (and runner?)

func getEvent(location string, eventNum int) []map[string]string {
	// TODO: Return object.
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
	results := make([]map[string]string, 0, 500) // TODO: REVIEW capacity of 500.
	doc.Find(".Results-table-row").Each(func(i int, s *goquery.Selection) {
		result := make(map[string]string)
		result["id"] = s.Find(".Results-table-td.Results-table-td--name").Find(".compact").Find("a").AttrOr("href", "")
		result["id"] = result["id"][strings.LastIndex(result["id"], "/")+1:]
		result["name"] = s.AttrOr("data-name", "")
		result["agegroup"] = s.AttrOr("data-agegroup", "")
		result["club"] = s.AttrOr("data-club", "")
		result["clubId"] = s.Find(".Results-table-club.Results-tablet").Find(".detailed").Find("a").AttrOr("href", "")
		result["clubId"] = result["clubId"][strings.LastIndex(result["clubId"], "=")+1:]
		result["gender"] = s.AttrOr("data-gender", "")
		result["position"] = s.AttrOr("data-position", "")
		result["runs"] = s.AttrOr("data-runs", "")
		result["agegrade"] = s.AttrOr("data-agegrade", "")
		result["achievement"] = s.AttrOr("data-achievement", "")
		result["time"] = s.Find(".Results-table-td.Results-table-td--time").Find(".compact").Text()
		result["previousPB"], _ = s.Find(".Results-table-td.Results-table-td--time").Find(".detailed").Html() // TODO: Handle error.
		result["previousPB"] = strings.ReplaceAll(result["previousPB"], "<span class=\"Results-table--normal\">", "")
		result["previousPB"] = strings.ReplaceAll(result["previousPB"], "</span>", "")
		result["previousPB"] = strings.ReplaceAll(result["previousPB"], "</span>&nbsp;", " ")
		results = append(results, result)
	})
	return results
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
	eventData := getEvent(location, eventNum)
	fmt.Println(eventData)
}
