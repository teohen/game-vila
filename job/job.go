package job

import (
	"github/teohen/mgm-tto/cnts"
	"strings"

	"github.com/google/uuid"
)

type JobType string

const JobChopTreeType JobType = "ChopTreeJob"

type Object interface {
	ID() string
	Pos() cnts.Point
}

type Job struct {
	id      string
	name    string
	typeJob JobType
	Object  Object
}

func NewJob(typeJob JobType, obj Object) *Job {
	j := Job{
		id:      strings.ReplaceAll(uuid.NewString(), "-", ""),
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

func (j *Job) GetObject() Object {
	return j.Object
}

func (j *Job) ID() string {
	return j.id
}
