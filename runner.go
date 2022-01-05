package main

// runner is a Parkrun attendee.
type runner struct { // TODO: REVIEW: Add in separate 'Volunteer' or 'Person' or 'Account' for name and id and mix in here?
	id     int
	name   string
	gender string // TODO: REVIEW: Lock this down?
}
