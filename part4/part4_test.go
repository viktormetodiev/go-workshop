package part4_test

import (
	. "part4"
	"reflect"
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

func TestScheduler_should_use_max_goroutines(t *testing.T) {
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

	s.Add(sum, 1, 2, 3)
	s.Add(multiply, 3, 4, 5)
	s.Add(sum, 1, 2, 3)
	s.Add(multiply, 3, 4, 5)
	s.Add(sum, 1, 2, 3)
	s.Add(multiply, 3, 4, 5)

	actual := s.Run()
	expected := []int{6, 60, 6, 60, 6, 60}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Wanted %v, got %v", expected, actual)
	}

}

var results = make([]int, 0)

func benchmarkRun(b *testing.B, scheduleNum int) {
	s := NewScheduler(100000)

	sum := func(args ...int) (total int) {
		for _, v := range args {
			total += v
		}

		time.Sleep(10 * time.Millisecond)

		return
	}

	for i := 0; i < scheduleNum; i++ {
		s.Add(sum, 1, 2, 3, 4)
	}

	for n := 0; n < b.N; n++ {
		results = s.Run()
	}
}

func BenchmarkScheduler_Run1m(b *testing.B) {
	benchmarkRun(b, 1000000)
}
