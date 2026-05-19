package constants

const (
	ScreenW  = 1200
	ScreenH  = 1200
	GridCols = 50
	GridRows = 50
	TileSize = ScreenW / GridCols

	GridOffsetX = 0
	GridOffsetY = 0
)

func GridToScreen(col, row int) (x, y float32) {
	return float32(col)*TileSize + GridOffsetX, float32(row)*TileSize + GridOffsetY
}

func ScreenToGrid(mx, my int32) (int, int, bool) {
	gx := int(mx-GridOffsetX) / TileSize
	gy := int(my-GridOffsetY) / TileSize
	if gx < 0 || gx >= GridCols || gy < 0 || gy >= GridRows {
		return 0, 0, false
	}
	return gx, gy, true
}
