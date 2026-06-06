package simulation

import (
	"fmt"
	"time"

	"github/teohen/mgm-tto/buildings"
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/entity"
	"github/teohen/mgm-tto/events"
	"github/teohen/mgm-tto/goap"
	"github/teohen/mgm-tto/job"
	"github/teohen/mgm-tto/world"
)

type Tool int

const (
	ToolSelect Tool = iota
	ToolAxe
)

type Simulation struct {
	tickCount int
	world     *world.World
	villagers []*entity.Villager
	trees     []*entity.Tree
	buildings []*buildings.Storage

	ActiveTool Tool
	Selected   map[[2]int]bool
	WorldState *goap.State
	Goals      []*goap.State
}

const (
	forestFrequency = 0.07
	forestThreshold = 0.1
	treeHealth      = 100
	treeWoodYield   = 20
)

func New() *Simulation {
	seed := time.Now().UnixNano()
	w := world.NewWorld(cnts.GridRows, cnts.GridCols)
	w.Generate(seed)

	forestNoise := world.NewNoise(seed + 1)
	var trees []*entity.Tree
	treeCount := 0
	for r := 0; r < w.Rows(); r++ {
		for c := 0; c < w.Cols(); c++ {
			cell := w.GetCell(c, r)
			if cell.Type != world.Grass {
				continue
			}
			if forestNoise.Noise2D(float64(c)*forestFrequency, float64(r)*forestFrequency) < forestThreshold {
				continue
			}
			treeCount++
			id := fmt.Sprintf("tree_%d", treeCount)

			t := entity.NewTree(id, c, r, treeHealth, treeWoodYield)
			w.Occupy(c, r)
			trees = append(trees, t)
		}
	}

	return &Simulation{
		world:     &w,
		villagers: nil,
		trees:     trees,
	}
}

func (s *Simulation) Tick() {
	all := s.Entities()
	for _, v := range s.villagers {
		v.Tick(&all, s.world)
	}

	s.processEvents()
	s.tickCount++
}

func (s *Simulation) GetEntityPosition(entityID string) cnts.Point {
	for _, v := range s.villagers {
		if v.ID == entityID {
			return v.Pos()
		}
	}
	for _, t := range s.trees {
		if t.ID == entityID {
			return t.Pos()
		}
	}
	return cnts.Point{X: -1, Y: -1}
}

func (s *Simulation) AddVillager(v *entity.Villager) {
	s.villagers = append(s.villagers, v)
	// s.world.Occupy(v.Pos().X, v.Pos().Y)
}

func (s *Simulation) AddTree(tree *entity.Tree) {
	s.trees = append(s.trees, tree)
	s.world.Occupy(tree.Pos().X, tree.Pos().Y)
}

func (s *Simulation) RemoveTree(x, y int) bool {
	p := cnts.Point{
		X: x,
		Y: y,
	}

	for i, t := range s.trees {
		if t.Pos() == p {
			s.world.Vacate(x, y)
			s.trees = append(s.trees[:i], s.trees[i+1:]...)
			return true
		}
	}
	return false
}

func (s *Simulation) AddStorage(storage *buildings.Storage) {
	s.buildings = append(s.buildings, storage)
	s.World().Occupy(storage.Pos().X, storage.Pos().Y)
}

func (s *Simulation) PushJob(j job.Job) {
	job.GetJobQueue().Push(j)
}

func (s *Simulation) ProcessAxeSelection(cells [][2]int) {
	for _, cell := range cells {
		col, row := cell[0], cell[1]
		tree := s.TreeAt(col, row)
		if tree == nil {
			continue
		}
		s.PushJob(job.Job{
			Type:       job.JobChopTreeType,
			TargetPos:  tree.Pos(),
			TargetID:   tree.ID,
			WorldState: s.WorldState,
		})
	}
}

func (s *Simulation) TreeAt(x, y int) *entity.Tree {
	p := cnts.Point{
		X: x, Y: y,
	}
	for _, t := range s.trees {
		if t.Pos() == p {
			return t
		}
	}
	return nil
}

func (s *Simulation) World() *world.World {
	return s.world
}

func (s *Simulation) Entities() []entity.Entity {
	total := len(s.villagers) + len(s.trees)
	all := make([]entity.Entity, 0, total)
	for _, v := range s.villagers {
		all = append(all, v)
	}
	for _, t := range s.trees {
		all = append(all, t)
	}
	return all
}

func (s *Simulation) Buildings() []*buildings.Storage {
	return s.buildings
}

func (s *Simulation) OnSelectionComplete() {
	switch s.ActiveTool {
	case ToolAxe:
		cells := make([][2]int, 0, len(s.Selected))
		for pos := range s.Selected {
			cells = append(cells, pos)
		}
		s.ProcessAxeSelection(cells)
	}
}

func (s *Simulation) processEvents() {
	for {
		select {
		case evt := <-events.EventQueue:
			switch evt.Type {
			case events.EventTreeCut:
				treePos := evt.Payload["pos"].(cnts.Point)
				s.RemoveTree(treePos.X, treePos.Y)
			default:
				// Canal está vazio, sai do loop de eventos e segue o frame
				return
			}
		default:
			// Canal está vazio, sai do loop de eventos e segue o frame
			return
		}
	}
}
