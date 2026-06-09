# RFC: Restructure Crush Model Configuration

## Status

**Draft** — Initial proposal for discussion.

## Problem

The current Crush config has a confusing two-layer model system that doesn't align with how users think about agents and models:

### Current Structure

```json
{
  "providers": {
    "openai": {
      "models": [
        { "id": "gpt-4o", "name": "GPT-4o", "context_window": 128000 }
      ]
    }
  },
  "models": {
    "large": { "provider": "openai", "model": "gpt-4o" },
    "small": { "provider": "openai", "model": "gpt-4o-mini" }
  }
}
```

**Sources of confusion:**

1. **`models` appears twice** — once under `providers` (a flat array of model definitions with metadata like `id`, `name`, `context_window`) and once at root level (a map of model *references* using `provider` + `model` fields). Users naturally expect one place to define models.

2. **`large`/`small` don't map to agents** — The two built-in agents are `"coder"` and `"task"` (defined in `config.go:60-61`), but the config uses `"large"` and `"small"` as keys. Users expect `models.coder` and `models.task`.

3. **The `model` field is ambiguous** — In `SelectedModel`, the `model` field refers to the model's `id` in the provider's models array, but there's no `id` field on `SelectedModel` itself. Users might think `model` is the model name or type.

4. **No extensibility** — Adding a third agent (e.g., `"researcher"`) requires understanding the `SelectedModelType` enum and the `large`/`small` distinction. The config doesn't naturally support custom agents.

5. **Auto-selection is opaque** — When no `models` section is defined, Crush auto-selects the first model from the first enabled provider (see `load.go:580-597`). This works but is hidden — users don't understand why their model was chosen or how to override it.

## Current Architecture

### Model Resolution Chain

1. `coordinator.go:671` — Looks up `Models["large"]` and `Models["small"]`
2. `config.go:631` — `GetModelByType()` resolves to `GetModel(provider, model)`
3. `config.go:609` — `GetModel()` searches `providerConfig.Models[]` for `m.ID == model`
4. If not found, returns `nil` → UI shows `0%` usage, `"via <provider>"` with no model name

### Agent Configuration

```go
// config.go:60-61
const (
    AgentCoder string = "coder"
    AgentTask  string = "task"
)

// config.go:504
type Agent struct {
    Model SelectedModelType `json:"model"` // "large" or "small"
}
```

Both `"coder"` and `"task"` agents default to `SelectedModelTypeLarge` (`"large"`). The `Agent.Model` field is a reference to a model type, not a direct model reference.

### Auto-Selection (Fallback)

When no `models` section is defined:
- `load.go:580-597` — picks the first enabled provider, first model in its `models[]` array
- `load.go:600-620` — persists these defaults to the config store
- The user gets `models.large` and `models.small` auto-populated without realizing it

## Proposal

### New Structure

```json
{
  "providers": {
    "openai": {
      "models": [
        { "id": "gpt-4o", "name": "GPT-4o", "context_window": 128000 },
        { "id": "gpt-4o-mini", "name": "GPT-4o Mini", "context_window": 128000 }
      ]
    }
  },
  "models": {
    "coder": { "provider": "openai", "model": "gpt-4o" },
    "task": { "provider": "openai", "model": "gpt-4o-mini" }
  }
}
```

### Key Changes

1. **`models.large` → `models.coder`** — The key matches the agent ID (`config.AgentCoder`). Users configure models by agent, not by size.

2. **`models.small` → `models.task`** — Same logic. The "task" agent is the secondary/smaller agent by convention.

3. **Backward compatibility** — Accept both `"large"`/`"small"` and `"coder"`/`"task"` keys, mapping them internally:
   ```
   "large"  → "coder"
   "small"  → "task"
   ```

4. **Auto-selection with explicit defaults** — When no `models` section is defined:
   - Auto-populate `models.coder` and `models.task` with the first provider's first model
   - Write these defaults to the config (or a separate defaults file) so users can see and edit them

5. **Extensible agent model references** — The `Agent` struct's `Model` field should reference a model key directly:
   ```go
   type Agent struct {
       Model string `json:"model"` // key into models map, e.g., "coder", "task", "researcher"
   }
   ```
   This allows custom agents to reference any model key without adding new `SelectedModelType` constants.

6. **Remove the `SelectedModelType` enum** — Replace with string keys that map into the `models` map. The `large`/`small` distinction becomes a convention, not a type.

### Migration Path

1. **Phase 1 (non-breaking):** Accept both `"large"`/`"small"` and `"coder"`/`"task"` as keys in the `models` map. Internally normalize to `"coder"`/`"task"`.

2. **Phase 2 (migration):** When loading a config with `"large"`/`"small"`, log a deprecation warning and auto-migrate to `"coder"`/`"task"`.

3. **Phase 3 (breaking):** Remove support for `"large"`/`"small"` keys.

### Benefits

- **Intuitive config** — `models.coder` matches the agent that uses it
- **Extensible** — Adding a `"researcher"` agent only requires adding `"models.researcher"`
- **Clearer model references** — The `model` field in `SelectedModel` already refers to the model's `id` in the provider's `models[]` array. This doesn't change.
- **Backward compatible** — Existing configs continue to work

## Open Questions

1. **Should `models.task` default to the same model as `models.coder` if not specified?** — This would simplify configs for users who only have one model.

2. **Should we support model aliases?** — e.g., `"models.coder": "gpt-4o"` instead of `{"provider": "openai", "model": "gpt-4o"}` when there's only one provider.

3. **How to handle the auto-selection fallback?** — Should it write defaults to disk, or keep them in-memory only?

4. **Should we add a `default` model key** that acts as a fallback for any agent/model reference that doesn't have an explicit config?

## Implementation Notes

### Files to Change

- `internal/config/config.go` — `SelectedModelType` constants, `Agent.Model` field type
- `internal/config/load.go` — Model selection and auto-default logic
- `internal/agent/coordinator.go` — `buildAgentModels()` lookup
- `internal/ui/model/header.go` — Model display (may need updates for new model lookup)
- `internal/ui/common/elements.go` — Token usage display

### Key Code Paths

- `coordinator.go:671` — `c.cfg.Config().Models[config.SelectedModelTypeLarge]`
- `config.go:631` — `GetModelByType(modelType SelectedModelType)`
- `load.go:600` — `configureSelectedModels()`
- `load.go:540` — `defaultModelSelection()`
