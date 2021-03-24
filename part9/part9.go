package part9

import (
	"context"
	"errors"
	"fmt"
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

var Timeout = errors.New("job timed out")

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
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		ch := make(chan Result)
		
		go func() {
			defer cancel()
			defer close(ch)
			defer func() {
				if err := recover(); err != nil {
					ch <- Result{
						0,
						fmt.Errorf("Panicked: %s", err.(string)),
					}
				}
			}()

			ch <- Result{
				workToDo.job.w(workToDo.job.args...),
				nil,
			}
		}()

		select {
		case result := <-ch:
			resultStream <- jobCompletion{
				result,
				workToDo.index,
			}
		case <-ctx.Done():
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
	jobs, totalJobs := make([]job, len(s.jobs)), len(s.jobs)

	copy(jobs, s.jobs)

	s.jobs = []job{}
	results := make([]Result, totalJobs)

	workStream, resultStream := make(chan jobRequest, s.maxThreads), make(chan jobCompletion, totalJobs)
	defer close(workStream)
	defer close(resultStream)

	for i := 0; i < s.maxThreads; i++ {
		go doWork(workStream, resultStream, s.timeout)
	}

	for index, jobToDo := range jobs {
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
