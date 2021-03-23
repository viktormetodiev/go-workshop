package part5

import (
	"runtime"
)

type work func(...int) int

type job struct {
	w    work
	args []int
}

type Scheduler struct {
	maxThreads int
	jobs       []job
}

func NewScheduler(maxThreads int) *Scheduler {
	if maxThreads == 0 {
		maxThreads = runtime.NumCPU()
	}

	return &Scheduler{
		maxThreads: maxThreads,
	}
}

func (s *Scheduler) Add(w work, args ...int) {
	s.jobs = append(s.jobs, job{w, args})
}

type workItem struct {
	job
	index int
}

type jobResult struct {
	result int
	index  int
}

func doWork(workStream chan workItem, resultStream chan jobResult, cancel chan struct{}) {
	for workToDo := range workStream {
		resultStream <- jobResult{
			workToDo.job.w(workToDo.job.args...),
			workToDo.index,
		}
	}
}

func (s *Scheduler) Run() []int {
	totalJobs := len(s.jobs)
	results := make([]int, totalJobs)

	workStream, resultStream := make(chan workItem, s.maxThreads), make(chan jobResult, totalJobs)
	cancel := make(chan struct{})

	defer close(workStream)
	defer close(resultStream)

	for i := 0; i < s.maxThreads; i++ {
		go doWork(workStream, resultStream, cancel)
	}

	for index, jobToDo := range s.jobs {
		workStream <- workItem{
			jobToDo,
			index,
		}
	}

	for i := 0; i < totalJobs; i++ {
		jobResult := <-resultStream
		results[jobResult.index] = jobResult.result
	}

	return results
}
