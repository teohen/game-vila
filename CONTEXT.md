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
A fixed-size 2D grid of cells  where everything exists. Supports pan and zoom camera navigation.
_Avoid_: Map, level, board

**Cell**:
A single tile on the world grid. Has a terrain type (Grass, Dirt, Water) and a grid position (row, col).
_Avoid_: Tile, square, block

**Terrain**:
The surface type of a Cell. Determines visual appearance and background color.
_Avoid_: Ground, floor, biome

**Villager**:
A non-player character that inhabits the village. Has an ID and a position on the world grid. Composes Movement, Lumberjack, Collecter, and Storager traits. Uses GOAP (Goal-Oriented Action Planning) for decision-making. Drawn using a spritesheet texture.
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
A work order in the JobQueue. Has a Type (currently only JobChopTreeType) and an Object (the target entity). Created by the player via tool actions (e.g., axe tool creates ChopTree jobs).
_Avoid_: Task, order, work item

**JobQueue**:
A global queue that holds Jobs for Villagers to consume. Villagers scan the queue for the closest reachable job (using A* pathfinding) and remove the chosen job. Created by tool actions; consumed by Idle Villagers.
_Avoid_: TaskList, work queue, order book

**JobType**:
The kind of work a Job represents (currently only JobChopTreeType). Determines what a Villager does upon arrival. A JobType maps to a GOAP action which triggers the corresponding trait — e.g., JobChopTreeType triggers ActionChopTree, executed by the Lumberjack trait.
_Avoid_: Job kind, work type, action type

**Camera**:
A 2D viewport that the player controls with right-click drag (pan) and mouse wheel (zoom). Clamped between 25% and 700%.
_Avoid_: Viewport, view

**SpriteBank**:
A package-level variable that owns loaded textures and makes them available to any entity that needs to draw. Exposes `Terrain`, `Human`, `Structures`, and `Animals` textures. Created by the `spritebank` package, lifecycle managed via `LoadAll()` / `UnloadAll()`.
_Avoid_: TextureManager, AssetRegistry, resource cache

**Tool**:
A mode the Player can activate to change what happens when they drag-select Cells. Pressing '1' toggles between ToolSelect (highlight cells only) and ToolAxe (highlight + create ChopTrees Jobs).
_Avoid_: Mode, weapon, item

**Debug**:
A simple debug flag (`cnts.DEBUGGING`) set via the `-debug` CLI flag. When enabled, draws grid coordinates on cells. No per-category debug keys exist.
_Avoid_: Logging, tracing, verbose

**Trait**:
A composable unit of behavior and state embedded into an Entity. Each trait owns its own logic and lifecycle. Traits are mixed into entities via Go struct embedding. An entity may compose zero or more traits.
_Avoid_: Component, module, plugin, system

**Movement**:
A trait that handles grid-based movement. Owns position (`cnts.Point`), target (`cnts.Point`), path (Waypoints), and a state machine (Idle, Moving, Waiting, Arrived). Exposes `Update()` which advances one tick of movement, `SetTarget(target cnts.Point)` which initiates pathfinding toward a destination, and `Pos()` which returns current coordinates. When the next waypoint is occupied, Movement enters the Waiting state for up to `WaitDuration` (5 ticks) per retry, with a maximum of `MaxRetries` (10) before giving up and returning to idle. Uses A* pathfinding via the pathfinding package and manages Occupy/Vacate on the World as it moves. Meant to be embedded in any mobile entity (Villager, future Animals, Vehicles).
_Avoid_: Navigation, locomotion, mover

**Tick**:
The fundamental unit of game time. A pulse that fires at a fixed interval (configurable, independent of frame rate). On each tick the game advances its simulation: entities Tick, actions are processed, resources are updated. Ticks are deterministic — same interval always yields same behavior regardless of FPS.
_Avoid_: Frame, step, beat, cycle, turn

**Clock**:
A fixed-tick accumulator that decouples simulation speed from frame rate. Owns an `accumulator` (milliseconds) and an `interval` of 200ms (5 ticks/sec). Each frame, `Advance(dtMs)` adds the frame delta and fires as many ticks as the accumulator holds. Lives in the `game` package.
_Avoid_: Timer, stopwatch, scheduler

**Noise**:
A 2D Perlin noise function used for procedural world generation. Lives in the `world` package. Produces smooth continuous values in [-1, 1]. Used to determine terrain layout and forest density.
_Avoid_: Random, Perlin (implementation detail)

**Noise invocation**:
A call to Noise at a given frequency and seed. Two invocations are used: one at frequency ≈0.035 for terrain type (Water/Dirt/Grass) and another at frequency ≈0.07 for forest density. The terrain invocation uses two thresholds: cells below -0.15 become Water, between -0.15 and 0.05 become Dirt, and above 0.05 become Grass. The forest invocation uses a threshold of 0.1 — only Grass cells above this value receive a Tree. Each uses an independent seed to avoid correlation.
_Avoid_: Channel, layer, octave

**Save**:
A serialized snapshot of the whole Game state written to disk as a JSON file. Contains the World grid, Villagers (with movement state), Trees, Jobs, and Camera position/zoom. Types exist in the `save` package but are not currently wired to keyboard shortcuts.
_Avoid_: Save file, save data, save slot

**Load**:
The act of reading a Save from disk and replacing the current Game state entirely with the reconstructed state. `LoadFromFile()` exists in the `save` package but is currently stubbed (not wired to the game loop).
_Avoid_: Restore, open, import

**Game**:
The top-level struct that wires together the Simulation, UI, and Clock. Created by `game.New()` (procedural world). The main loop calls `UI.Input()`, `Update()` (advances the Clock and ticks the Simulation), and `UI.Draw()` each frame.
_Avoid_: App, engine, state
