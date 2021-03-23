package part3_test

import (
	. "part3"
	"reflect"
	"testing"
	"time"
)

func TestScheduler_should_lazily_take_a_function_and_arguments(t *testing.T) {
	s := NewScheduler()

	f := func(...int) int {
		t.Fatalf("Scheduler is not lazy")
		return 0
	}

	s.Add(f, 1, 2)
}

func TestScheduler_should_return_expected_results_in_scheduled_order(t *testing.T) {
	s := NewScheduler()

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

var results = make([]int, 0)

func benchmarkRun(b *testing.B, scheduleNum int) {
	s := NewScheduler()

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

// We've refactored the benchmarking to make it easier to try out different scenarios
func BenchmarkScheduler_Run100(b *testing.B) {
	// You should see that perf has improved ~100x since part 2
	benchmarkRun(b, 100)
}

func BenchmarkScheduler_Run1m(b *testing.B) {
	benchmarkRun(b, 1000000)
}
