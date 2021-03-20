package part2_test

import (
	. "part2"
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

// We capture the results globally to avoid any cache optimisations
var results = make([]int, 0)

// Much like the `Test` convention, there is also the `Benchmark` convention. If we were to run `go test -bench=.` in
// our terminal Go will pass a variable `N` into this function numerous times until it gets a relatively consistent result.
// We can use benchmarking for performance critical pieces of our code. In this case we'll use it to show the difference
// between serialising our work and running it concurrently
func BenchmarkScheduler_Run(b *testing.B) {
	s := NewScheduler()

	sum := func(args ...int) (total int) {
		for _, v := range args {
			total += v
		}

		// pretend to do some work here
		time.Sleep(10 * time.Millisecond)

		return
	}

	// Schedule 100 functions
	for i := 0; i < 100; i++ {
		s.Add(sum, 1, 2, 3, 4)
	}

	// Keep running until we get a consistent result
	for n := 0; n < b.N; n++ {
		// we assign the result to avoid cache optimisations here too
		results = s.Run()
	}
}
