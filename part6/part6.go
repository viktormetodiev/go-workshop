package part6

import (
	"errors"
	"runtime"
	"time"
)

type work func(...int) int

type job struct {
	w    work
	args []int
}

type Scheduler struct {
	maxThreads int
	timeout    time.Duration
	jobs       []job
}

func NewScheduler(maxThreads int, timeout time.Duration) *Scheduler {
	if maxThreads == 0 {
		maxThreads = runtime.NumCPU()
	}

	return &Scheduler{
		maxThreads: maxThreads,
		timeout:    timeout,
	}
}

func (s *Scheduler) Add(w work, args ...int) {
	s.jobs = append(s.jobs, job{w, args})
}

var Timeout = errors.New("Job timed out")

type jobRequest struct {
	job
	index int
}

type Result struct {
	Value int
	Err   error
}

type jobCompletion struct {
	result Result
	index  int
}

func doWork(workStream chan jobRequest, resultStream chan jobCompletion, timeout time.Duration) {
	for workToDo := range workStream {
		ch := make(chan int)
		
		go func() {
			defer close(ch)
			ch <- workToDo.job.w(workToDo.job.args...)
		}()

		select {
		case result := <-ch:
			resultStream <- jobCompletion{
				Result{
					result,
					nil,
				},
				workToDo.index,
			}
		case <-time.After(timeout):
			resultStream <- jobCompletion{
				Result{
					0,
					Timeout,
				},
				workToDo.index,
			}
		}
	}
}

func (s *Scheduler) Run() []Result {
	totalJobs := len(s.jobs)
	results := make([]Result, totalJobs)

	workStream, resultStream := make(chan jobRequest, s.maxThreads), make(chan jobCompletion, totalJobs)

	defer close(workStream)
	defer close(resultStream)

	for i := 0; i < s.maxThreads; i++ {
		go doWork(workStream, resultStream, s.timeout)
	}

	for index, jobToDo := range s.jobs {
		workStream <- jobRequest{
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
