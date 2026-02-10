package components

import (
	"fmt"
	"strings"

	"github.com/starfederation/datastar-go/datastar"
)

// ============================================================================
// Datastar Action Builder
// ============================================================================

// Action represents a Datastar action that can be used in event handlers.
// It supports method chaining for adding modifiers like debounce.
type Action struct {
	action   string
	debounce int // milliseconds, 0 means no debounce
}

// Debounce adds a debounce modifier to the action (in milliseconds).
// Example: Get("/search").Debounce(300) -> "__debounce_300ms: @get('/search')"
func (a Action) Debounce(ms int) Action {
	a.debounce = ms
	return a
}

// String returns the final action string for use in data-on attributes.
func (a Action) String() string {
	if a.debounce > 0 {
		return fmt.Sprintf("__debounce_%dms: %s", a.debounce, a.action)
	}
	return a.action
}

// ============================================================================
// SSE Request Builders (wraps datastar helpers with Action type)
// ============================================================================

// GetSSE creates a GET SSE request action.
// Example: GetSSE("/api/users") -> @get('/api/users')
func GetSSE(urlFormat string, args ...any) Action {
	return Action{action: datastar.GetSSE(urlFormat, args...)}
}

// PostSSE creates a POST SSE request action.
// Example: PostSSE("/api/users") -> @post('/api/users')
func PostSSE(urlFormat string, args ...any) Action {
	return Action{action: datastar.PostSSE(urlFormat, args...)}
}

// PutSSE creates a PUT SSE request action.
// Example: PutSSE("/api/users/%d", id) -> @put('/api/users/123')
func PutSSE(urlFormat string, args ...any) Action {
	return Action{action: datastar.PutSSE(urlFormat, args...)}
}

// PatchSSE creates a PATCH SSE request action.
// Example: PatchSSE("/api/users/%d", id) -> @patch('/api/users/123')
func PatchSSE(urlFormat string, args ...any) Action {
	return Action{action: datastar.PatchSSE(urlFormat, args...)}
}

// DeleteSSE creates a DELETE SSE request action.
// Example: DeleteSSE("/api/users/%d", id) -> @delete('/api/users/123')
func DeleteSSE(urlFormat string, args ...any) Action {
	return Action{action: datastar.DeleteSSE(urlFormat, args...)}
}

// RawAction creates an action from a raw action string.
// Example: RawAction("$count++") -> $count++
func RawAction(action string) Action {
	return Action{action: action}
}

// ============================================================================
// Dedicated Event Handler Options
// ============================================================================

// OnClick creates a click event handler option.
// Example: OnClick(Post("/api/delete/%d", id))
func OnClick(action Action) onOption {
	return onOption{event: "click", action: action.String()}
}

// OnChange creates a change event handler option.
// Fires when input value changes and element loses focus.
// Example: OnChange(Get("/api/validate"))
func OnChange(action Action) onOption {
	return onOption{event: "change", action: action.String()}
}

// OnInput creates an input event handler option.
// Fires on every keystroke/input change.
// Example: OnInput(Get("/api/search").Debounce(300))
func OnInput(action Action) onOption {
	return onOption{event: "input", action: action.String()}
}

// OnSubmit creates a submit event handler option.
// Example: OnSubmit(Post("/api/form"))
func OnSubmit(action Action) onOption {
	return onOption{event: "submit", action: action.String()}
}

// OnLoad creates a load event handler option.
// Fires when the element is loaded/mounted.
// Example: OnLoad(Get("/api/init"))
func OnLoad(action Action) onOption {
	return onOption{event: "load", action: action.String()}
}

// OnFocus creates a focus event handler option.
// Example: OnFocus(Raw("$focused = true"))
func OnFocus(action Action) onOption {
	return onOption{event: "focus", action: action.String()}
}

// OnBlur creates a blur event handler option.
// Example: OnBlur(Get("/api/validate"))
func OnBlur(action Action) onOption {
	return onOption{event: "blur", action: action.String()}
}

// OnKeydown creates a keydown event handler option.
// Example: OnKeydown(Raw("$evt.key === 'Enter' && @post('/api/submit')"))
func OnKeydown(action Action) onOption {
	return onOption{event: "keydown", action: action.String()}
}

// OnKeyup creates a keyup event handler option.
// Example: OnKeyup(Raw("$evt.key === 'Escape' && $open = false"))
func OnKeyup(action Action) onOption {
	return onOption{event: "keyup", action: action.String()}
}

// OnMouseenter creates a mouseenter event handler option.
// Example: OnMouseenter(Raw("$hovered = true"))
func OnMouseenter(action Action) onOption {
	return onOption{event: "mouseenter", action: action.String()}
}

// OnMouseleave creates a mouseleave event handler option.
// Example: OnMouseleave(Raw("$hovered = false"))
func OnMouseleave(action Action) onOption {
	return onOption{event: "mouseleave", action: action.String()}
}

// ============================================================================
// Conditional Actions
// ============================================================================

// When creates a conditional action that only executes if the condition is true.
// Example: When("$search.length >= 2", Get("/api/search")) -> "$search.length >= 2 && @get('/api/search')"
func When(condition string, action Action) Action {
	return Action{
		action:   fmt.Sprintf("%s && %s", condition, action.action),
		debounce: action.debounce,
	}
}

// IfElse creates a ternary conditional action.
// Example: IfElse("$active", Raw("$active = false"), Raw("$active = true"))
func IfElse(condition string, ifTrue, ifFalse Action) Action {
	return Action{
		action: fmt.Sprintf("%s ? %s : %s", condition, ifTrue.action, ifFalse.action),
	}
}

// ============================================================================
// Action Combinators
// ============================================================================

// Then chains multiple actions together with semicolons.
// Example: Raw("$loading = true").Then(Post("/api/save")) -> "$loading = true; @post('/api/save')"
func (a Action) Then(next Action) Action {
	return Action{
		action:   fmt.Sprintf("%s; %s", a.action, next.action),
		debounce: a.debounce, // preserve debounce from first action
	}
}

// Chain combines multiple actions into one, separated by semicolons.
// Example: Chain(Raw("$loading = true"), Post("/api/save"), Raw("$loading = false"))
func Chain(actions ...Action) Action {
	if len(actions) == 0 {
		return Action{}
	}
	result := actions[0]
	for _, next := range actions[1:] {
		result = result.Then(next)
	}
	return result
}

// ============================================================================
// Key Types and Constants
// ============================================================================

// Key represents a keyboard key for use with OnKeydownKey/OnKeyupKey.
// Keys generate conditions like "evt.key === 'Enter'" for use in Datastar expressions.
type Key struct {
	key       string // The evt.key value (e.g., "Enter", "Escape", "s")
	ctrl      bool
	alt       bool
	shift     bool
	meta      bool
	ctrlOrCmd bool // Cross-platform: ctrl on Windows/Linux, meta on Mac
}

// Condition returns the JavaScript condition for matching this key.
// Example: KeyEnter.Condition() -> "evt.key === 'Enter'"
// Example: KeyS.Ctrl().Condition() -> "evt.ctrlKey && evt.key === 's'"
func (k Key) Condition() string {
	var parts []string

	if k.ctrl {
		parts = append(parts, "evt.ctrlKey")
	}
	if k.alt {
		parts = append(parts, "evt.altKey")
	}
	if k.shift {
		parts = append(parts, "evt.shiftKey")
	}
	if k.meta {
		parts = append(parts, "evt.metaKey")
	}
	if k.ctrlOrCmd {
		parts = append(parts, "(evt.ctrlKey || evt.metaKey)")
	}

	parts = append(parts, fmt.Sprintf("evt.key === '%s'", k.key))

	return strings.Join(parts, " && ")
}

// Ctrl adds the Ctrl modifier to the key.
// Example: KeyS.Ctrl() -> evt.ctrlKey && evt.key === 's'
func (k Key) Ctrl() Key {
	k.ctrl = true
	return k
}

// Alt adds the Alt modifier to the key.
// Example: KeyEnter.Alt() -> evt.altKey && evt.key === 'Enter'
func (k Key) Alt() Key {
	k.alt = true
	return k
}

// Shift adds the Shift modifier to the key.
// Example: KeyTab.Shift() -> evt.shiftKey && evt.key === 'Tab'
func (k Key) Shift() Key {
	k.shift = true
	return k
}

// Meta adds the Meta (Cmd on Mac, Win on Windows) modifier to the key.
// Example: KeyS.Meta() -> evt.metaKey && evt.key === 's'
func (k Key) Meta() Key {
	k.meta = true
	return k
}

// CtrlOrCmd adds a cross-platform modifier (Ctrl on Windows/Linux, Cmd on Mac).
// Example: KeyS.CtrlOrCmd() -> (evt.ctrlKey || evt.metaKey) && evt.key === 's'
func (k Key) CtrlOrCmd() Key {
	k.ctrlOrCmd = true
	return k
}

// Common key constants (using standard KeyboardEvent.key values)
var (
	// Navigation keys
	KeyEnter     = Key{key: "Enter"}
	KeyEscape    = Key{key: "Escape"}
	KeyTab       = Key{key: "Tab"}
	KeyBackspace = Key{key: "Backspace"}
	KeyDelete    = Key{key: "Delete"}
	KeySpace     = Key{key: " "}

	// Arrow keys
	KeyArrowUp    = Key{key: "ArrowUp"}
	KeyArrowDown  = Key{key: "ArrowDown"}
	KeyArrowLeft  = Key{key: "ArrowLeft"}
	KeyArrowRight = Key{key: "ArrowRight"}

	// Function keys
	KeyF1  = Key{key: "F1"}
	KeyF2  = Key{key: "F2"}
	KeyF3  = Key{key: "F3"}
	KeyF4  = Key{key: "F4"}
	KeyF5  = Key{key: "F5"}
	KeyF6  = Key{key: "F6"}
	KeyF7  = Key{key: "F7"}
	KeyF8  = Key{key: "F8"}
	KeyF9  = Key{key: "F9"}
	KeyF10 = Key{key: "F10"}
	KeyF11 = Key{key: "F11"}
	KeyF12 = Key{key: "F12"}

	// Common letter keys (lowercase - standard for evt.key with no shift)
	KeyA = Key{key: "a"}
	KeyB = Key{key: "b"}
	KeyC = Key{key: "c"}
	KeyD = Key{key: "d"}
	KeyE = Key{key: "e"}
	KeyF = Key{key: "f"}
	KeyG = Key{key: "g"}
	KeyH = Key{key: "h"}
	KeyI = Key{key: "i"}
	KeyJ = Key{key: "j"}
	KeyK = Key{key: "k"}
	KeyL = Key{key: "l"}
	KeyM = Key{key: "m"}
	KeyN = Key{key: "n"}
	KeyO = Key{key: "o"}
	KeyP = Key{key: "p"}
	KeyQ = Key{key: "q"}
	KeyR = Key{key: "r"}
	KeyS = Key{key: "s"}
	KeyT = Key{key: "t"}
	KeyU = Key{key: "u"}
	KeyV = Key{key: "v"}
	KeyW = Key{key: "w"}
	KeyX = Key{key: "x"}
	KeyY = Key{key: "y"}
	KeyZ = Key{key: "z"}

	// Number keys
	Key0 = Key{key: "0"}
	Key1 = Key{key: "1"}
	Key2 = Key{key: "2"}
	Key3 = Key{key: "3"}
	Key4 = Key{key: "4"}
	Key5 = Key{key: "5"}
	Key6 = Key{key: "6"}
	Key7 = Key{key: "7"}
	Key8 = Key{key: "8"}
	Key9 = Key{key: "9"}

	// Other common keys
	KeyHome     = Key{key: "Home"}
	KeyEnd      = Key{key: "End"}
	KeyPageUp   = Key{key: "PageUp"}
	KeyPageDown = Key{key: "PageDown"}
	KeyInsert   = Key{key: "Insert"}
)

// ============================================================================
// Key-Specific Event Handlers
// ============================================================================

// OnKeydownKey creates a keydown handler that only fires for a specific key.
// Uses evt.key matching as per Datastar documentation.
// Example: OnKeydownKey(KeyEnter, PostSSE("/api/submit"))
//   -> data-on:keydown="evt.key === 'Enter' && @post('/api/submit')"
// Example: OnKeydownKey(KeyS.Ctrl(), PostSSE("/api/save"))
//   -> data-on:keydown="evt.ctrlKey && evt.key === 's' && @post('/api/save')"
func OnKeydownKey(key Key, action Action) onOption {
	return onOption{
		event:  "keydown",
		action: fmt.Sprintf("%s && %s", key.Condition(), action.String()),
	}
}

// OnKeyupKey creates a keyup handler that only fires for a specific key.
// Example: OnKeyupKey(KeyEscape, RawAction("$open = false"))
//   -> data-on:keyup="evt.key === 'Escape' && $open = false"
func OnKeyupKey(key Key, action Action) onOption {
	return onOption{
		event:  "keyup",
		action: fmt.Sprintf("%s && %s", key.Condition(), action.String()),
	}
}

// OnKeydownWindow creates a global keydown handler (listens on window).
// Useful for keyboard shortcuts that should work regardless of focus.
// Example: OnKeydownWindow(KeyS.CtrlOrCmd(), PostSSE("/api/save"))
//   -> data-on:keydown__window="(evt.ctrlKey || evt.metaKey) && evt.key === 's' && @post('/api/save')"
func OnKeydownWindow(key Key, action Action) onOption {
	return onOption{
		event:  "keydown__window",
		action: fmt.Sprintf("%s && %s", key.Condition(), action.String()),
	}
}

// OnKeyupWindow creates a global keyup handler (listens on window).
func OnKeyupWindow(key Key, action Action) onOption {
	return onOption{
		event:  "keyup__window",
		action: fmt.Sprintf("%s && %s", key.Condition(), action.String()),
	}
}
