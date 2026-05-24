package game

import (
	"fmt"
	"math"

	"github/teohen/mgm-tto/constants"
	"github/teohen/mgm-tto/debug"
	"github/teohen/mgm-tto/entity"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Tool int

const (
	ToolSelect Tool = iota
	ToolAxe
)

const debugPrintInterval = 60

type Game struct {
	simulation *Simulation
	camera     rl.Camera2D

	isDragging bool
	dragStart  rl.Vector2
	dragEnd    rl.Vector2
	selected   map[rl.Vector2]bool

	activeTool Tool

	debugMode       bool
	debugFrameCount int

	clock Clock
}

func New() Game {
	return Game{
		simulation: NewSimulationDefault(),
		camera: rl.Camera2D{
			Target:   rl.NewVector2(float32(constants.GridCols)*constants.TileSize/2, float32(constants.GridRows)*constants.TileSize/2),
			Offset:   rl.NewVector2(constants.ScreenW/2, constants.ScreenH/2),
			Rotation: 0,
			Zoom:     1.0,
		},
		selected: make(map[rl.Vector2]bool),
		clock:    newClock(),
	}
}

func (g *Game) AddVillager(v *entity.Villager) {
	g.simulation.AddVillager(v)
}

func (g *Game) AddTree(tree *entity.Tree) {
	g.simulation.AddTree(tree)
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

	if rl.IsMouseButtonDown(rl.MouseLeftButton) && g.isDragging {
		g.dragEnd = worldPos
	}

	if rl.IsMouseButtonReleased(rl.MouseLeftButton) && g.isDragging {
		g.dragEnd = worldPos
		g.isDragging = false
		g.selectCellsInBox()
		g.onSelectionComplete()
	}

	if rl.IsKeyPressed(rl.KeyF5) {
		g.debugMode = !g.debugMode
		debug.SetEnabled(g.debugMode)
		if g.debugMode {
			fmt.Printf("[DEBUG] Debug mode enabled — %s\n", debug.ActiveString())
		} else {
			fmt.Printf("[DEBUG] Debug mode disabled\n")
		}
	}

	if rl.IsKeyPressed(rl.KeyOne) {
		if g.activeTool == ToolAxe {
			g.activeTool = ToolSelect
			fmt.Println("[TOOL] Selection")
		} else {
			g.activeTool = ToolAxe
			fmt.Println("[TOOL] Axe")
		}
	}

	if g.debugMode {
		switch {
		case rl.IsKeyPressed(rl.KeyZero):
			debug.Toggle(debug.Sim)
		case rl.IsKeyPressed(rl.KeyTwo):
			debug.Toggle(debug.Move)
		case rl.IsKeyPressed(rl.KeyThree):
			debug.Toggle(debug.Path)
		case rl.IsKeyPressed(rl.KeyFour):
			debug.Toggle(debug.Clock)
		case rl.IsKeyPressed(rl.KeyFive):
			debug.Toggle(debug.Job)
		case rl.IsKeyPressed(rl.KeySix):
			debug.Toggle(debug.World)
		}
	}
}

func (g *Game) Update() {
	dt := float64(rl.GetFrameTime()) * 1000.0
	ticks := g.clock.Advance(dt)

	for i := 0; i < ticks; i++ {
		g.simulation.Tick()
	}

	if !g.debugMode {
		return
	}

	g.debugFrameCount++
	if g.debugFrameCount < debugPrintInterval {
		return
	}
	g.debugFrameCount = 0

	entities := g.simulation.Entities()
	villagers := 0
	for _, e := range entities {
		if _, ok := e.(*entity.Villager); ok {
			villagers++
		}
	}
	//g.debugPrint("FPS=%d zoom=%.2f tick=%d ents=%d villagers=%d tool=%d selected=%d",
	//rl.GetFPS(), g.camera.Zoom, g.clock.TickCount(), len(entities), villagers, g.activeTool, len(g.selected))
}

func (g *Game) Draw() {
	rl.BeginMode2D(g.camera)

	g.simulation.World().Draw()

	for _, e := range g.simulation.Entities() {
		e.Draw()
	}

	if g.isDragging {
		x := float32(math.Min(float64(g.dragStart.X), float64(g.dragEnd.X)))
		y := float32(math.Min(float64(g.dragStart.Y), float64(g.dragEnd.Y)))
		w := float32(math.Abs(float64(g.dragEnd.X - g.dragStart.X)))
		h := float32(math.Abs(float64(g.dragEnd.Y - g.dragStart.Y)))
		rl.DrawRectangleLines(int32(x), int32(y), int32(w), int32(h), rl.Blue)
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

func (g *Game) selectCellsInBox() {
	g.selected = make(map[rl.Vector2]bool)

	minX := math.Min(float64(g.dragStart.X), float64(g.dragEnd.X))
	maxX := math.Max(float64(g.dragStart.X), float64(g.dragEnd.X))
	minY := math.Min(float64(g.dragStart.Y), float64(g.dragEnd.Y))
	maxY := math.Max(float64(g.dragStart.Y), float64(g.dragEnd.Y))

	if maxX-minX < 4 && maxY-minY < 4 {
		col := int(g.dragEnd.X) / constants.TileSize
		row := int(g.dragEnd.Y) / constants.TileSize
		if col >= 0 && col < constants.GridCols && row >= 0 && row < constants.GridRows {
			g.selected[rl.NewVector2(float32(col), float32(row))] = true
		}
		return
	}

	for row := 0; row < g.simulation.World().Rows(); row++ {
		for col := 0; col < g.simulation.World().Cols(); col++ {
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

func (g *Game) onSelectionComplete() {
	switch g.activeTool {
	case ToolAxe:
		cells := make([][2]int, 0, len(g.selected))
		for pos := range g.selected {
			cells = append(cells, [2]int{int(pos.X), int(pos.Y)})
		}
		g.simulation.ProcessAxeSelection(cells)
	}
}

func (g *Game) debugPrint(format string, args ...any) {
	if !g.debugMode {
		return
	}
	fmt.Printf("[DEBUG] "+format+"\n", args...)
}
