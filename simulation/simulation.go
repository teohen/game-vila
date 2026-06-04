package simulation

import (
	"fmt"
	"time"

	"github/teohen/mgm-tto/constants"
	"github/teohen/mgm-tto/entity"
	"github/teohen/mgm-tto/goap"
	"github/teohen/mgm-tto/save"
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

func NewFromSave(s save.Save) *Simulation {
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
		q.Push(entity.Job{Type: entity.JobType(js.Type), TargetX: js.TargetX, TargetY: js.TargetY})
	}

	return &Simulation{
		world:     w,
		villagers: villagers,
		trees:     trees,
	}
}

func New() *Simulation {
	seed := time.Now().UnixNano()
	w := world.NewWorld(constants.GridRows, constants.GridCols)
	w.Generate(seed)

	// state := goap.StateOf("near_tree=0")

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
		v.Tick(all)
	}
	for _, t := range s.trees {
		t.Tick(nil)
	}

	s.tickCount++
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

func (s *Simulation) RemoveTree(x, y int) bool {
	for i, t := range s.trees {
		if t.X == x && t.Y == y {
			s.world.Vacate(x, y)
			s.trees = append(s.trees[:i], s.trees[i+1:]...)
			return true
		}
	}
	return false
}

func (s *Simulation) PushJob(job entity.Job) {
	entity.GetJobQueue().Push(job)
}

func (s *Simulation) ProcessAxeSelection(cells [][2]int) {
	for _, cell := range cells {
		col, row := cell[0], cell[1]
		tree := s.TreeAt(col, row)
		if tree == nil {
			continue
		}
		s.PushJob(entity.Job{
			Type:       entity.JobTypeChopTrees,
			TargetX:    tree.X,
			TargetY:    tree.Y,
			TargetID:   tree.ID,
			WorldState: s.WorldState,
		})
	}
}

func (s *Simulation) TreeAt(x, y int) *entity.Tree {
	for _, t := range s.trees {
		if t.X == x && t.Y == y {
			return t
		}
	}
	return nil
}

func (s *Simulation) World() *world.World {
	return s.world
}

func (s *Simulation) ToSave() save.Save {
	cells := make([][]int, s.world.Rows())
	for r := range cells {
		cells[r] = make([]int, s.world.Cols())
		for c := range cells[r] {
			cells[r][c] = int(s.world.GetCell(c, r).Type)
		}
	}

	villagers := make([]save.VillagerSave, len(s.villagers))
	for i, v := range s.villagers {
		vs := save.VillagerSave{
			ID:   v.ID,
			Name: v.Name(),
			Type: int(v.Type),
			X:    v.X,
			Y:    v.Y,
		}
		if v.State != entity.StateIdle {
			vs.TargetX = &v.TargetX
			vs.TargetY = &v.TargetY
			vs.State = v.State.String()
		}
		villagers[i] = vs
	}

	var trees []save.TreeSave
	for _, t := range s.trees {
		if t.Health <= 0 {
			continue
		}
		trees = append(trees, save.TreeSave{
			ID:        t.ID,
			X:         t.X,
			Y:         t.Y,
			Health:    t.Health,
			WoodYield: t.WoodYield,
		})
	}

	jobs := make([]save.JobSave, len(entity.GetJobQueue().GetJobs()))
	for i, j := range entity.GetJobQueue().GetJobs() {
		jobs[i] = save.JobSave{
			Type:    int(j.Type),
			TargetX: j.TargetX,
			TargetY: j.TargetY,
		}
	}

	return save.Save{
		Version:   1,
		World:     save.WorldSave{Rows: s.world.Rows(), Cols: s.world.Cols(), Cells: cells},
		Villagers: villagers,
		Trees:     trees,
		Jobs:      jobs,
	}
}

func (s *Simulation) QueueJobs() []entity.Job {
	return entity.GetJobQueue().GetJobs()
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
