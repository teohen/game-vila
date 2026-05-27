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
A single tile on the world grid. Has a terrain type (Grass, Dirt, Water) and a grid position (row, col).
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

**Plan**:
A sequence of PlanSteps that a Villager executes to fulfill a Job. Hardcoded per JobType in the Villager. When the current step completes (e.g., Movement arrives adjacent to the target), the Villager advances to the next step and activates the matching Trait.
_Avoid_: Task list, itinerary, route

**PlanStep**:
A single sub-objective in a Plan. Has a TraitType (e.g., Move, Chop) and a target grid position (TargetX, TargetY). A Trait inspects the current PlanStep to determine what to do.
_Avoid_: Waypoint, sub-task, goal

**Job**:
A work order in the JobQueue. Has a Type (JobTypeChopTrees, JobTypeMove) and a target cell (TargetX, TargetY). Created by the player via tool actions (e.g., axe tool creates ChopTrees jobs).
_Avoid_: Task, order, work item

**JobQueue**:
A queue that holds Jobs for Villagers to consume. Villagers pull from this queue to determine their movement target. Created by tool actions; consumed by any Idle Villager.
_Avoid_: TaskList, work queue, order book

**JobType**:
The kind of work a Job represents (e.g., JobTypeChopTrees, JobTypeMove). Determines which Plan the Villager constructs when it picks up the Job. The Villager hardcodes the mapping from JobType to a sequence of PlanSteps.
_Avoid_: Job kind, work type, action type

**Lumberjack**:
A Trait that handles the chopping of Trees. Embedded into a Villager. Owns a state machine (Idle, Hitting). Activated when the current PlanStep requires chopping at the Villager's adjacent cell. Each tick while Hitting deals damage to the target Tree. When Tree Health reaches zero, the Lumberjack collects the WoodYield and the Villager enters a Carrying state.
_Avoid_: Woodcutter, forester, chopper

**Camera**:
A 2D viewport that the player controls with right-click drag (pan) and mouse wheel (zoom). Clamped between 25% and 700%.
_Avoid_: Viewport, view

**SpriteBank**:
A package-level variable that owns loaded textures and makes them available to any entity that needs to draw. Exposes `Terrain` and `Human` textures. Created by the `spritebank` package, lifecycle managed via `LoadAll()` / `UnloadAll()`.
_Avoid_: TextureManager, AssetRegistry, resource cache

**Tool**:
A mode the Player can activate to change what happens when they drag-select Cells. Pressing '1' toggles between ToolSelect (highlight cells only) and ToolAxe (highlight + create ChopTrees Jobs).
_Avoid_: Mode, weapon, item

**Console**:
An in-game command console opened by pressing the backtick (`) key. Accepts text input for developer commands. Supported commands: `help`, `spawnvillager <name> <x> <y>`, `addtree <x> <y>`, `removetree <x> <y>`, `cleartrees`, `addjob <move|chop> <x> <y>`, `clearjobs`. Lives in the `ui` package.
_Avoid_: Terminal, shell, CLI

**Debug**:
A flag-based debug system. Press F5 to toggle debug mode on/off. When active, number keys 0-6 toggle individual debug categories: 0=Sim, 2=Move, 3=Path, 4=Clock, 5=Job, 6=World. Each category prints diagnostic messages to stdout when its subsystem runs. Lives in the `debug` package.
_Avoid_: Logging, tracing, verbose

**Trait**:
A composable unit of behavior and state embedded into an Entity. Each trait owns its own logic and lifecycle. Traits are mixed into entities via Go struct embedding. An entity may compose zero or more traits.
_Avoid_: Component, module, plugin, system

**Movement**:
A trait that handles grid-based movement. Owns position (X, Y), target (TargetX, TargetY), path (Waypoints), and a state machine (Idle, Moving, Waiting, Arrived). Exposes `Update(world)` which advances one tick of movement, `SetTarget(x, y, world)` which initiates pathfinding toward a destination, and `Pos()` which returns current coordinates. When the next waypoint is occupied, Movement enters the Waiting state for up to `WaitDuration` (5 ticks) per retry, with a maximum of `MaxRetries` (10) before emitting `EventStuck`. Uses A* pathfinding via the pathfinding package and manages Occupy/Vacate on the World as it moves. Meant to be embedded in any mobile entity (Villager, future Animals, Vehicles).
_Avoid_: Navigation, locomotion, mover

**MovementEvent**:
A value returned by an Entity's `Tick()` to signal what happened during the tick. Used by the Simulation to decide what to do next. Possible values: `EventNone` (no event, still moving/waiting), `EventIdle` (entity is idle and can accept a job), `EventArrived` (entity reached its target), `EventStuck` (entity failed to find a path after repeated retries).
_Avoid_: Signal, message, result

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
A serialized snapshot of the whole Game state written to disk as a JSON file. Contains the World grid, Villagers (with movement state), Trees, Jobs, and Camera position/zoom. Created by pressing F9.
_Avoid_: Save file, save data, save slot

**Load**:
The act of reading a Save from disk and replacing the current Game state entirely with the reconstructed state. Triggered by pressing F10. If the file is missing or corrupt, an error is printed and the game continues unchanged.
_Avoid_: Restore, open, import

**Game**:
The top-level struct that wires together the Simulation, UI, and Clock. Created by `game.New()` (procedural world) or `game.NewFromSave(save)` (from a save file). The main loop calls `UI.Input()`, `Update()` (advances the Clock and ticks the Simulation), and `UI.Draw()` each frame.
_Avoid_: App, engine, state
