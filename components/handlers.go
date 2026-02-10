package components

import (
	"fmt"

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

// Get creates a GET SSE request action.
// Example: Get("/api/users") -> @get('/api/users')
func Get(urlFormat string, args ...any) Action {
	return Action{action: datastar.GetSSE(urlFormat, args...)}
}

// Post creates a POST SSE request action.
// Example: Post("/api/users") -> @post('/api/users')
func Post(urlFormat string, args ...any) Action {
	return Action{action: datastar.PostSSE(urlFormat, args...)}
}

// Put creates a PUT SSE request action.
// Example: Put("/api/users/%d", id) -> @put('/api/users/123')
func Put(urlFormat string, args ...any) Action {
	return Action{action: datastar.PutSSE(urlFormat, args...)}
}

// Patch creates a PATCH SSE request action.
// Example: Patch("/api/users/%d", id) -> @patch('/api/users/123')
func Patch(urlFormat string, args ...any) Action {
	return Action{action: datastar.PatchSSE(urlFormat, args...)}
}

// Delete creates a DELETE SSE request action.
// Example: Delete("/api/users/%d", id) -> @delete('/api/users/123')
func Delete(urlFormat string, args ...any) Action {
	return Action{action: datastar.DeleteSSE(urlFormat, args...)}
}

// Raw creates an action from a raw action string.
// Example: Raw("$count++") -> $count++
func Raw(action string) Action {
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
