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
A fixed-size 2D grid of cells (30×30) where everything exists. Supports pan and zoom camera navigation.
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
A set of Cells highlighted by the player via drag-to-select. Used to target tools.
_Avoid_: Highlight, pick, mark

**Tree**:
A world entity that yields Wood when chopped by a Villager. Has a grid position (x, y), Health, and a WoodYield value. Occupies its cell (blocks movement). Does not move or act autonomously.
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
The system by which a Villager determines the sequence of cells to traverse toward a target. Uses A* on the 30×30 grid. Produces a list of waypoints the Villager follows one cell per tick. Inputs: walkable cells (static terrain) and occupied cells (dynamic entities).
_Avoid_: Navigation, routing, move planning

**Job**:
A work order in the JobQueue. Has a Type (JobTypeChopTrees, JobTypeMove) and a target cell (TargetX, TargetY). Created by the player via tool actions (e.g., axe tool creates ChopTrees jobs).
_Avoid_: Task, order, work item

**JobQueue**:
A queue that holds Jobs for Villagers to consume. Villagers pull from this queue to determine their movement target. Created by tool actions; consumed by any Idle/Arrived Villager.
_Avoid_: TaskList, work queue, order book

**JobType**:
The kind of work a Job represents (e.g., JobTypeChopTrees, JobTypeMove). Determines what a Villager does upon arrival. A JobType maps to a Trait that provides the capacity to execute it — e.g., the Lumberjack Trait handles JobTypeChopTrees.
_Avoid_: Job kind, work type, action type

**Camera**:
A 2D viewport that the player controls with right-click drag (pan) and mouse wheel (zoom). Clamped between 25% and 700%.
_Avoid_: Viewport, view

**SpriteBank**:
A package-level variable that owns loaded textures and makes them available to any entity that needs to draw. Exposes `Terrain` and `Human` textures. Created by the `spritebank` package, lifecycle managed via `LoadAll()` / `UnloadAll()`.
_Avoid_: TextureManager, AssetRegistry, resource cache

**Tool**:
A mode the Player can activate to change what happens when they drag-select Cells. Pressing '1' toggles between ToolSelect (highlight cells only) and ToolAxe (highlight + create ChopTrees Jobs).
_Avoid_: Mode, weapon, item

**Trait**:
A composable unit of behavior and state embedded into an Entity. Each trait owns its own logic and lifecycle. Traits are mixed into entities via Go struct embedding. An entity may compose zero or more traits.
_Avoid_: Component, module, plugin, system

**Movement**:
A trait that handles grid-based movement. Owns position (X, Y), target (TargetX, TargetY), path (Waypoints), and a state machine (Idle, Moving, Waiting, Arrived). Exposes `Update(world)` which advances one tick of movement, `SetTarget(x, y, world)` which initiates pathfinding toward a destination, and `Pos()` which returns current coordinates. Uses A* pathfinding via the pathfinding package and manages Occupy/Vacate on the World as it moves. Meant to be embedded in any mobile entity (Villager, future Animals, Vehicles).
_Avoid_: Navigation, locomotion, mover

**Tick**:
The fundamental unit of game time. A pulse that fires at a fixed interval (configurable, independent of frame rate). On each tick the game advances its simulation: entities Tick, actions are processed, resources are updated. Ticks are deterministic — same interval always yields same behavior regardless of FPS.
_Avoid_: Frame, step, beat, cycle, turn

**Noise**:
A 2D Perlin noise function used for procedural world generation. Lives in the `world` package. Produces smooth continuous values in [-1, 1]. Used to determine terrain layout and forest density.
_Avoid_: Random, Perlin (implementation detail)

**Noise invocation**:
A call to Noise at a given frequency and seed. Two invocations are used: one at frequency ≈0.035 for terrain type (Water/Dirt/Grass) and another at frequency ≈0.07 for forest density. Each uses an independent seed to avoid correlation.
_Avoid_: Channel, layer, octave

**Save**:
A serialized snapshot of the whole Game state written to disk as a JSON file. Contains the World grid, Villagers (with movement state), Trees, Jobs, and Camera position/zoom. Created by pressing F9.
_Avoid_: Save file, save data, save slot

**Load**:
The act of reading a Save from disk and replacing the current Game state entirely with the reconstructed state. Triggered by pressing F10. If the file is missing or corrupt, an error is printed and the game continues unchanged.
_Avoid_: Restore, open, import
