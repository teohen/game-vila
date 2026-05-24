package constants

const (
	ScreenW  = 1200
	ScreenH  = 1200
	GridCols = 100
	GridRows = 100
	TileSize = ScreenW / GridCols

	GridOffsetX = 0
	GridOffsetY = 0

	CameraZoomMin = 0.25
	CameraZoomMax = 7.0

	TickInterval = 200 // milliseconds between ticks (5 ticks/sec)
)

func WorldToScreen(col, row int) (x, y float32) {
	return float32(col)*TileSize + GridOffsetX, float32(row)*TileSize + GridOffsetY
}

func ScreenToWorld(mx, my int32) (int, int, bool) {
	gx := int(mx-GridOffsetX) / TileSize
	gy := int(my-GridOffsetY) / TileSize
	if gx < 0 || gx >= GridCols || gy < 0 || gy >= GridRows {
		return 0, 0, false
	}
	return gx, gy, true
}
