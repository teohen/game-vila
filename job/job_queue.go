package job

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
	if len(q.Jobs) > 0 {
		job := q.Jobs[0]
		return &job
	}
	return nil
}

func (q *JobQueue) Remove(jobName, objectId string) {
	for i, job := range q.Jobs {
		if jobName == job.name && objectId == job.Object.ID() {
			q.Jobs = append(q.Jobs[:i], q.Jobs[i+1:]...)
		}
	}
}

func GetJobQueue() *JobQueue {
	return jobQueue
}

var jobQueue = &JobQueue{}
