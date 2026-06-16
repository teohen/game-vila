package resource

import "github/teohen/mgm-tto/cnts"

type ResourceType string

type Resources struct {
	ID           string
	typeResource ResourceType
}

type IResource interface {
	ID() string
	Type() ResourceType
	Draw()
	Pos() cnts.Point
}

const (
	ResourceTreeType ResourceType = "ResourceTreeType"
	ResourceWoodType ResourceType = "ResourceWoodType"
)
