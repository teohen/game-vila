package entity

type Job struct {
	TargetX int
	TargetY int
}

type JobQueue struct {
	jobs []Job
}

func NewJobQueue() JobQueue {
	return JobQueue{
		jobs: []Job{
			{TargetX: 30, TargetY: 11},
		},
	}
}

func (q *JobQueue) Pop() *Job {
	if len(q.jobs) == 0 {
		return nil
	}
	job := q.jobs[0]
	q.jobs = q.jobs[1:]
	return &job
}
