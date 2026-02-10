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
type Key struct {
	name      string
	modifiers []string
}

// String returns the key as a dot-separated modifier string for Datastar.
// Example: KeyEnter.Ctrl() -> "ctrl.enter"
func (k Key) String() string {
	if len(k.modifiers) == 0 {
		return k.name
	}
	var sb strings.Builder
	for _, mod := range k.modifiers {
		sb.WriteString(mod)
		sb.WriteByte('.')
	}
	sb.WriteString(k.name)
	return sb.String()
}

// Ctrl adds the Ctrl modifier to the key.
// Example: KeyS.Ctrl() -> Ctrl+S
func (k Key) Ctrl() Key {
	return Key{name: k.name, modifiers: append(k.modifiers, "ctrl")}
}

// Alt adds the Alt modifier to the key.
// Example: KeyEnter.Alt() -> Alt+Enter
func (k Key) Alt() Key {
	return Key{name: k.name, modifiers: append(k.modifiers, "alt")}
}

// Shift adds the Shift modifier to the key.
// Example: KeyTab.Shift() -> Shift+Tab
func (k Key) Shift() Key {
	return Key{name: k.name, modifiers: append(k.modifiers, "shift")}
}

// Meta adds the Meta (Cmd on Mac, Win on Windows) modifier to the key.
// Example: KeyS.Meta() -> Cmd+S / Win+S
func (k Key) Meta() Key {
	return Key{name: k.name, modifiers: append(k.modifiers, "meta")}
}

// CtrlOrMeta adds ctrl on non-Mac and meta on Mac (common for cross-platform shortcuts).
// Example: KeyS.CtrlOrMeta() -> Ctrl+S on Windows/Linux, Cmd+S on Mac
func (k Key) CtrlOrMeta() Key {
	return Key{name: k.name, modifiers: append(k.modifiers, "ctrlormeta")}
}

// Common key constants
var (
	// Navigation keys
	KeyEnter     = Key{name: "enter"}
	KeyEscape    = Key{name: "escape"}
	KeyTab       = Key{name: "tab"}
	KeyBackspace = Key{name: "backspace"}
	KeyDelete    = Key{name: "delete"}
	KeySpace     = Key{name: "space"}

	// Arrow keys
	KeyArrowUp    = Key{name: "arrowup"}
	KeyArrowDown  = Key{name: "arrowdown"}
	KeyArrowLeft  = Key{name: "arrowleft"}
	KeyArrowRight = Key{name: "arrowright"}

	// Function keys
	KeyF1  = Key{name: "f1"}
	KeyF2  = Key{name: "f2"}
	KeyF3  = Key{name: "f3"}
	KeyF4  = Key{name: "f4"}
	KeyF5  = Key{name: "f5"}
	KeyF6  = Key{name: "f6"}
	KeyF7  = Key{name: "f7"}
	KeyF8  = Key{name: "f8"}
	KeyF9  = Key{name: "f9"}
	KeyF10 = Key{name: "f10"}
	KeyF11 = Key{name: "f11"}
	KeyF12 = Key{name: "f12"}

	// Common letter keys (for shortcuts)
	KeyA = Key{name: "a"}
	KeyB = Key{name: "b"}
	KeyC = Key{name: "c"}
	KeyD = Key{name: "d"}
	KeyE = Key{name: "e"}
	KeyF = Key{name: "f"}
	KeyG = Key{name: "g"}
	KeyH = Key{name: "h"}
	KeyI = Key{name: "i"}
	KeyJ = Key{name: "j"}
	KeyK = Key{name: "k"}
	KeyL = Key{name: "l"}
	KeyM = Key{name: "m"}
	KeyN = Key{name: "n"}
	KeyO = Key{name: "o"}
	KeyP = Key{name: "p"}
	KeyQ = Key{name: "q"}
	KeyR = Key{name: "r"}
	KeyS = Key{name: "s"}
	KeyT = Key{name: "t"}
	KeyU = Key{name: "u"}
	KeyV = Key{name: "v"}
	KeyW = Key{name: "w"}
	KeyX = Key{name: "x"}
	KeyY = Key{name: "y"}
	KeyZ = Key{name: "z"}

	// Number keys
	Key0 = Key{name: "0"}
	Key1 = Key{name: "1"}
	Key2 = Key{name: "2"}
	Key3 = Key{name: "3"}
	Key4 = Key{name: "4"}
	Key5 = Key{name: "5"}
	Key6 = Key{name: "6"}
	Key7 = Key{name: "7"}
	Key8 = Key{name: "8"}
	Key9 = Key{name: "9"}

	// Other common keys
	KeyHome     = Key{name: "home"}
	KeyEnd      = Key{name: "end"}
	KeyPageUp   = Key{name: "pageup"}
	KeyPageDown = Key{name: "pagedown"}
	KeyInsert   = Key{name: "insert"}
)

// ============================================================================
// Key-Specific Event Handlers
// ============================================================================

// OnKeydownKey creates a keydown handler that only fires for a specific key.
// Example: OnKeydownKey(KeyEnter, Post("/api/submit"))
// Example: OnKeydownKey(KeyS.Ctrl(), Post("/api/save"))
func OnKeydownKey(key Key, action Action) onOption {
	return onOption{event: "keydown." + key.String(), action: action.String()}
}

// OnKeyupKey creates a keyup handler that only fires for a specific key.
// Example: OnKeyupKey(KeyEscape, Raw("$open = false"))
func OnKeyupKey(key Key, action Action) onOption {
	return onOption{event: "keyup." + key.String(), action: action.String()}
}
