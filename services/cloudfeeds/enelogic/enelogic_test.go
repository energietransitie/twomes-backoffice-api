package enelogic

import (
	"testing"
	"time"
)

func TestSplitDays_multipleDays(t *testing.T) {
	from, err := time.Parse(time.RFC3339, "2023-12-15T12:57:24Z")
	if err != nil {
		t.Fatal(err)
	}

	to, err := time.Parse(time.RFC3339, "2023-12-20T11:48:16Z")
	if err != nil {
		t.Fatal(err)
	}

	days := splitDays(from, to)

	t.Log(days)

	if len(days) != 5 {
		t.Fatalf("splitDays(from=2023-12-15, to=2023-12-20) len = %d; want 5", len(days))
	}

	got := RequestTime{days[0].Start}
	want := "2023-12-15"
	if got.String() != want {
		t.Errorf("splitDays(from=2023-12-15, to=2023-12-20) [0].Start = %s; want %s", got, want)
	}
	got = RequestTime{days[0].End}
	want = "2023-12-16"
	if got.String() != want {
		t.Errorf("splitDays(from=2023-12-15, to=2023-12-20) [0].End = %s; want %s", got, want)
	}

	got = RequestTime{days[1].Start}
	want = "2023-12-16"
	if got.String() != want {
		t.Errorf("splitDays(from=2023-12-15, to=2023-12-20) [1].Start = %s; want %s", got, want)
	}
	got = RequestTime{days[1].End}
	want = "2023-12-17"
	if got.String() != want {
		t.Errorf("splitDays(from=2023-12-15, to=2023-12-20) [1].End = %s; want %s", got, want)
	}

	got = RequestTime{days[2].Start}
	want = "2023-12-17"
	if got.String() != want {
		t.Errorf("splitDays(from=2023-12-15, to=2023-12-20) [2].Start = %s; want %s", got, want)
	}
	got = RequestTime{days[2].End}
	want = "2023-12-18"
	if got.String() != want {
		t.Errorf("splitDays(from=2023-12-15, to=2023-12-20) [2].End = %s; want %s", got, want)
	}

	got = RequestTime{days[3].Start}
	want = "2023-12-18"
	if got.String() != want {
		t.Errorf("splitDays(from=2023-12-15, to=2023-12-20) [3].Start = %s; want %s", got, want)
	}
	got = RequestTime{days[3].End}
	want = "2023-12-19"
	if got.String() != want {
		t.Errorf("splitDays(from=2023-12-15, to=2023-12-20) [3].End = %s; want %s", got, want)
	}

	got = RequestTime{days[4].Start}
	want = "2023-12-19"
	if got.String() != want {
		t.Errorf("splitDays(from=2023-12-15, to=2023-12-20) [4].Start = %s; want %s", got, want)
	}
	got = RequestTime{days[4].End}
	want = "2023-12-20"
	if got.String() != want {
		t.Errorf("splitDays(from=2023-12-15, to=2023-12-20) [4].End = %s; want %s", got, want)
	}
}

func TestSplitDays_singeDay(t *testing.T) {
	from, err := time.Parse(time.RFC3339, "2023-12-15T04:01:24Z")
	if err != nil {
		t.Fatal(err)
	}

	to, err := time.Parse(time.RFC3339, "2023-12-16T04:02:16Z")
	if err != nil {
		t.Fatal(err)
	}

	days := splitDays(from, to)

	t.Log(days)

	if len(days) != 1 {
		t.Fatalf("splitDays(from=2023-12-15, to=2023-12-16) len = %d; want 1", len(days))
	}

	got := RequestTime{days[0].Start}
	want := "2023-12-15"
	if got.String() != want {
		t.Errorf("splitDays(from=2023-12-15, to=2023-12-16) Start = %s; want %s", got, want)
	}

	got = RequestTime{days[0].End}
	want = "2023-12-16"
	if got.String() != want {
		t.Errorf("splitDays(from=2023-12-15, to=2023-12-16) End = %s; want %s", got, want)
	}
}
