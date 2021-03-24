package part6_test

import (
	. "part6"
	"reflect"
	"runtime"
	"testing"
	"time"
)

func TestScheduler_should_lazily_take_a_function_and_arguments(t *testing.T) {
	s := NewScheduler(0, 1000*time.Millisecond)

	f := func(...int) int {
		t.Fatalf("Scheduler is not lazy")
		return 0
	}

	s.Add(f, 1, 2)
}

func TestScheduler_should_return_expected_results_in_scheduled_order(t *testing.T) {
	s := NewScheduler(0, 1000*time.Millisecond)

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
	expected := []Result{Result{6, nil}, Result{60, nil}}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Wanted %v, got %v", expected, actual)
	}

}

func TestScheduler_should_clean_up_goroutines(t *testing.T) {
	s := NewScheduler(2, 1000*time.Millisecond)

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
	s.Add(sum, 1, 2, 3)
	s.Add(multiply, 3, 4, 5)
	s.Add(sum, 1, 2, 3)
	s.Add(multiply, 3, 4, 5)

	s.Run()
	time.Sleep(1 * time.Millisecond)

	if ng := runtime.NumGoroutine(); ng != 2 {
		t.Errorf("There were %v active goroutines, expected 2", ng)
	}
}

func TestScheduler_should_timeout_long_running_funcs(t *testing.T) {
	// We add another param to our constructor which allows us to set a maximum time for each piece of work
	s := NewScheduler(0, 1*time.Second)

	sum := func(args ...int) (total int) {
		for _, v := range args {
			total += v
		}

		// This piece of work should exceed our timeout
		time.Sleep(2 * time.Second)

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

	start := time.Now()
	actual := s.Run()
	// We now return a `Result` data structure that details the value and error.
	// If we have a Timeout it will match the exported error type with a zero value.
	expected := []Result{
		Result{0, Timeout},
		Result{60, nil},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Wanted %v, got %v", expected, actual)
	}

	if time.Since(start) > (2 * time.Second) {
		t.Errorf("The work did not correctly time out")
	}

}
