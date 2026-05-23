package game

import (
	"fmt"
	"math/rand"
	"time"

	"github/teohen/mgm-tto/constants"
	"github/teohen/mgm-tto/entity"
	"github/teohen/mgm-tto/save"
	"github/teohen/mgm-tto/world"
)

type Simulation struct {
	tickCount int
	world     *world.World
	villagers []*entity.Villager
	trees     []*entity.Tree
	jobQueue  entity.JobQueue
}

func NewSimulationFromSave(s save.Save) *Simulation {
	cells := make([][]world.CellType, s.World.Rows)
	for r := range cells {
		cells[r] = make([]world.CellType, s.World.Cols)
		for c := range cells[r] {
			cells[r][c] = world.CellType(s.World.Cells[r][c])
		}
	}
	w := world.NewWorldFromCells(cells)

	villagers := make([]*entity.Villager, len(s.Villagers))
	for i, vs := range s.Villagers {
		v := entity.NewVillager(vs.ID, vs.Name, vs.X, vs.Y)
		w.Occupy(vs.X, vs.Y)
		if vs.TargetX != nil && vs.TargetY != nil {
			v.SetTarget(*vs.TargetX, *vs.TargetY, w)
		}
		villagers[i] = v
	}

	trees := make([]*entity.Tree, len(s.Trees))
	for i, ts := range s.Trees {
		t := entity.NewTree(ts.ID, ts.X, ts.Y, ts.Health, ts.WoodYield)
		w.Occupy(ts.X, ts.Y)
		trees[i] = t
	}

	q := entity.NewJobQueue()
	for _, js := range s.Jobs {
		q.Push(js.TargetX, js.TargetY)
	}

	return &Simulation{
		world:     w,
		villagers: villagers,
		trees:     trees,
		jobQueue:  q,
	}
}

func (s *Simulation) Tick() {
	for _, v := range s.villagers {
		event := v.Tick(s.world)
		switch event {
		case entity.EventIdle, entity.EventArrived:
			job := s.jobQueue.Pop()
			if job == nil {
				job = s.randomTarget(v)
			}
			v.SetTarget(job.TargetX, job.TargetY, s.world)
		}
	}
	for _, t := range s.trees {
		t.Tick(s.world)
	}
	s.tickCount++
}

func (s *Simulation) randomTarget(v *entity.Villager) *entity.Job {
	cx, cy := v.Pos()
	for {
		x := rand.Intn(s.world.Cols())
		y := rand.Intn(s.world.Rows())
		if x == cx && y == cy {
			continue
		}
		if !s.world.IsWalkable(x, y) {
			continue
		}
		return &entity.Job{TargetX: x, TargetY: y}
	}
}

func (s *Simulation) AdvanceTicks(n int) {
	for i := 0; i < n; i++ {
		s.Tick()
	}
}

func (s *Simulation) SetTarget(villagerID string, x, y int) {
	for _, v := range s.villagers {
		if v.ID == villagerID {
			v.SetTarget(x, y, s.world)
			return
		}
	}
}

func (s *Simulation) Pos(entityID string) (int, int) {
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
	return -1, -1
}

func (s *Simulation) TickCount() int {
	return s.tickCount
}

func (s *Simulation) AddVillager(v *entity.Villager) {
	s.villagers = append(s.villagers, v)
	s.world.Occupy(v.X, v.Y)
}

func (s *Simulation) AddTree(tree *entity.Tree) {
	s.trees = append(s.trees, tree)
	s.world.Occupy(tree.X, tree.Y)
}

func (s *Simulation) World() *world.World {
	return s.world
}

const (
	forestFrequency = 0.07
	forestThreshold = 0.1
	treeHealth      = 3
	treeWoodYield   = 5
)

func NewSimulationDefault() *Simulation {
	seed := time.Now().UnixNano()
	w := world.NewWorld(constants.GridRows, constants.GridCols)
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
			id := fmt.Sprintf("tree-%d", treeCount)
			t := entity.NewTree(id, c, r, treeHealth, treeWoodYield)
			w.Occupy(c, r)
			trees = append(trees, t)
		}
	}

	return &Simulation{
		world:     &w,
		villagers: nil,
		trees:     trees,
		jobQueue:  entity.NewJobQueue(),
	}
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
