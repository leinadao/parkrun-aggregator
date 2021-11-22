package main

// clubsCapacity constant is the initial slice capacity to be used for event clubs.
const clubsCapacity = 500 // TODO: REVIEW capacity of 500.

// club is a running club registered against a Parkrun attendee.
type club struct {
	id   int
	name string
}
