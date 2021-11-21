package main

// resultsCapacity constant is the initial slice capacity to be used for event results.
const resultsCapacity = 500 // TODO: REVIEW capacity of 500.

// result is a runner's result during a single Parkrun event.
// It contains no details of or link to the event itself.
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
	currentPB   string
}
