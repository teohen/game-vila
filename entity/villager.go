package entity

import (
	"github/teohen/mgm-tto/constants"
	"github/teohen/mgm-tto/spritebank"
	"github/teohen/mgm-tto/world"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	WoodWeightPerUnit    = 2
	DefaultMaxCarryWeight = 50
)

type VillagerType int

const (
	Human VillagerType = 1
)

type Villager struct {
	Movement
	Lumberjack
	ID             string
	name           string
	Type           VillagerType
	Wood           int
	MaxCarryWeight int

	plan []PlanStep
	step int
}

func NewVillager(id, name string, x, y int) *Villager {
	return &Villager{
		Movement: Movement{
			X: x,
			Y: y,
		},
		ID:             id,
		name:           name,
		Type:           Human,
		MaxCarryWeight: DefaultMaxCarryWeight,
	}
}

func (v *Villager) CurrentWeight() int {
	return v.Wood * WoodWeightPerUnit
}

func (v *Villager) IsCarrying() bool {
	return v.CurrentWeight() >= v.MaxCarryWeight
}

func (v *Villager) Tick(w *world.World) MovementEvent {
	if len(v.plan) == 0 {
		if v.IsCarrying() {
			return EventNone
		}
		return EventIdle
	}

	step := v.plan[v.step]
	switch step.Trait {
	case TraitMove:
		if v.State == StateIdle {
			v.Movement.SetTarget(step.TargetX, step.TargetY, w)
		}
		event := v.Movement.Update(w)
		if event == EventArrived {
			v.step++
			if v.step >= len(v.plan) {
				v.clearPlan()
				if v.IsCarrying() {
					return EventNone
				}
				return EventIdle
			}
		}
		return EventNone

	case TraitChop:
		if !v.IsHitting() {
			if step.Tree == nil {
				v.step++
				if v.step >= len(v.plan) {
					v.clearPlan()
					if v.IsCarrying() {
						return EventNone
					}
					return EventIdle
				}
				return EventNone
			}
			v.Lumberjack.Start(step.Tree)
		}
		wood, done := v.Lumberjack.Update(w)
		if wood > 0 {
			v.Wood += wood
		}
		if done {
			v.step++
			if v.step >= len(v.plan) {
				v.clearPlan()
				if v.IsCarrying() {
					return EventNone
				}
				return EventIdle
			}
		}
		return EventNone

	case TraitDeposit:
		v.Wood = 0
		v.step++
		v.State = StateIdle
		v.Waypoints = nil
		if v.step >= len(v.plan) {
			v.clearPlan()
			return EventIdle
		}
		return EventNone
	}

	return EventNone
}

func (v *Villager) SetPlan(plan []PlanStep) {
	v.plan = plan
	v.step = 0
	v.State = StateIdle
	v.WaitTicks = 0
	v.WaitCount = 0
	v.Waypoints = nil
}

func (v *Villager) IsIdle() bool {
	return len(v.plan) == 0
}

func (v *Villager) clearPlan() {
	v.plan = nil
	v.step = 0
}

func (v *Villager) Name() string {
	return v.name
}

func (v *Villager) Pos() (int, int) {
	return v.Movement.Pos()
}

func getSource(v *Villager) (rl.Rectangle, rl.Rectangle) {
	src := rl.NewRectangle(0, 0, 0, 0)
	dst := rl.NewRectangle(0, 0, 0, 0)

	switch v.Type {
	case Human:
		x, y := constants.WorldToScreen(v.X, v.Y)
		dst.X = x
		dst.Y = y
		dst.Width = constants.TileSize
		dst.Height = constants.TileSize
		src.X = 41
		src.Y = 21
		src.Width = 16
		src.Height = 19
	}

	return src, dst
}

func (v *Villager) Draw() {
	src, dst := getSource(v)
	rl.DrawTexturePro(spritebank.Human, src, dst, rl.NewVector2(0, 0), 0, rl.White)
}
