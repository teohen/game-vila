package game

import (
	"fmt"
	"math"

	"github/teohen/mgm-tto/constants"
	"github/teohen/mgm-tto/entity"
	"github/teohen/mgm-tto/pathfinding"
	"github/teohen/mgm-tto/world"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const debugPrintInterval = 60

type Game struct {
	world   world.World
	npcs    []*entity.NPC
	trees   []entity.Tree
	camera  rl.Camera2D

	isDragging      bool
	dragStart       rl.Vector2
	dragEnd         rl.Vector2
	selected        map[rl.Vector2]bool

	debugMode       bool
	debugFrameCount int

	clock    Clock
	jobQueue entity.JobQueue
}

func New() Game {
	w := world.NewWorld(constants.GridRows, constants.GridCols)

	g := Game{
		world: w,
		npcs:  []*entity.NPC{},
		trees: []entity.Tree{},
		camera: rl.Camera2D{
			Target:   rl.NewVector2(float32(constants.GridCols)*constants.TileSize/2, float32(constants.GridRows)*constants.TileSize/2),
			Offset:   rl.NewVector2(constants.ScreenW/2, constants.ScreenH/2),
			Rotation: 0,
			Zoom:     1.0,
		},
		selected: make(map[rl.Vector2]bool),
		clock:    newClock(),
		jobQueue: entity.NewJobQueue(),
	}

	return g
}

func (g *Game) AddNPC(npc *entity.NPC) {
	g.npcs = append(g.npcs, npc)
	g.world.Occupy(npc.X, npc.Y)
}

func (g *Game) AddTree(tree entity.Tree) {
	g.trees = append(g.trees, tree)
	g.world.Occupy(tree.X, tree.Y)
}

func (g *Game) Input() {
	screenPos := rl.GetMousePosition()
	worldPos := g.screenToWorld(screenPos)

	if rl.IsMouseButtonDown(rl.MouseRightButton) {
		delta := rl.GetMouseDelta()
		g.pan(delta.X, delta.Y)
	}

	wheel := rl.GetMouseWheelMove()
	if wheel != 0 {
		g.zoom(1.0 + wheel*0.1)
	}

	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		g.isDragging = true
		g.dragStart = worldPos
		g.dragEnd = worldPos
	}

	if rl.IsMouseButtonReleased(rl.MouseLeftButton) && g.isDragging {
		g.dragEnd = worldPos
	}

	if rl.IsMouseButtonReleased(rl.MouseLeftButton) && g.isDragging {
		g.isDragging = false
		g.selectTilesInBox()
	}

	if rl.IsKeyPressed(rl.KeyF5) {
		g.debugMode = !g.debugMode
	}
}

func (g *Game) Update() {
	dt := float64(rl.GetFrameTime()) * 1000.0
	ticks := g.clock.Advance(dt)

	for i := 0; i < ticks; i++ {
		for _, npc := range g.npcs {
			g.updateNPC(npc)
		}
	}

	if !g.debugMode {
		return
	}

	g.debugFrameCount++
	if g.debugFrameCount < debugPrintInterval {
		return
	}
	g.debugFrameCount = 0

	screenPos := rl.GetMousePosition()
	worldPos := g.screenToWorld(screenPos)
	col := int(worldPos.X) / constants.TileSize
	row := int(worldPos.Y) / constants.TileSize

	g.debugPrint("FPS=%d zoom=%.2f tick=%d mouse=(%.0f,%.0f) cell=(%d,%d)",
		rl.GetFPS(), g.camera.Zoom, g.clock.TickCount(), screenPos.X, screenPos.Y, col, row)
}

func (g *Game) updateNPC(npc *entity.NPC) {
	switch npc.State {
	case entity.StateIdle:
		job := g.jobQueue.Pop()
		if job == nil {
			return
		}
		npc.TargetX = job.TargetX
		npc.TargetY = job.TargetY
		from := pathfinding.Point{X: npc.X, Y: npc.Y}
		to := pathfinding.Point{X: npc.TargetX, Y: npc.TargetY}
		path := pathfinding.FindPath(&g.world, from, to)
		if path == nil {
			return
		}
		npc.Waypoints = path
		npc.State = entity.StateMoving

	case entity.StateMoving:
		if len(npc.Waypoints) == 0 {
			npc.State = entity.StateArrived
			return
		}
		next := npc.Waypoints[0]
		if next.X == npc.TargetX && next.Y == npc.TargetY {
			npc.State = entity.StateArrived
			return
		}
		if g.world.IsOccupied(next.X, next.Y) {
			npc.State = entity.StateWaiting
			npc.WaitTicks = 0
			npc.WaitCount++
			return
		}
		g.world.Vacate(npc.X, npc.Y)
		npc.X = next.X
		npc.Y = next.Y
		g.world.Occupy(npc.X, npc.Y)
		npc.Waypoints = npc.Waypoints[1:]
		npc.WaitCount = 0

	case entity.StateWaiting:
		npc.WaitTicks++
		if npc.WaitTicks >= entity.WaitDuration {
			if npc.WaitCount >= entity.MaxRetries {
				npc.State = entity.StateIdle
				npc.WaitCount = 0
				npc.TargetX = 0
				npc.TargetY = 0
				npc.Waypoints = nil
				return
			}
			from := pathfinding.Point{X: npc.X, Y: npc.Y}
			to := pathfinding.Point{X: npc.TargetX, Y: npc.TargetY}
			path := pathfinding.FindPath(&g.world, from, to)
			if len(path) == 0 {
				npc.State = entity.StateIdle
				npc.WaitCount = 0
				return
			}
			npc.Waypoints = path
			npc.State = entity.StateMoving
		}

	case entity.StateArrived:
		npc.State = entity.StateIdle
		npc.WaitCount = 0
		npc.TargetX = 0
		npc.TargetY = 0
		npc.Waypoints = nil
	}
}

func (g *Game) Draw() {
	rl.BeginMode2D(g.camera)

	g.world.Draw()

	for _, tree := range g.trees {
		tree.Draw()
	}

	for _, npc := range g.npcs {
		npc.Draw()
	}

	if g.isDragging {
		w := g.dragEnd.X - g.dragStart.X
		h := g.dragEnd.Y - g.dragStart.Y
		rl.DrawRectangleLines(
			int32(g.dragStart.X),
			int32(g.dragStart.Y),
			int32(w),
			int32(h),
			rl.Blue,
		)
	}

	for pos := range g.selected {
		col, row := int(pos.X), int(pos.Y)
		x, y := constants.WorldToScreen(col, row)
		rl.DrawRectangleLines(
			int32(x), int32(y),
			int32(constants.TileSize), int32(constants.TileSize),
			rl.Red,
		)
	}

	rl.EndMode2D()
}

func (g *Game) screenToWorld(screenPos rl.Vector2) rl.Vector2 {
	return rl.GetScreenToWorld2D(screenPos, g.camera)
}

func (g *Game) pan(dx, dy float32) {
	g.camera.Target.X -= dx / g.camera.Zoom
	g.camera.Target.Y -= dy / g.camera.Zoom
}

func (g *Game) zoom(factor float32) {
	newZoom := g.camera.Zoom * factor
	if newZoom < constants.CameraZoomMin {
		newZoom = constants.CameraZoomMin
	}
	if newZoom > constants.CameraZoomMax {
		newZoom = constants.CameraZoomMax
	}
	g.camera.Zoom = newZoom
}

func (g *Game) selectTilesInBox() {
	minX := math.Min(float64(g.dragStart.X), float64(g.dragEnd.X))
	maxX := math.Max(float64(g.dragStart.X), float64(g.dragEnd.X))
	minY := math.Min(float64(g.dragStart.Y), float64(g.dragEnd.Y))
	maxY := math.Max(float64(g.dragStart.Y), float64(g.dragEnd.Y))

	g.selected = make(map[rl.Vector2]bool)

	for row := 0; row < g.world.Rows(); row++ {
		for col := 0; col < g.world.Cols(); col++ {
			x, y := constants.WorldToScreen(col, row)
			tileCenterX := float64(x + constants.TileSize/2)
			tileCenterY := float64(y + constants.TileSize/2)

			if tileCenterX >= minX && tileCenterX <= maxX &&
				tileCenterY >= minY && tileCenterY <= maxY {
				g.selected[rl.NewVector2(float32(col), float32(row))] = true
			}
		}
	}
}

func (g *Game) debugPrint(format string, args ...any) {
	if !g.debugMode {
		return
	}
	fmt.Printf("[DEBUG] "+format+"\n", args...)
}
