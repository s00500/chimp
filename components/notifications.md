# Notifications

Stacking toast notifications with auto-dismiss and manual close support.

## Setup

The `NotificationContainer` is already included in `BaseHTML` layout. If using a custom layout, add it to your body:

```go
@components.NotificationContainer()
```

## Sending Notifications

Use the helper functions in your HTTP handlers:

```go
func HandleSave(w http.ResponseWriter, r *http.Request) {
    sse := datastar.NewSSE(w, r)

    // Save logic...
    if err != nil {
        components.SendError(sse, "Failed to save")
        return
    }

    components.SendSuccess(sse, "Saved successfully")
}
```

### Available Helpers

| Function | Description |
|----------|-------------|
| `SendSuccess(sse, text)` | Green success notification |
| `SendError(sse, text)` | Red error notification |
| `SendWarning(sse, text)` | Orange warning notification |
| `SendInfo(sse, text)` | Blue info notification |
| `SendNotification(sse, type, text)` | Generic with custom type |

### Notification Types

```go
components.NotificationSuccess  // "success"
components.NotificationError    // "error"
components.NotificationWarning  // "warning"
components.NotificationInfo     // "info"
```

## Direct Usage

If you need more control, use the components directly:

```go
sse.PatchElementTempl(
    components.Notification(components.NotificationSuccess, "Custom message"),
    datastar.WithModeAppend(),
    datastar.WithSelector("#notifications"),
)
```

Or the convenience wrappers:

```go
sse.PatchElementTempl(
    components.SuccessNotification("Saved!"),
    datastar.WithModeAppend(),
    datastar.WithSelector("#notifications"),
)
```

## Behavior

- Notifications stack vertically in the top-right corner
- Slide-in animation from right on appear
- Close button to dismiss manually
- Multiple notifications can be visible simultaneously
- z-index of 9999 ensures visibility above other elements

## Why This Uses Explicit Selectors

Unlike DataTable helpers which rely on idiomorph (Datastar's default outer merge mode matching elements by ID), notifications use **append** mode to stack multiple toasts into `#notifications`. Append mode requires an explicit selector. See `datatable.md` for more on the idiomorph pattern.
