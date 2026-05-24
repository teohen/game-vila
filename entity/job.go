package entity

import (
	"fmt"

	"github/teohen/mgm-tto/debug"
)

type JobType int

const (
	JobTypeMove JobType = iota
	JobTypeChopTrees
)

type Job struct {
	Type    JobType
	TargetX int
	TargetY int
}

type JobQueue struct {
	jobs []Job
}

func NewJobQueue() JobQueue {
	return JobQueue{}
}

func (q *JobQueue) Push(job Job) {
	q.jobs = append(q.jobs, job)
	q.debugJobs("push", job)
}

func (q *JobQueue) Pop() *Job {
	if len(q.jobs) == 0 {
		return nil
	}
	job := q.jobs[0]
	q.jobs = q.jobs[1:]
	q.debugJobs("pop", job)
	return &job
}

func (q *JobQueue) Get() []Job {
	return q.jobs
}

func (q *JobQueue) debugJobs(action string, job Job) {
	if debug.IsEnabled(debug.Job) {
		fmt.Printf("[DEBUG] JobQueue %s type=%d target=(%d,%d) queue=%d\n",
			action, job.Type, job.TargetX, job.TargetY, len(q.jobs))
	}
}
