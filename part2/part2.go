package part2

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
	results := make([]int, len(s.jobs))

	for index, job := range s.jobs {
		results[index] = job.w(job.args...)
	}

	return results
}
