package part1

type work func(...int) int

type Scheduler struct{}

func NewScheduler() *Scheduler {
	return &Scheduler{}
}

func (s *Scheduler) Add(w work, args ...int) {}
