package simulation

import (
	"fmt"
	"math/rand"
	"time"

	"github/teohen/mgm-tto/constants"
	"github/teohen/mgm-tto/debug"
	"github/teohen/mgm-tto/entity"
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
	deposits  []*entity.Deposit
	jobQueue  entity.JobQueue

	ActiveTool Tool
	Selected   map[[2]int]bool

	planner *entity.GoapPlanner
}

const (
	forestFrequency = 0.07
	forestThreshold = 0.1
	treeHealth      = 3
	treeWoodYield   = 5
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
		v.Wood = vs.Wood
		if vs.MaxCarryWeight > 0 {
			v.MaxCarryWeight = vs.MaxCarryWeight
		}
		w.Occupy(vs.X, vs.Y)
		if vs.TargetX != nil && vs.TargetY != nil {
			plan := []entity.PlanStep{
				{Trait: entity.TraitMove, TargetX: *vs.TargetX, TargetY: *vs.TargetY},
			}
			v.SetPlan(plan)
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

	deposits := make([]*entity.Deposit, len(s.Deposits))
	for i, ds := range s.Deposits {
		deposits[i] = entity.NewDeposit(ds.X, ds.Y)
	}

	return &Simulation{
		world:     w,
		villagers: villagers,
		trees:     trees,
		deposits:  deposits,
		jobQueue:  q,
		planner:   entity.NewGoapPlanner(),
	}
}

func New() *Simulation {
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

	var deposits []*entity.Deposit
	rng := rand.New(rand.NewSource(seed + 2))
	for attempt := 0; attempt < 100; attempt++ {
		dx := rng.Intn(constants.GridCols)
		dy := rng.Intn(constants.GridRows)
		cell := w.GetCell(dx, dy)
		if cell != nil && cell.Walkable() && !w.IsOccupied(dx, dy) {
			deposits = append(deposits, entity.NewDeposit(dx, dy))
			break
		}
	}

	return &Simulation{
		world:     &w,
		villagers: nil,
		trees:     trees,
		deposits:  deposits,
		jobQueue:  entity.NewJobQueue(),
		planner:   entity.NewGoapPlanner(),
	}
}

func (s *Simulation) Tick() {
	for _, v := range s.villagers {
		event := v.Tick(s.world)
		if event == entity.EventIdle {
			job := s.jobQueue.PopClosest(v.X, v.Y)
			if job != nil {
				s.setPlan(v, *job)
			}
		}
	}

	if len(s.deposits) > 0 {
		for _, v := range s.villagers {
			if v.IsIdle() && v.IsCarrying() {
				s.setDepositPlan(v)
			}
		}
	}

	for _, t := range s.trees {
		t.Tick(s.world)
	}

	s.removeDeadTrees()

	s.debugSimulation()

	s.tickCount++
}

func treeHealthKey(x, y int) string {
	return fmt.Sprintf("tree_health_%d_%d", x, y)
}

func (s *Simulation) buildState(v *entity.Villager) entity.GoapState {
	state := entity.GoapState{
		"at_x":        v.X,
		"at_y":        v.Y,
		"wood":        v.Wood,
		"weight":      v.CurrentWeight(),
		"max_weight":  v.MaxCarryWeight,
	}
	for _, t := range s.trees {
		if t.Health > 0 {
			state[treeHealthKey(t.X, t.Y)] = t.Health
		}
	}
	if len(s.deposits) > 0 {
		state["deposit_exists"] = 1
	} else {
		state["deposit_exists"] = 0
	}
	return state
}

func (s *Simulation) jobToGoal(job entity.Job) entity.GoapState {
	switch job.Type {
	case entity.JobTypeMove:
		return entity.GoapState{"at_x": job.TargetX, "at_y": job.TargetY}
	case entity.JobTypeChopTrees:
		return entity.GoapState{treeHealthKey(job.TargetX, job.TargetY): 0}
	default:
		return nil
	}
}

func (s *Simulation) buildActions(v *entity.Villager, job entity.Job) []entity.GoapAction {
	var actions []entity.GoapAction

	switch job.Type {
	case entity.JobTypeChopTrees:
		tree := s.TreeAt(job.TargetX, job.TargetY)
		actions = append(actions, s.moveToAction(job.TargetX, job.TargetY))
		if tree != nil {
			actions = append(actions, s.chopAction(job.TargetX, job.TargetY, tree))
		}
	case entity.JobTypeMove:
		actions = append(actions, s.moveToAction(job.TargetX, job.TargetY))
	}

	for _, d := range s.deposits {
		actions = append(actions, s.moveToAction(d.X, d.Y))
		actions = append(actions, s.depositWoodAction(d.X, d.Y))
	}

	return actions
}

func (s *Simulation) moveToAction(x, y int) entity.GoapAction {
	return entity.GoapAction{
		Name: fmt.Sprintf("MoveTo(%d,%d)", x, y),
		Check: func(_ entity.GoapState) bool { return true },
		Apply: func(state entity.GoapState) entity.GoapState {
			state["at_x"] = x
			state["at_y"] = y
			return state
		},
		Cost: func(state entity.GoapState) int {
			dx := state["at_x"] - x
			dy := state["at_y"] - y
			if dx < 0 {
				dx = -dx
			}
			if dy < 0 {
				dy = -dy
			}
			return dx + dy
		},
		BuildStep: func() entity.PlanStep {
			return entity.PlanStep{Trait: entity.TraitMove, TargetX: x, TargetY: y}
		},
	}
}

func (s *Simulation) chopAction(x, y int, tree *entity.Tree) entity.GoapAction {
	return entity.GoapAction{
		Name: fmt.Sprintf("Chop(%d,%d)", x, y),
		Check: func(state entity.GoapState) bool {
			if state["at_x"] != x || state["at_y"] != y {
				return false
			}
			if state[treeHealthKey(x, y)] <= 0 {
				return false
			}
			newWood := state["wood"] + treeWoodYield
			return newWood*entity.WoodWeightPerUnit <= state["max_weight"]
		},
		Apply: func(state entity.GoapState) entity.GoapState {
			state[treeHealthKey(x, y)] = 0
			state["wood"] = state["wood"] + treeWoodYield
			state["weight"] = state["wood"] * entity.WoodWeightPerUnit
			return state
		},
		Cost: func(_ entity.GoapState) int { return treeHealth },
		BuildStep: func() entity.PlanStep {
			return entity.PlanStep{
				Trait:   entity.TraitChop,
				TargetX: x,
				TargetY: y,
				Tree:    tree,
			}
		},
	}
}

func (s *Simulation) depositWoodAction(x, y int) entity.GoapAction {
	return entity.GoapAction{
		Name: fmt.Sprintf("DepositWood(%d,%d)", x, y),
		Check: func(state entity.GoapState) bool {
			return state["at_x"] == x && state["at_y"] == y && state["wood"] > 0 && state["deposit_exists"] == 1
		},
		Apply: func(state entity.GoapState) entity.GoapState {
			state["wood"] = 0
			state["weight"] = 0
			return state
		},
		Cost: func(_ entity.GoapState) int { return 1 },
		BuildStep: func() entity.PlanStep {
			return entity.PlanStep{Trait: entity.TraitDeposit, TargetX: x, TargetY: y}
		},
	}
}

func (s *Simulation) buildDepositActions() []entity.GoapAction {
	var actions []entity.GoapAction
	for _, d := range s.deposits {
		actions = append(actions, s.moveToAction(d.X, d.Y))
		actions = append(actions, s.depositWoodAction(d.X, d.Y))
	}
	return actions
}

func (s *Simulation) setDepositPlan(v *entity.Villager) {
	state := s.buildState(v)
	goal := entity.GoapState{"wood": 0}
	actions := s.buildDepositActions()
	plan := s.planner.Plan(state, goal, actions)
	if plan != nil {
		v.SetPlan(plan)
	}
}

func (s *Simulation) setPlan(v *entity.Villager, job entity.Job) {
	state := s.buildState(v)
	goal := s.jobToGoal(job)
	actions := s.buildActions(v, job)
	if goal == nil || actions == nil {
		return
	}
	plan := s.planner.Plan(state, goal, actions)
	if plan == nil {
		return
	}
	v.SetPlan(plan)
}

func (s *Simulation) removeDeadTrees() {
	alive := s.trees[:0]
	for _, t := range s.trees {
		if t.Health > 0 {
			alive = append(alive, t)
		}
	}
	s.trees = alive
}

func (s *Simulation) AdvanceTicks(n int) {
	for i := 0; i < n; i++ {
		s.Tick()
	}
}

func (s *Simulation) SetTarget(villagerID string, x, y int) {
	for _, v := range s.villagers {
		if v.ID == villagerID {
			s.setPlan(v, entity.Job{Type: entity.JobTypeMove, TargetX: x, TargetY: y})
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
	s.jobQueue.Push(job)
}

func (s *Simulation) ProcessAxeSelection(cells [][2]int) {
	for _, cell := range cells {
		col, row := cell[0], cell[1]
		tree := s.TreeAt(col, row)
		if tree == nil {
			continue
		}
		s.PushJob(entity.Job{
			Type:    entity.JobTypeChopTrees,
			TargetX: tree.X,
			TargetY: tree.Y,
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

func (s *Simulation) AddDeposit(x, y int) {
	s.deposits = append(s.deposits, entity.NewDeposit(x, y))
}

func (s *Simulation) VillagerWood(id string) int {
	for _, v := range s.villagers {
		if v.ID == id {
			return v.Wood
		}
	}
	return -1
}

func (s *Simulation) VillagerMaxCarryWeight(id string) int {
	for _, v := range s.villagers {
		if v.ID == id {
			return v.MaxCarryWeight
		}
	}
	return -1
}

func (s *Simulation) SetMaxCarryWeight(id string, weight int) {
	for _, v := range s.villagers {
		if v.ID == id {
			v.MaxCarryWeight = weight
			return
		}
	}
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
			ID:             v.ID,
			Name:           v.Name(),
			Type:           int(v.Type),
			X:              v.X,
			Y:              v.Y,
			Wood:           v.Wood,
			MaxCarryWeight: v.MaxCarryWeight,
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

	jobs := make([]save.JobSave, len(s.jobQueue.Get()))
	for i, j := range s.jobQueue.Get() {
		jobs[i] = save.JobSave{
			Type:    int(j.Type),
			TargetX: j.TargetX,
			TargetY: j.TargetY,
		}
	}

	deposits := make([]save.DepositSave, len(s.deposits))
	for i, d := range s.deposits {
		deposits[i] = save.DepositSave{X: d.X, Y: d.Y}
	}

	return save.Save{
		Version:   1,
		World:     save.WorldSave{Rows: s.world.Rows(), Cols: s.world.Cols(), Cells: cells},
		Villagers: villagers,
		Trees:     trees,
		Deposits:  deposits,
		Jobs:      jobs,
	}
}

func (s *Simulation) QueueJobs() []entity.Job {
	return s.jobQueue.Get()
}

func (s *Simulation) Entities() []entity.Entity {
	total := len(s.villagers) + len(s.trees) + len(s.deposits)
	all := make([]entity.Entity, 0, total)
	for _, v := range s.villagers {
		all = append(all, v)
	}
	for _, t := range s.trees {
		all = append(all, t)
	}
	for _, d := range s.deposits {
		all = append(all, d)
	}
	return all
}

func (s *Simulation) debugSimulation() {
	if debug.IsEnabled(debug.Sim) {
		fmt.Printf("[SIMULATION] Sim tick=%d villagers=%d trees=%d deposits=%d jobs=%d\n",
			s.tickCount, len(s.villagers), len(s.trees), len(s.deposits), len(s.jobQueue.Get()))
	}
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
