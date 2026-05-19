package scence

import (
	"github/teohen/mgm-tto/grid"
)

type Scene struct {
	grid grid.Grid
}

func New() *Scene {
	scene := Scene{}
	return &scene
}
