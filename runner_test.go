package main

import "testing"

func TestRunnerFields(t *testing.T) {
	var i1, n1, g1 = 1, "test", "Male"
	var r1 = runner{
		id:     i1,
		name:   n1,
		gender: g1,
	}
	r1id := r1.id
	if r1id != i1 {
		t.Errorf("runner{id=%d}.id = %d; want %d", i1, r1id, i1)
	}
	r1name := r1.name
	if r1name != n1 {
		t.Errorf("runner{name=%s}.name = %s; want %s", n1, r1name, n1)
	}
	r1gender := r1.gender
	if r1gender != g1 {
		t.Errorf("runner{gender=%s}.gender = %s; want %s", g1, r1gender, g1)
	}
}

func TestRunnerEquality(t *testing.T) {
	var i1, i2, n1, n2, g1, g2 = 1, 2, "test", "Testing", "Male", "Female"
	var r1 = runner{
		id:     i1,
		name:   n1,
		gender: g1,
	}
	var r2 = runner{
		id:     i1,
		name:   n1,
		gender: g1,
	}
	var r3 = runner{
		id:     i2,
		name:   n2,
		gender: g2,
	}
	t1 := r1 == r2
	if !t1 {
		t.Errorf("runner %v == runner %v = %t; want true", r1, r2, t1)
	}
	t2 := r1 == r3
	if t2 {
		t.Errorf("runner %v == runner %v = %t; want false", r1, r3, t2)
	}
}
