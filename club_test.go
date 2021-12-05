package main

import "testing"

func TestClubFields(t *testing.T) {
	var i1, n1 = 1, "test"
	var c1 = club{
		id:   i1,
		name: n1,
	}
	c1id := c1.id
	if c1id != i1 {
		t.Errorf("club{id=%d}.id = %d; want %d", i1, c1id, i1)
	}
	c1name := c1.name
	if c1name != n1 {
		t.Errorf("club{name=%s}.name = %s; want %s", n1, c1name, n1)
	}
}

func TestClubEquality(t *testing.T) {
	var i1, i2, n1, n2 = 1, 2, "test", "Testing"
	var c1 = club{
		id:   i1,
		name: n1,
	}
	var c2 = club{
		id:   i1,
		name: n1,
	}
	var c3 = club{
		id:   i2,
		name: n2,
	}
	t1 := c1 == c2
	if !t1 {
		t.Errorf("club %v == club %v = %t; want true", c1, c2, t1)
	}
	t2 := c1 == c3
	if t2 {
		t.Errorf("club %v == club %v = %t; want false", c1, c3, t2)
	}
}
