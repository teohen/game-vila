package job

import (
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/goap"
)

type JobType string

const JobChopTreeType JobType = "ChopTreeJob"

type Job struct {
	Type       JobType
	TargetPos  cnts.Point
	TargetID   string
	WorldState *goap.State
}

type JobQueue struct {
	Jobs []Job
}

func (q *JobQueue) Push(job Job) {
	q.Jobs = append(q.Jobs, job)
}

func (q *JobQueue) Pop() *Job {
	if len(q.Jobs) == 0 {
		return nil
	}
	job := q.Jobs[0]
	q.Jobs = q.Jobs[1:]
	return &job
}

func (q *JobQueue) Peek() *Job {
	job := q.Jobs[0]
	return &job
}

func (q *JobQueue) Remove(j *Job) {

}

func GetJobQueue() *JobQueue {
	return jobQueue
}

var jobQueue = &JobQueue{}
