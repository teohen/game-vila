package entity

import "github/teohen/mgm-tto/goap"

type JobType int

const (
	JobTypeMove JobType = iota
	JobTypeChopTrees
)

type Job struct {
	Type       JobType
	TargetX    int
	TargetY    int
	TargetID   string
	WorldState *goap.State
}

type JobQueue struct {
	jobs []Job
}

func NewJobQueue() JobQueue {
	return JobQueue{}
}

func (q *JobQueue) Push(job Job) {
	q.jobs = append(q.jobs, job)
}

func (q *JobQueue) Pop() *Job {
	if len(q.jobs) == 0 {
		return nil
	}
	job := q.jobs[0]
	q.jobs = q.jobs[1:]
	return &job
}

func (q *JobQueue) GetJobs() []Job {
	return q.jobs
}

var jobQueue = &JobQueue{}

func GetJobQueue() *JobQueue {
	return jobQueue
}
