# scriptor

A "scripting" framework, to coordinate actions: For example, this can be useful to construct testing sequences in Go.

# HOW TO USE

## 1. Create a `scene.Scene` object

```go
s := scene.New()
```

## 2. Add `scene.Action` objects to `scene.Scene`

```go
// inline Action
s.Add(scene.ActionFunc(func(ctx context.Context) error {
  // .. do something ...
  return nil
}))

// built-in Action
s.Add(actions.Delay(5*time.Second))

// custom Action object
s.Add(myCustomAction{})
```

## 3. Prepare a `context.Context` object

All transient data in this tool is expected to be passed down along with a `context.Context` object.

For example, logging is done through a `slog.Logger` passed down with the context. This means that the `context.Context` object must be primed with the logger object before the `Action`s are fired.

To do this, you can manually create a context object using the appropriate injector functions, such as `log.InjectContext()`:

```go
ctx := log.InjectContext(context.Background(), slog.New(....))
```

Or, to get the default set of values injected, you can use `scriptor.DefaultContext()`:

```go
ctx := scriptor.DefaultContext(context.Background())
```

The values that need to be injected defer based on the actions that you provide.
As of this writing the documentation is lacking, but in the future each component
in this module should clearly state which values need to be injected into the
`context.Context` object.

As of this writing it is safest to just use `scriptor.DefaultContext()` for most everything.

## 4. Execute the `scene.Scene` object

Finally, put everything together and execute the scene.

```go
s.Execute(ctx)
```