# mgm-tto

A management game where the player builds a village with Villagers, animals, and buildings on a top-down 2D grid world.

## Language

**Player**:
A person playing the game. Controls a camera to pan/zoom the world and selects tiles to issue commands.
_Avoid_: User, client

**Village**:
What the player builds over time. Composed of buildings, Villagers, and animals. Not yet explicitly modeled (implicit goal).
_Avoid_: Town, city, settlement

**World**:
A fixed-size 2D grid of cells (50×50) where everything exists. Supports pan and zoom camera navigation.
_Avoid_: Map, level, board

**Cell**:
A single tile on the world grid. Has a terrain type (Empty, Grass, Dirt, Water) and a grid position (row, col).
_Avoid_: Tile, square, block

**Terrain**:
The surface type of a Cell. Determines visual appearance and background color.
_Avoid_: Ground, floor, biome

**Villager**:
A non-player character that inhabits the village. Has an ID, name, type (currently only Human), and a position on the world grid. Is drawn using a spritesheet texture.
_Avoid_: Character, person, unit

**Selection**:
A set of Cells highlighted by the player via drag-to-select. Used to target commands (commands not yet implemented).
_Avoid_: Highlight, pick, mark

**Tree**:
A world entity that yields Wood when chopped by a Villager. Has a grid position (x, y), health, and a WoodYield value. Occupies its cell (blocks movement). Does not move or act autonomously.
_Avoid_: Resource, harvestable, node, object

**WoodYield**:
The amount of Wood a Tree drops when it is fully chopped by a Villager. A property of the Tree, not a separate entity.
_Avoid_: Drop, loot, output

**Walkable**:
A property of a Cell that is true when Villagers can move through it. Determined statically by CellType — Water is not walkable, Empty/Grass/Dirt are walkable. Never changes at runtime.
_Avoid_: Passable, traversable, pathable

**Occupied**:
A dynamic property of a Cell tracked by the World. True when a solid entity (Villager, Tree, etc.) currently stands on that cell. Used by pathfinding to avoid collision. Changes as entities move or are destroyed.
_Avoid_: Blocked, taken, full, collided

**Pathfinding**:
The system by which a Villager determines the sequence of cells to traverse toward a target. Uses A* on the 50×50 grid. Produces a list of waypoints the Villager follows one cell per tick. Inputs: walkable cells (static terrain) and occupied cells (dynamic entities).
_Avoid_: Navigation, routing, move planning

**JobQueue**:
A future system that holds work orders (e.g., chop trees) for Villagers to consume. Villagers pull from this queue to determine their movement target. Not yet implemented; a stub is used during initial movement implementation.
_Avoid_: TaskList, work queue, order book

**Camera**:
A 2D viewport that the player controls with right-click drag (pan) and mouse wheel (zoom). Clamped between 25% and 300%.
_Avoid_: Viewport, view

**SpriteBank**:
A package-level variable that owns loaded textures and makes them available to any entity that needs to draw. Exposes `Terrain` and `Human` textures. Created by the `spritebank` package, lifecycle managed via `LoadAll()` / `UnloadAll()`.
_Avoid_: TextureManager, AssetRegistry, resource cache

**Tick**:
The fundamental unit of game time. A pulse that fires at a fixed interval (configurable, independent of frame rate). On each tick the game advances its simulation: Villagers step, actions are processed, resources are updated. Ticks are deterministic — same interval always yields same behavior regardless of FPS.
_Avoid_: Frame, step, beat, cycle, turn
