package part1_test

import (
	. "part1"
	"testing"
)

func TestScheduler_should_lazily_take_a_function_and_arguments(t *testing.T) {
	s := New()

	f := func(...int) int {
		t.Fatalf("Scheduler is not lazy")
		return 0
	}

	s.Add(f, 1, 2)
	s.Add(f, 1, 2)
}
