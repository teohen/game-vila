package entity

type Job struct {
	TargetX int
	TargetY int
}

type JobQueue struct {
	jobs []Job
}

func NewJobQueue() JobQueue {
	return JobQueue{}
}

func (q *JobQueue) Push(targetX, targetY int) {
	q.jobs = append(q.jobs, Job{TargetX: targetX, TargetY: targetY})
}

func (q *JobQueue) Pop() *Job {
	if len(q.jobs) == 0 {
		return nil
	}
	job := q.jobs[0]
	q.jobs = q.jobs[1:]
	return &job
}
