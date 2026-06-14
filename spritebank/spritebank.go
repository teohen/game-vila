package spritebank

import (
	"fmt"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var Terrain rl.Texture2D
var Human rl.Texture2D
var Structures rl.Texture2D
var Animals rl.Texture2D

func LoadAll() {
	Terrain = loadTexture("./res/assets/spr_terrains.png")
	Human = loadTexture("./res/assets/player_anims.png")
	Structures = loadTexture("./res/assets/spr_structures_exterior.png")
	Animals = loadTexture("./res/assets/animals.png")
}

func UnloadAll() {
	rl.UnloadTexture(Terrain)
	rl.UnloadTexture(Human)
	rl.UnloadTexture(Structures)
	rl.UnloadTexture(Animals)
}

func loadTexture(path string) rl.Texture2D {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic(fmt.Sprintf("spritebank: texture not found: %s", path))
	}
	return rl.LoadTexture(path)
}
