package main

import "testing"

func TestGetDayID(t *testing.T) {
	t.Run("on first day", func(t *testing.T) {
		var timestamp int64 = 1451606400
		got := GetDayID(timestamp)
		want := 0

		if got != want {
			t.Errorf("ParseDay: got %d but expected %d", got, want)
		}
	})

	t.Run("on last day", func(t *testing.T) {
		var timestamp int64 = 1452755160
		got := GetDayID(timestamp)
		want := 13

		if got != want {
			t.Errorf("ParseDay: got %d but expected %d", got, want)
		}
	})

	t.Run("before first day", func(t *testing.T) {
		got := GetDayID(0)
		want := -1

		if got != want {
			t.Errorf("ParseDay: got %d but expected %d", got, want)
		}
	})
}

func TestGetDisjointUsers(t *testing.T) {
	t.Run("no overlap", func(t *testing.T) {
		users1 := map[int64]bool{
			1: true,
			2: true,
			3: true,
		}

		users2 := map[int64]bool{
			4: true,
			5: true,
			6: true,
		}

		got := GetDisjointUsers(users1, users2)

		if len(got) != 3 {
			t.Errorf("Expected overlapping set of length 0 but set has length %d", len(got))
		}
	})

	t.Run("overlap", func(t *testing.T) {
		users1 := map[int64]bool{
			1: true,
			2: true,
			3: true,
		}

		users2 := map[int64]bool{
			2: true,
			3: true,
			4: true,
		}

		got := GetDisjointUsers(users1, users2)

		if len(got) != 1 {
			t.Errorf("Expected overlapping set of length 2 but set has length %d", len(got))
		}
	})
}

func TestComposeLine(t *testing.T) {
	var activities [NumDays]usersSet
	for i := 0; i < NumDays; i++ {
		activities[i] = usersSet{}
	}

	activities[0] = usersSet{
		1: true,
		2: true,
		3: true,
	}
	activities[1] = usersSet{
		1: true,
	}
	activities[2] = usersSet{
		2: true,
		3: true,
	}
	activities[5] = usersSet{
		1: true,
		2: true,
		3: true,
		4: true,
		5: true,
	}
	activities[6] = usersSet{
		2: true,
		3: true,
		4: true,
		5: true,
	}
	activities[13] = usersSet{
		1: true,
		2: true,
		3: true,
	}

	t.Run("day one", func(t *testing.T) {
		got := ComposeLine(0, activities)
		want := "1,2,1,0,0,0,0,0,0,0,0,0,0,0,0"

		if got != want {
			t.Errorf("got '%s' but expected '%s'", got, want)
		}
	})

	t.Run("day two", func(t *testing.T) {
		got := ComposeLine(1, activities)
		want := "2,0,0,0,0,0,0,0,0,0,0,0,0,0,0"

		if got != want {
			t.Errorf("got '%s' but expected '%s'", got, want)
		}
	})

	t.Run("day six", func(t *testing.T) {
		got := ComposeLine(5, activities)
		want := "6,1,4,0,0,0,0,0,0,0,0,0,0,0,0"

		if got != want {
			t.Errorf("got '%s' but expected '%s'", got, want)
		}
	})

	t.Run("last day", func(t *testing.T) {
		got := ComposeLine(13, activities)
		want := "14,3,0,0,0,0,0,0,0,0,0,0,0,0,0"

		if got != want {
			t.Errorf("got '%s' but expected '%s'", got, want)
		}
	})
}
