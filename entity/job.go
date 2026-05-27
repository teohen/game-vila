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

func (q *JobQueue) PopClosest(vx, vy int) *Job {
	if len(q.jobs) == 0 {
		return nil
	}
	best := 0
	bestDist := abs(vx-q.jobs[0].TargetX) + abs(vy-q.jobs[0].TargetY)
	for i := 1; i < len(q.jobs); i++ {
		dist := abs(vx-q.jobs[i].TargetX) + abs(vy-q.jobs[i].TargetY)
		if dist < bestDist {
			best = i
			bestDist = dist
		}
	}
	job := q.jobs[best]
	q.jobs = append(q.jobs[:best], q.jobs[best+1:]...)
	q.debugJobs("pop", job)
	return &job
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (q *JobQueue) debugJobs(action string, job Job) {
	if debug.IsEnabled(debug.Job) {
		fmt.Printf("[DEBUG] JobQueue %s type=%d target=(%d,%d) queue=%d\n",
			action, job.Type, job.TargetX, job.TargetY, len(q.jobs))
	}
}
