# Mouse Selection Box Implementation Plan

Based on your Raylib Go project with an 18x12 tile grid (64px tiles).

---

## Overview

Capture mouse drag selection to mark multiple tiles as selected and change their visual appearance.

---

## 1. Track Mouse State

Add selection state to your grid in `grid/grid.go`:

```go
type Grid struct {
    cells       [][]Tile
    texture     rl.Texture2D
    // Add these:
    isDragging  bool
    dragStart   rl.Vector2  // screen coords
    dragEnd     rl.Vector2  // screen coords
    selected    map[rl.Vector2]bool  // grid coords of selected tiles
}
```

---

## 2. Capture Selection Box

In `main.go`, update the `input()` function:

```go
func input() {
    mousePos := rl.GetMousePosition()

    // Start dragging (left button pressed)
    if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
        g.isDragging = true
        g.dragStart = mousePos
        g.dragEnd = mousePos
    }

    // While dragging, update end position
    if rl.IsMouseButtonDown(rl.MouseLeftButton) && g.isDragging {
        g.dragEnd = mousePos
    }

    // End dragging (left button released)
    if rl.IsMouseButtonReleased(rl.MouseLeftButton) && g.isDragging {
        g.isDragging = false
        g.selectTilesInBox()
    }
}
```

---

## 3. Select Tiles Within Box

Add this method to `grid/grid.go`:

```go
func (g *Grid) selectTilesInBox() {
    // Convert screen rect to grid bounds
    minX := math.Min(g.dragStart.X, g.dragEnd.X)
    maxX := math.Max(g.dragStart.X, g.dragEnd.X)
    minY := math.Min(g.dragStart.Y, g.dragEnd.Y)
    maxY := math.Max(g.dragStart.Y, g.dragEnd.Y)

    g.selected = make(map[rl.Vector2]bool)

    for row := 0; row < len(g.cells); row++ {
        for col := 0; col < len(g.cells[row]); col++ {
            // Get tile screen position
            x, y := GridToScreen(col, row)
            tileCenterX := x + constants.TileSize/2
            tileCenterY := y + constants.TileSize/2

            // Check if tile center is within selection box
            if tileCenterX >= minX && tileCenterX <= maxX &&
               tileCenterY >= minY && tileCenterY <= maxY {
                g.selected[rl.NewVector2(float32(col), float32(row))] = true
            }
        }
    }
}
```

---

## 4. Add Selected Field to Tile

In `grid/tile.go`:

```go
type Tile struct {
    Type     TileType
    texture  rl.Texture2D
    row      int
    col      int
    size     float32
    Selected bool  // add this
}
```

---

## 5. Change Color on Selected Tiles

Update `Tile.draw()` in `grid/tile.go`:

```go
func (t *Tile) draw() {
    src, dst := getSource(t)
    color := rl.White
    if t.Selected {
        color = rl.NewColor(200, 200, 255, 255) // lighter tint
    }
    rl.DrawTexturePro(t.texture, src, dst, rl.NewVector2(0, 0), 0, color)
}
```

---

## 6. Render Selection Box and Highlights

Update `Grid.Draw()` in `grid/grid.go`:

```go
func (g *Grid) Draw() {
    // Draw all tiles (they'll use their Selected state for color)
    for _, row := range g.cells {
        for _, tile := range row {
            tile.draw()
        }
    }

    // Draw selection box while dragging
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

    // Draw highlight border on selected tiles
    for pos := range g.selected {
        col, row := int(pos.X), int(pos.Y)
        x, y := GridToScreen(col, row)
        rl.DrawRectangleLines(
            int32(x), int32(y),
            int32(constants.TileSize), int32(constants.TileSize),
            rl.Red,
        )
    }
}
```

---

## 7. Add Required Import

In `grid/grid.go`, add `math` to imports:

```go
import (
    "fmt"
    "math"
    "github/teohen/mgm-tto/constants"
    rl "github.com/gen2brain/raylib-go/raylib"
)
```

---

## Key Functions Reference

| Step | Raylib Function |
|------|-----------------|
| Capture start drag | `IsMouseButtonPressed(rl.MouseLeftButton)` |
| Track drag | `IsMouseButtonDown(rl.MouseLeftButton)` + `GetMousePosition()` |
| End drag | `IsMouseButtonReleased(rl.MouseLeftButton)` |
| Draw box | `DrawRectangleLines()` |
| Draw highlight | `DrawRectangleLines()` with different color |

---

## Notes

- The existing `ScreenToGrid()` function can also be used to determine if tiles are within the selection box by checking their grid positions
- You may want to add a clear selection functionality (e.g., right-click or pressing a key)
- Consider adding a "lasso selection" toggle for additive selection vs. replacing selection