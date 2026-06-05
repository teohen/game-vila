package entity

import (
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/goap"
)

type JobType int

const (
	JobTypeMove JobType = iota
	JobTypeChopTrees
)

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

var jobQueue = &JobQueue{}

func GetJobQueue() *JobQueue {
	return jobQueue
}
