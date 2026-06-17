package simulation

import (
	"github/teohen/mgm-tto/building"
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/events"
	"github/teohen/mgm-tto/goap"
	"github/teohen/mgm-tto/inventory"
	"github/teohen/mgm-tto/job"
	"github/teohen/mgm-tto/npc"
	"github/teohen/mgm-tto/resource"
	"github/teohen/mgm-tto/world"
	"log"
	"time"
)

type Tool int

const (
	ToolSelect Tool = iota
	ToolAxe
)

type Simulation struct {
	tickCount int
	world     *world.World
	npcs      *npc.NPCList
	resources []resource.IResource
	buildings *building.BuildingsList
	inventory *inventory.Inventory

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

func New(worldSize int, initialSeed int64) *Simulation {
	seed := int64(0)
	if initialSeed == 0 {
		seed = time.Now().UnixNano()
	}

	w := world.NewWorld(worldSize, worldSize)
	w.Generate(seed)

	sim := Simulation{
		world:     &w,
		npcs:      npc.NewNPCList(),
		buildings: building.NewBuildingsList(),
		inventory: inventory.NewInventory(),
	}

	forestNoise := world.NewNoise(seed + 1)
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
			t := resource.NewTree(c, r, treeHealth, treeWoodYield)
			sim.AddTree(t)
		}
	}

	return &sim
}

func NewEmpty(worldSize int) *Simulation {
	w := world.NewWorld(worldSize, worldSize)

	sim := Simulation{
		world:     &w,
		npcs:      npc.NewNPCList(),
		buildings: building.NewBuildingsList(),
	}

	return &sim
}

func (s *Simulation) Tick() {
	for _, v := range s.npcs.NPCS() {
		v.Tick()
	}

	s.processEvents()
	s.tickCount++
}

func (s *Simulation) AddVillager(v *npc.Villager) {
	s.npcs.AddNPC(v)
	s.world.Occupy(v.Pos().X, v.Pos().Y)
}

func (s *Simulation) AddTree(tree *resource.Tree) {
	s.resources = append(s.resources, tree)
}

func (s *Simulation) AddResource(res resource.IResource) {
	s.resources = append(s.resources, res)
}

func (s *Simulation) RemoveResource(res resource.IResource) {
	for i, r := range s.resources {
		if res.ID() == r.ID() {
			s.resources = append(s.resources[:i], s.resources[i+1:]...)
			return
		}
	}
}

func (s *Simulation) AddDeer(deer *npc.Deer) {
	s.npcs.AddNPC(deer)
	s.world.Occupy(deer.Pos().X, deer.Pos().Y)
}

func (s *Simulation) RemoveTree(x, y int) bool {
	p := cnts.Point{
		X: x,
		Y: y,
	}

	for i, t := range s.resources {
		if t.Type() == resource.ResourceTreeType && t.Pos() == p {
			s.resources = append(s.resources[:i], s.resources[i+1:]...)
			return true
		}
	}
	return false
}

func (s *Simulation) AddStorage(storage *building.Storage) {
	s.buildings.AddBuilding(storage)
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
		tree.Mark(true)
		s.PushJob(*job.NewJob(job.JobChopTreeType, tree))
	}
}

func (s *Simulation) TreeAt(x, y int) *resource.Tree {
	p := cnts.Point{
		X: x, Y: y,
	}
	for _, t := range s.resources {
		if t.Type() == resource.ResourceTreeType && t.Pos() == p {
			tree, ok := t.(*resource.Tree)
			if !ok {
				log.Fatal("TREE CONVERTION NOT WOTK")
			}
			return tree
		}
	}
	return nil
}

func (s *Simulation) World() *world.World {
	return s.world
}

func (s *Simulation) Buildings() []building.Building {
	return s.buildings.Buldings()
}

func (s *Simulation) NPCs() []npc.NPC {
	return s.npcs.NPCS()
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

func (s *Simulation) Resources() []resource.IResource {
	return s.resources
}

func (s *Simulation) processEvents() {
	for {
		select {
		case evt := <-events.EventQueue:
			switch evt.Type {
			case events.EventTreeCut:
				treePos := evt.Payload["pos"].(cnts.Point)
				s.RemoveTree(treePos.X, treePos.Y)
				woodYeld := evt.Payload["woodYield"].(int)
				s.AddResource(resource.NewWood(treePos, woodYeld))
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

func (s *Simulation) CollectResourceAt(pos cnts.Point) int {
	for _, r := range s.resources {
		if r.Collectable() && r.Pos() == pos {
			s.RemoveResource(r)
			return r.Amount()
		}
	}
	return 0
}
