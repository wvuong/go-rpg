# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
go run .                          # run the game (must be from repo root — see Asset loading)
go run ./cli/sprites              # LPC spritesheet viewer: ↑↓ sheet, ←→ animation, Space reset
go run ./cli/splitter -file=x.png -x=0 -y=0 -width=64 -height=64 -frames=8 -out=walk
go build ./... && go vet ./...    # cgo deprecation warnings from ebiten's metal driver are expected noise

go test ./...                     # only engine and scene have tests; the rest report "no test files"
go test ./scene/ -run TestAdvanceStepsBlocksOnTimer -v   # single test
```

`cli/splitter` writes one PNG per frame as `<out>-000.png`, and fails rather than emitting blank frames if the requested rect runs past the edge of the sheet.

Test coverage is deliberately narrow: it covers the two places where a bug is invisible until it panics in-game — tile map bounds ([engine/map_test.go](engine/map_test.go)) and the battle step queue ([scene/battlescene_test.go](scene/battlescene_test.go)). Anything touching `ebiten.Image` decoding or `ebitenui` widgets needs the game loop and isn't unit-testable as structured, so scene *construction* is untested. The battle tests get around this by building a `BattleScene` literal with no battlers and no director, which doubles as the assertion: the guards under test must run before either field is touched, so deleting a guard panics instead of silently regressing.

The VS Code launch config (`Launch main.go`) debugs `main.go`.

In-game controls: `Escape` returns to title, `Tab` toggles debug overlays (bounding boxes, tile grid, tile IDs, coordinates). In the battle scene, `1`–`4` force battler states, `A` runs the attack sequence, `R` fires an arrow.

## Architecture

Ebiten game built on a scene/director pattern. Data flows one way down: `main.go` builds `GameConfig` + `GameState` → [game/game.go](game/game.go) implements ebiten's `Update`/`Draw`/`Layout` and delegates everything to [scene/director.go](scene/director.go), which holds exactly one active `Scene` (`Update()` + `Draw(screen)`).

**Scene switching is construction.** `Director.SwitchToX()` builds a brand-new scene each call — nothing is cached. Scene constructors are where sprite sheets get sliced into `SpriteIndex` animations (see [scene/battlescene.go](scene/battlescene.go)), so re-entering a scene rebuilds all of that. The only thing that survives a switch is `engine.GameState`, which every scene holds by pointer; the player mutates `state.WorldPosition` in place, which is how world position persists across title → world → title.

**Package layering.** `engine` is the reusable core and imports nothing from the game. `scene` composes engine types with `ebitenui` widgets. `assets` embeds/loads files and depends on `engine` for the types it constructs (`TileMap`, `SpriteIndex`). Keep `engine` free of imports from `scene`/`game`/`assets`.

### Time and animation

`engine.Timer` converts a `time.Duration` into a tick count via `ebiten.TPS()` and counts up on `Update()`. Everything animated is driven by this, not by wall-clock time — so per-frame durations are approximations rounded to ticks.

`engine.SpriteIndex` is the single animation primitive: a slice of `SubImage` references into one sheet, a per-frame interval list, and a cursor advanced by `NextFrame()` / `PreviousFrame()`. `PreviousFrame()` exists so a run animation can be played backwards for a retreat (`Battler.UseForwardAnimation(false)`). Non-looping indexes clamp on the last frame rather than wrapping. `skipFirstFrame` handles LPC sheets whose row-0 frame is a standing pose that shouldn't appear in a walk cycle.

Animations are grouped into role-specific structs rather than a generic map: `DirectionalSpriteIndex` (world movement), `LpcSpriteIndex` (LPC walk/death), `BattleAnimations` (idle/run/attack/guard). `Battler.Update()` switches on `BattlerState` and picks the matching field — nil animations are skipped, so a battler without a guard sheet simply won't change frames.

`engine.Sprite` treats `ScreenPosition` as the sprite **center**; `Draw` applies horizontal flip, then rotation about the center, then translates by `position - half extents`. Anything computing screen offsets must account for that centering.

### Tile maps

Tiled `.tmx` files in [assets/levels/](assets/levels/) are parsed once at startup by the loader in [assets/assets.go](assets/assets.go) and flattened into `engine.TileMap`:

- `tileImages` is indexed by **global tile ID**, sized `totalTiles+1` because GID 0 means "empty" and must remain a nil slot.
- `layers` is `[][]int` of GIDs, each layer a flat `width*height` array indexed `x + y*width`.
- `impassable` is a single flattened collision mask derived from the Tiled tile property `terrain-impassable`, and is only collected for **layers after index 0** (layer 0 is treated as ground and never blocks).

Adding collidable terrain means setting `terrain-impassable` on the tile in Tiled, not editing Go code. The loader logs the full tile grid and collision map to stdout on every launch — verbose but useful when a level looks wrong. Note it uses `fmt.Printf` for the grids and `log` for everything else, so redirecting `log` alone won't silence it.

Anything off the map reads as **solid** (`isSolidTileAtXY`), so the player is stopped at the edge rather than indexing outside the flattened mask. Keep that guard: the mask is a single flat `[]int`, so an unchecked negative column doesn't fail loudly, it wraps onto the previous row and silently tests the wrong tile.

`engine.Player` does the movement + collision resolution: move, clamp to map bounds, then test the four corners of the bounding box and snap back to the tile edge based on movement direction. The clamp has to come before the bounding box is derived, and it insets by half the sprite because `Position` is the sprite *centre* — clamping to `[0, mapWidth]` would still leave the box hanging half a sprite off the map. `engine.Camera` follows the player, clamps to map bounds, and only recomputes the player's *screen* position near map edges where centering is impossible.

Collision resolution is a single `if/else if` chain over `dirY` then `dirX`, so a diagonal collision only resolves one axis; diagonal movement is also un-normalized and therefore √2 faster than cardinal.

### Battle sequencing

Battle actions are a queue of `step{timer, action}` in `BattleScene.steps`, drained **one step per tick** by `advanceSteps()`. A step with no timer executes immediately; a step with a timer blocks the queue until it's ready. Multi-tick motion is expressed as N identical steps (`ticksNeeded = distance / speed`) rather than as velocity integration, so movement is exact and interruption-free. Build new battle actions by appending steps in the order they should occur — see `Attack()` and `FireArrow()`.

**A sequence's trailing steps are what restore state**, so they can't be dropped: `Attack()` ends by snapping the attacker back to its start X and returning it to idle, and `FireArrow()` ends by removing the arrow from `bs.arrows`. Both therefore return early when `busy()` — a new action replacing a running one would strand a battler mid-field or leak a permanently-drawn arrow. Any new battle action must follow the same rule.

### Asset loading

Two mechanisms coexist, and the difference matters:

- Most assets use `//go:embed *` on `assetsFS` in [assets/assets.go](assets/assets.go) — path-independent.
- LPC spritesheets (`mustLoadLpcSprites`) use `os.DirFS("assets")` and `os.ReadDir("assets/lpc")`, so **the game and `cli/sprites` must be run with the repo root as the working directory** or they panic.

LPC sheets are described declaratively by [assets/lpc/spritesheet.json](assets/lpc/spritesheet.json) (`row`, `width`, `height`, `frames`, `type`, `facing`, `loop`, `skipFirstFrame`). Every `*-spritesheet.png` in `assets/lpc/` is assumed to share that same row layout, so adding a character is a drop-in file — no code change.

Loaded sheets land in `assets.LpcSpritesMap`, keyed by the **trimmed** name (`knight`, not `knight-spritesheet.png`) so that `assets.LpcSpriteNames` can be used to iterate or look them up. Keep those two in step when changing the loader.

`cli/sprites` reimplements the same descriptor walk as `mustLoadLpcSprites`; changes to the parsing rules need to land in both. They have already drifted — the viewer uses 150ms frame intervals and the game 100ms, so the viewer does not preview what the game actually plays.

Tiny Swords assets are the opposite: each animation is an explicitly named package-level `*ebiten.Image` var, with frame counts and sheet dimensions hardcoded at the `NewHorizontalSpriteIndex` call sites in `battlescene.go` (192×192 for warrior/archer, 320×320 for lancer).

## Known rough edges

[TODO.md](TODO.md) tracks the current work list. Beyond it, these are known and unfixed:

- **Diagonal movement** is un-normalized and only resolves one axis on collision (see Tile maps above). `engine.Vector.Normalize` exists for this and is unused — it also has no zero-magnitude guard, so it returns NaN for a zero vector.
- **`filepath.Join` is used for `io/fs` paths** in [assets/assets.go](assets/assets.go) (the tilemap glob and the tileset image load). `io/fs` requires forward slashes, so both break on Windows; use `path.Join`.
- **`MustLoadJson` returns an `error`** despite the `Must` prefix, and its callers `panic` with a string that discards the actual error.
- **`SwitchToTileMap` indexes `tileMaps["levels/world.tmx"]` unchecked** — a missing key yields a nil `*TileMap` that panics later inside `NewPlayer`, far from the cause.
- **`tileImages` is sized `sum(TileCount)+1`**, which assumes Tiled assigns `firstgid` contiguously. True for the current level (max GID 67, len 68) but a gap would overflow the slice.
- **`SpriteIndex.Reset()` doesn't reset its timer**, so the first frame after a reset can last ~0 ticks; `PreviousFrame()` also ignores `skipFirstFrame`, so a reversed walk can land on the standing pose.
- **`BattleScene.ChangeState` hardcodes battler indices** `0,2,1,3` and panics if the roster shrinks.
- **Dead code**: `engine.Ticker` (unused entirely), `engine.ClampInt`, `Battler.ResetAnimation`, `BattleScene.arrowSprite`.
- **Startup cost**: the tile grid dump (above) and `DefaultFont()` re-parsing the goregular TTF on every call.
