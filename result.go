package main

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
	currentPB   string
}
