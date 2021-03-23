package part3

import "sync"

type work func(...int) int

type job struct {
	w work
	args []int
}

type Scheduler struct{
	jobs []job
}

func NewScheduler() *Scheduler {
	return &Scheduler{}
}

func (s *Scheduler) Add(w work, args ...int) {
	s.jobs = append(s.jobs, job{w, args})
}

func (s *Scheduler) Run() []int {
	totalJobs := len(s.jobs)
	results := make([]int, totalJobs)
	
	wg := sync.WaitGroup{}
	wg.Add(totalJobs)

	for i, j := range s.jobs {
		go func(i int, j job) {
			defer wg.Done()
			results[i] = j.w(j.args...)
		}(i, j)
	}

	wg.Wait()

	return results
}
