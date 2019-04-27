package todotxt

import (
	"strings"
	"testing"
	"time"
)

func TestPriorityFromLetter(t *testing.T) {
	tests := []struct {
		letter        string
		expectedError bool
		expectedVal   int
	}{
		{"A", false, 0},
		{"B", false, 1},
		{"Z", false, 25},
		{"a", true, 0},
		{"AA", true, 0},
		{"0", true, 0},
		{"!", true, 0},
	}
	for i, test := range tests {
		p, err := PriorityFromLetter(test.letter)
		if test.expectedError && err == nil {
			t.Errorf("expected error on test %d", i+1)
		}
		if !test.expectedError && err != nil {
			t.Errorf("unexpected error on test %d: %v", i+1, err)
		}
		if p != Priority(test.expectedVal) {
			t.Errorf("expected %d, got %d on test %d", test.expectedVal, p, i+1)
		}
	}
}

func TestPriorityToString(t *testing.T) {
	a, _ := PriorityFromLetter("A")
	z, _ := PriorityFromLetter("Z")
	if a.String() != "A" {
		t.Errorf("Expected A, got %s", a.String())
	}
	if z.String() != "Z" {
		t.Errorf("Expected Z, got %s", z.String())
	}
}

func TestItemToString(t *testing.T) {
	cd, _ := time.Parse("2006-01-02", "2019-04-27")
	tests := map[string]struct {
		item     Item
		expected string
	}{
		"simple":           {item: Item{Description: "simple"}, expected: "simple"},
		"complete":         {item: Item{Complete: true, Description: "complete"}, expected: "x complete"},
		"completeWithDate": {item: Item{Complete: true, CreationDate: &cd, Description: "complete"}, expected: "x 2019-04-27 complete"},
		"completeWithTags": {item: Item{Complete: true, Description: "complete", Tags: []Tag{{Key: TagProject, Value: "proj"}, {Key: TagContext, Value: "test"}, {Key: "custom", Value: "tag"}}}, expected: "x complete +proj @test custom:tag"},
	}
	for test, v := range tests {
		t.Run(test, func(t *testing.T) {
			if v.item.String() != v.expected {
				t.Errorf("Expected %s got %s", v.expected, v.item.String())
			}
		})
	}
}

func TestParse(t *testing.T) {
	input := `(A) Call Mom @Phone +Family
(A) Schedule annual checkup +Health
(B) Outline chapter 5 +Novel @Computer
(C) Add cover sheets @Office +TPSReports
2019-04-27 Plan backyard herb garden @Home
2019-05-27 2019-04-27 Pick up milk @GroceryStore
Research self-publishing services +Novel @Computer
x Download Todo.txt mobile app @Phone custom:tag`
	items, err := Parse(input)
	if err != nil {
		t.Error(err)
	}
	for ln, v := range items {
		if strings.Split(input, "\n")[ln] != v.String() {
			t.Logf("line %d expected %s got %s", ln, strings.Split(input, "\n")[ln], v.String())
		}
	}
}
