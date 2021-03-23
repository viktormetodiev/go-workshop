package part5_test

import (
	. "part5"
	"reflect"
	"runtime"
	"testing"
	"time"
)

func TestScheduler_should_lazily_take_a_function_and_arguments(t *testing.T) {
	s := NewScheduler(0)

	f := func(...int) int {
		t.Fatalf("Scheduler is not lazy")
		return 0
	}

	s.Add(f, 1, 2)
}

func TestScheduler_should_return_expected_results_in_scheduled_order(t *testing.T) {
	s := NewScheduler(0)

	sum := func(args ...int) (total int) {
		for _, v := range args {
			total += v
		}

		return
	}

	multiply := func(args ...int) (total int) {
		total = 1
		for _, v := range args {
			total *= v
		}

		return
	}

	s.Add(sum, 1, 2, 3)
	s.Add(multiply, 3, 4, 5)

	actual := s.Run()
	expected := []int{6, 60}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Wanted %v, got %v", expected, actual)
	}

}

func TestScheduler_should_clean_up_goroutines(t *testing.T) {
	s := NewScheduler(2)

	sum := func(args ...int) (total int) {
		for _, v := range args {
			total += v
		}

		return
	}

	multiply := func(args ...int) (total int) {
		total = 1
		for _, v := range args {
			total *= v
		}

		return
	}

	// Schedule several functions that are greater than the number of max goroutines we've set (2)
	s.Add(sum, 1, 2, 3)
	s.Add(multiply, 3, 4, 5)
	s.Add(sum, 1, 2, 3)
	s.Add(multiply, 3, 4, 5)
	s.Add(sum, 1, 2, 3)
	s.Add(multiply, 3, 4, 5)

	prevRoutines := runtime.NumGoroutine()

	// Count the number of goroutines before and after the run
	s.Run()
	// Wait for the next tick to ensure the goroutines have been cleaned up
	time.Sleep(1 * time.Millisecond)

	if ng := runtime.NumGoroutine(); ng != prevRoutines {
		t.Errorf("There were %v active goroutines, expected %v", ng, prevRoutines)
	}
}
