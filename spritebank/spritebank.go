package spritebank

import (
	"fmt"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var Terrain rl.Texture2D
var Human rl.Texture2D

func LoadAll() {
	Terrain = loadTexture("./res/assets/spr_terrains.png")
	Human = loadTexture("./res/assets/player_anims.png")
}

func UnloadAll() {
	rl.UnloadTexture(Terrain)
	rl.UnloadTexture(Human)
}

func loadTexture(path string) rl.Texture2D {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic(fmt.Sprintf("spritebank: texture not found: %s", path))
	}
	return rl.LoadTexture(path)
}
