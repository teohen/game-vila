package job

import (
	"github/teohen/mgm-tto/cnts"
)

type JobType string

const JobChopTreeType JobType = "ChopTreeJob"

type Object interface {
	GetID() string
	Pos() cnts.Point
}

type Job struct {
	ID      string
	name    string
	typeJob JobType
	Object  Object
}

func NewJob(typeJob JobType, obj Object) *Job {
	j := Job{
		name:    string(typeJob),
		typeJob: typeJob,
		Object:  obj,
	}

	return &j
}

func (j *Job) Name() string {
	return j.name
}

func (j *Job) Type() JobType {
	return j.typeJob
}

func (j *Job) GetJob() Object {
	return j.Object
}
