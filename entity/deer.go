package entity

import (
	"github/teohen/mgm-tto/agent"
	"github/teohen/mgm-tto/building"
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/spritebank"
	"github/teohen/mgm-tto/world"
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type DeerState string

const (
	StateIdle    DeerState = "Idle"
	StateRoaming DeerState = "roaming"
)

type Deer struct {
	agent    agent.IAgent
	movement *Movement
	id       string
	State    DeerState
	w        *world.World
}

func NewDeer(x, y int, w *world.World) *Deer {
	d := Deer{id: cnts.NewID()}
	mv := NewMovement(x, y, w)
	d.agent = agent.NewAgent(x, y, w, "Deer", "walkable")
	d.agent.RegisterActor(agent.ActionMoveType, mv)
	d.movement = mv
	d.State = StateIdle

	d.w = w
	return &d
}

func (d *Deer) Tick(w *world.World, entities *[]Entity, buildings *building.BuildingsList) {
	if d.State == StateIdle {
		d.AddRoamGoal(w, entities)
	}

	switch d.State {
	case StateIdle:
		if found := d.agent.ChooseGoal(d.w, d.Pos()); found {
			d.State = StateRoaming
		}

	case StateRoaming:
		if finalAction := d.agent.ExecutePlan(); finalAction {
			d.State = StateIdle
		}
	}
}

func (d *Deer) Pos() cnts.Point {
	return d.movement.pos
}

func (d *Deer) ID() string {
	return d.id
}

func (d *Deer) Draw() {
	x, y := cnts.WorldToScreen(d.Pos().X, d.Pos().Y)
	src := rl.NewRectangle(12, 134, 128, 184)
	dst := rl.NewRectangle(x, y, cnts.TileSize, cnts.TileSize)
	rl.DrawTexturePro(spritebank.Animals, src, dst, rl.NewVector2(0, 0), 0, rl.White)
}

func (d *Deer) Type() EntityType {
	return EntityTypeVillager
}

func (d *Deer) AddRoamGoal(w *world.World, entities *[]Entity) {
	newX := RandomIntRangeFast(d.Pos().X-5, d.Pos().X+5)
	newY := RandomIntRangeFast(d.Pos().Y-5, d.Pos().Y+5)

	pin := cnts.Pin{Id: cnts.NewID(), Position: cnts.Point{X: newX, Y: newY}}
	g := agent.NewGoalRoam(w, &pin, d.movement.pos)
	d.agent.AddGoal(g)
}

func RandomIntRangeFast(min, max int) int {
	return rand.N(max-min+1) + min
}
