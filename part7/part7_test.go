package part7_test

import (
	. "part7"
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
	s := NewScheduler(0, 1000*time.Millisecond)

	sum := func(args ...int) (total int) {
		for _, v := range args {
			total += v
		}

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

	actual := s.Run()
	expected := []Result{Result{0, Timeout}, Result{60, nil}}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Wanted %v, got %v", expected, actual)
	}

}

func TestScheduler_should_manage_multiple_execs(t *testing.T) {
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
	s.Add(sum, 2, 3, 4)

	actual1 := s.Run()
	expected1 := []Result{
		Result{6, nil},
		Result{60, nil},
		Result{9, nil},
	}

	if !reflect.DeepEqual(actual1, expected1) {
		t.Errorf("First run: Wanted %v, got %v", expected1, actual1)
	}

	s.Add(multiply, 3, 4, 5, 6)
	s.Add(sum, 1, 2, 5)
	s.Add(multiply, 3, 4, 6)

	actual2 := s.Run()
	expected2 := []Result{
		Result{360, nil},
		Result{8, nil},
		Result{72, nil},
	}

	if !reflect.DeepEqual(actual2, expected2) {
		t.Errorf("Second run: Wanted %v, got %v", expected2, actual2)
	}
}
