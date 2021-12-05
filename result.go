package main

// TODO: Add test code!

// resultsCapacity constant is the initial slice capacity to be used for event results.
const resultsCapacity = 500 // TODO: REVIEW capacity of 500.

// result is a runner's result during a single Parkrun event.
// It contains no details of or link to the event itself.
type result struct {
	runner      runner
	ageGroup    string // TODO: REVIEW: More likely to change so not moved to runner.
	club        club   // TODO: REVIEW: More likely to change so not moved to runner.
	position    int
	runs        int
	ageGrade    float32
	achievement string
	time        string
	currentPB   string
}
