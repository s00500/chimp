package components

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/starfederation/datastar-go/datastar"
)

// ============================================================================
// Notification Helpers
// ============================================================================

// SendNotification sends a notification via SSE with append mode to #notifications.
func SendNotification(sse *datastar.ServerSentEventGenerator, msgType string, text string) error {
	return sse.PatchElementTempl(
		Notification(msgType, text),
		datastar.WithModeAppend(),
		datastar.WithSelector("#notifications"),
	)
}

// SendError sends an error notification.
func SendError(sse *datastar.ServerSentEventGenerator, text string) error {
	return SendNotification(sse, NotificationError, text)
}

// SendSuccess sends a success notification.
func SendSuccess(sse *datastar.ServerSentEventGenerator, text string) error {
	return SendNotification(sse, NotificationSuccess, text)
}

// SendWarning sends a warning notification.
func SendWarning(sse *datastar.ServerSentEventGenerator, text string) error {
	return SendNotification(sse, NotificationWarning, text)
}

// SendInfo sends an info notification.
func SendInfo(sse *datastar.ServerSentEventGenerator, text string) error {
	return SendNotification(sse, NotificationInfo, text)
}

// SendAutocompleteResults sends search results to an autocomplete component via SSE.
// It patches the rendered options into the listbox. The <chimp-autocomplete>
// web component automatically detects the DOM change and reinitializes BaseCoat.
// The name parameter must match the name passed to FormAutocomplete.
//
// Example:
//
//	results := []components.SelectOption{
//	    {Value: "1", Label: "Customer A"},
//	    {Value: "2", Label: "Customer B"},
//	}
//	components.SendAutocompleteResults(sse, "customer_id", results)
func SendAutocompleteResults(sse *datastar.ServerSentEventGenerator, name string, results []SelectOption) error {
	selectID := name + "_ac"
	listboxID := "#" + selectID + "_listbox"

	// Patch the rendered options into the listbox.
	// The <chimp-autocomplete> web component's MutationObserver will
	// automatically reinit BaseCoat and preserve popover state.
	return sse.PatchElementTempl(
		autocompleteResults(results),
		datastar.WithModeInner(),
		datastar.WithSelector(listboxID),
	)
}

// ============================================================================
// DataTable Helpers
// ============================================================================

// SendDataTableRows sends rendered table rows to a DataTable component via SSE.
// The rows parameter should be a templ component that renders <tr> elements.
// Rows are wrapped in a <tbody> with the table's body ID so idiomorph matches
// by element ID automatically — no explicit selector or mode needed.
// If the DataTable uses WithSignalPrefix, pass the same prefix as signalPrefix.
//
// Example:
//
//	func handleUsers(w http.ResponseWriter, r *http.Request) {
//	    sse := datastar.NewSSE(w, r)
//	    page, _ := strconv.Atoi(r.URL.Query().Get("page"))
//	    pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
//	    sortField := r.URL.Query().Get("sortField")
//	    sortDir := r.URL.Query().Get("sortDir")
//	    users, total := fetchUsers(page, pageSize, sortField, sortDir)
//	    components.SendDataTableRows(sse, "users", total, UsersRows(users))
//	}
func SendDataTableRows(sse *datastar.ServerSentEventGenerator, id string, totalRows int, rows templ.Component, signalPrefix ...string) error {
	prefix := id
	if len(signalPrefix) > 0 && signalPrefix[0] != "" {
		prefix = signalPrefix[0]
	}

	if err := sse.PatchElementTempl(dataTableBody(id, rows)); err != nil {
		return err
	}

	return sse.MarshalAndPatchSignals(map[string]any{
		prefix: map[string]any{
			"totalRows": totalRows,
			"loading":   false,
		},
	})
}

// DataTableRowID returns the conventional element ID for a table row.
// Use this both in your row templates and in handlers to keep IDs consistent.
//
// Example:
//
//	<tr id={ components.DataTableRowID("users", user.ID) } class="table-row">
func DataTableRowID(tableID string, rowID string) string {
	return tableID + "-row-" + rowID
}

// SendDataTableRow patches a single table row via SSE.
// The row component must render a <tr> with an id matching DataTableRowID(tableID, rowID).
// Idiomorph matches the element by ID automatically — no explicit selector needed.
//
// Example:
//
//	// GET /users/{id}/edit — swap a row into edit mode
//	func handleUserEdit(w http.ResponseWriter, r *http.Request) {
//	    sse := datastar.NewSSE(w, r)
//	    id := chi.URLParam(r, "id")
//	    user := fetchUser(id)
//	    components.SendDataTableRow(sse, "users", id, UserEditRow(user))
//	}
//
//	// PUT /users/{id} — save and swap back to display mode
//	func handleUserSave(w http.ResponseWriter, r *http.Request) {
//	    sse := datastar.NewSSE(w, r)
//	    id := chi.URLParam(r, "id")
//	    // ... save ...
//	    user := fetchUser(id)
//	    components.SendDataTableRow(sse, "users", id, UserRow(user))
//	}
func SendDataTableRow(sse *datastar.ServerSentEventGenerator, tableID string, rowID string, row templ.Component) error {
	return sse.PatchElementTempl(row)
}

// ReadEditSignals reads edit signals from a Datastar request for the given columns.
// Returns a map of field name → new value for each editable column.
func ReadEditSignals(r *http.Request, columns []Column) (map[string]string, error) {
	var allSignals map[string]any
	if err := datastar.ReadSignals(r, &allSignals); err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, col := range columns {
		if col.EditType == "" || col.EditType == EditNone {
			continue
		}
		sigName := EditSignalName(col.Field)
		if val, ok := allSignals[sigName]; ok {
			result[col.Field] = fmt.Sprintf("%v", val)
		}
	}
	return result, nil
}

// PatchEditSignals sends an SSE signal patch to overwrite edit signals with the given values.
// Use this before patching an edit row to prevent stale signals from a previous edit bleeding through
// (data-signals uses merge semantics and won't overwrite existing signals).
func PatchEditSignals(sse *datastar.ServerSentEventGenerator, values map[string]string, columns []Column) error {
	signals := make(map[string]any)
	for _, col := range columns {
		if col.EditType == "" || col.EditType == EditNone {
			continue
		}
		signals[EditSignalName(col.Field)] = values[col.Field]
	}
	return sse.MarshalAndPatchSignals(signals)
}

// ============================================================================
// DataTable Selection Helpers
// ============================================================================

// SelectionSignalName returns the signal name for a row's selection checkbox.
// Example: SelectionSignalName("devices", "dev-001") -> "devices.selectedRows.dev-001"
func SelectionSignalName(prefix, rowID string) string {
	return prefix + ".selectedRows." + rowID
}

// SelectionHighlightExpr returns a Datastar expression for conditional row highlighting.
// Example: SelectionHighlightExpr("devices", "dev-001") -> "$devices.selectedRows['dev-001']"
func SelectionHighlightExpr(prefix, rowID string) string {
	return fmt.Sprintf("$%s.selectedRows['%s']", prefix, rowID)
}

// SelectedCountExpr returns a JS expression that counts selected rows.
// Example: SelectedCountExpr("devices") -> "Object.values($devices.selectedRows).filter(Boolean).length"
func SelectedCountExpr(prefix string) string {
	return fmt.Sprintf("Object.values($%s.selectedRows).filter(Boolean).length", prefix)
}

// ReadSelectedRows reads the selectedRows signal from a Datastar request.
// Returns a map of rowID -> selected (true/false).
func ReadSelectedRows(r *http.Request, prefix string) (map[string]bool, error) {
	var allSignals map[string]any
	if err := datastar.ReadSignals(r, &allSignals); err != nil {
		return nil, err
	}

	result := make(map[string]bool)

	// Look for the nested prefix.selectedRows
	if prefixData, ok := allSignals[prefix]; ok {
		if prefixMap, ok := prefixData.(map[string]any); ok {
			if selectedRows, ok := prefixMap["selectedRows"]; ok {
				if rowsMap, ok := selectedRows.(map[string]any); ok {
					for id, val := range rowsMap {
						if b, ok := val.(bool); ok {
							result[id] = b
						}
					}
				}
			}
		}
	}

	return result, nil
}

// DeselectAllExpr returns a JS expression that clears all row selections client-side.
// Use this in "deselect all" buttons or batch action overlays.
// Example: DeselectAllExpr("devices") -> "Object.keys($devices.selectedRows).forEach(k => $devices.selectedRows[k] = false); $devices.selectAll = false"
func DeselectAllExpr(prefix string) string {
	return fmt.Sprintf("Object.keys($%s.selectedRows).forEach(k => $%s.selectedRows[k] = false); $%s.selectAll = false", prefix, prefix, prefix)
}

// ResetSelectionSignals sends SSE signal patches to clear all selections.
func ResetSelectionSignals(sse *datastar.ServerSentEventGenerator, prefix string) error {
	return sse.MarshalAndPatchSignals(map[string]any{
		prefix: map[string]any{
			"selectedRows": map[string]bool{},
			"selectAll":    false,
		},
	})
}

// ============================================================================
// Datastar Action Builder
// ============================================================================

// Action represents a Datastar action expression.
type Action struct {
	action string
}

// String returns the action string.
func (a Action) String() string {
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
// Conditional Actions
// ============================================================================

// When creates a conditional action that only executes if the condition is true.
// Example: When("$search.length >= 2", GetSSE("/api/search")) -> "$search.length >= 2 && @get('/api/search')"
func When(condition string, action Action) Action {
	return Action{action: fmt.Sprintf("%s && %s", condition, action.action)}
}

// IfElse creates a ternary conditional action.
// Example: IfElse("$active", RawAction("$active = false"), RawAction("$active = true"))
func IfElse(condition string, ifTrue, ifFalse Action) Action {
	return Action{action: fmt.Sprintf("%s ? %s : %s", condition, ifTrue.action, ifFalse.action)}
}

// ============================================================================
// Action Combinators
// ============================================================================

// Then chains multiple actions together with semicolons.
// Example: RawAction("$loading = true").Then(PostSSE("/api/save")) -> "$loading = true; @post('/api/save')"
func (a Action) Then(next Action) Action {
	return Action{action: fmt.Sprintf("%s; %s", a.action, next.action)}
}

// Chain combines multiple actions into one, separated by semicolons.
// Example: Chain(RawAction("$loading = true"), PostSSE("/api/save"), RawAction("$loading = false"))
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
// Event Handler with Modifiers
// ============================================================================

// EventHandler represents a Datastar event handler with modifiers.
// Modifiers are appended to the event name with double underscores.
// Example: click__window__debounce.300ms
type EventHandler struct {
	event  string   // base event name (e.g., "click", "keydown")
	action string   // the action expression
	mods   []string // modifiers (e.g., "window", "debounce.300ms", "prevent")
}

// buildEvent constructs the full event string with modifiers.
func (h EventHandler) buildEvent() string {
	if len(h.mods) == 0 {
		return h.event
	}
	return h.event + "__" + strings.Join(h.mods, "__")
}

// toOption converts the EventHandler to an onOption for use with components.
func (h EventHandler) toOption() onOption {
	return onOption{event: h.buildEvent(), action: h.action}
}

// Window adds the __window modifier (listen on window instead of element).
// Example: OnKeydown(KeyEscape, action).Window() -> data-on:keydown__window="..."
func (h EventHandler) Window() EventHandler {
	h.mods = append(h.mods, "window")
	return h
}

// Debounce adds a debounce modifier with the specified milliseconds.
// Example: OnInput(action).Debounce(300) -> data-on:input__debounce.300ms="..."
func (h EventHandler) Debounce(ms int) EventHandler {
	h.mods = append(h.mods, fmt.Sprintf("debounce.%dms", ms))
	return h
}

// Throttle adds a throttle modifier with the specified milliseconds.
// Example: OnScroll(action).Throttle(100) -> data-on:scroll__throttle.100ms="..."
func (h EventHandler) Throttle(ms int) EventHandler {
	h.mods = append(h.mods, fmt.Sprintf("throttle.%dms", ms))
	return h
}

// Prevent adds the __prevent modifier (calls preventDefault).
// Example: OnSubmit(action).Prevent() -> data-on:submit__prevent="..."
func (h EventHandler) Prevent() EventHandler {
	h.mods = append(h.mods, "prevent")
	return h
}

// Stop adds the __stop modifier (calls stopPropagation).
func (h EventHandler) Stop() EventHandler {
	h.mods = append(h.mods, "stop")
	return h
}

// Once adds the __once modifier (handler fires only once).
func (h EventHandler) Once() EventHandler {
	h.mods = append(h.mods, "once")
	return h
}

// Passive adds the __passive modifier (passive event listener).
func (h EventHandler) Passive() EventHandler {
	h.mods = append(h.mods, "passive")
	return h
}

// Capture adds the __capture modifier (capture phase listener).
func (h EventHandler) Capture() EventHandler {
	h.mods = append(h.mods, "capture")
	return h
}

// Duration adds a duration modifier with the specified milliseconds.
// Used with OnInterval to set the interval timing.
// Example: OnInterval(action).Duration(500) -> data-on:interval__duration.500ms="..."
func (h EventHandler) Duration(ms int) EventHandler {
	h.mods = append(h.mods, fmt.Sprintf("duration.%dms", ms))
	return h
}

// Leading adds the __leading modifier (execute immediately, then at interval).
// Example: OnInterval(action).Duration(500).Leading() -> data-on:interval__duration.500ms__leading="..."
func (h EventHandler) Leading() EventHandler {
	h.mods = append(h.mods, "leading")
	return h
}

// Implement all component option interfaces for EventHandler
func (h EventHandler) applyToForm(c *FormConfig)               { h.toOption().applyToForm(c) }
func (h EventHandler) applyToButton(c *ButtonConfig)           { h.toOption().applyToButton(c) }
func (h EventHandler) applyToAutocomplete(c *AutocompleteConfig) { h.toOption().applyToAutocomplete(c) }
func (h EventHandler) applyToDataTable(c *DataTableConfig)     { h.toOption().applyToDataTable(c) }
func (h EventHandler) applyToStack(c *StackConfig)             { h.toOption().applyToStack(c) }
func (h EventHandler) applyToFormGroup(c *FormGroupConfig)     { h.toOption().applyToFormGroup(c) }
func (h EventHandler) applyToCard(c *CardConfig)               { h.toOption().applyToCard(c) }
func (h EventHandler) applyToSection(c *SectionConfig)         { h.toOption().applyToSection(c) }
func (h EventHandler) applyToElement(c *ElementConfig)         { h.toOption().applyToElement(c) }

// ============================================================================
// Event Handler Constructors
// ============================================================================

// OnClick creates a click event handler.
// Example: OnClick(PostSSE("/api/delete"))
// Example: OnClick(PostSSE("/api/save")).Debounce(300)
func OnClick(action Action) EventHandler {
	return EventHandler{event: "click", action: action.String()}
}

// OnChange creates a change event handler.
// Fires when input value changes and element loses focus.
func OnChange(action Action) EventHandler {
	return EventHandler{event: "change", action: action.String()}
}

// OnInput creates an input event handler.
// Fires on every keystroke/input change.
// Example: OnInput(GetSSE("/api/search")).Debounce(300)
func OnInput(action Action) EventHandler {
	return EventHandler{event: "input", action: action.String()}
}

// OnSubmit creates a submit event handler.
// Example: OnSubmit(PostSSE("/api/form")).Prevent()
func OnSubmit(action Action) EventHandler {
	return EventHandler{event: "submit", action: action.String()}
}

// OnLoad creates a load event handler.
// Fires when the element is loaded/mounted.
func OnLoad(action Action) EventHandler {
	return EventHandler{event: "load", action: action.String()}
}

// DataInit creates a load event handler for initializing data when an element mounts.
// Shortcut for OnLoad — reads more naturally for data-fetching on mount.
// Example: DataInit(GetSSE("/api/users"))
// Example: DataInit(PostSSE("/api/init")).Once()
func DataInit(action Action) EventHandler {
	return OnLoad(action)
}

// OnFocus creates a focus event handler.
func OnFocus(action Action) EventHandler {
	return EventHandler{event: "focus", action: action.String()}
}

// OnBlur creates a blur event handler.
func OnBlur(action Action) EventHandler {
	return EventHandler{event: "blur", action: action.String()}
}

// OnKeydown creates a keydown event handler for any key.
func OnKeydown(action Action) EventHandler {
	return EventHandler{event: "keydown", action: action.String()}
}

// OnKeyup creates a keyup event handler for any key.
func OnKeyup(action Action) EventHandler {
	return EventHandler{event: "keyup", action: action.String()}
}

// OnMouseenter creates a mouseenter event handler.
func OnMouseenter(action Action) EventHandler {
	return EventHandler{event: "mouseenter", action: action.String()}
}

// OnMouseleave creates a mouseleave event handler.
func OnMouseleave(action Action) EventHandler {
	return EventHandler{event: "mouseleave", action: action.String()}
}

// OnScroll creates a scroll event handler.
// Example: OnScroll(action).Throttle(100)
func OnScroll(action Action) EventHandler {
	return EventHandler{event: "scroll", action: action.String()}
}

// OnInterval creates an interval event handler that fires repeatedly.
// Default interval is 1 second.
// Example: OnInterval(GetSSE("/api/poll"))
// Example: OnInterval(GetSSE("/api/poll")).Duration(5000)
// Example: OnInterval(RawAction("$count++")).Duration(500).Leading()
func OnInterval(action Action) EventHandler {
	return EventHandler{event: "interval", action: action.String()}
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
func (k Key) Ctrl() Key {
	k.ctrl = true
	return k
}

// Alt adds the Alt modifier to the key.
func (k Key) Alt() Key {
	k.alt = true
	return k
}

// Shift adds the Shift modifier to the key.
func (k Key) Shift() Key {
	k.shift = true
	return k
}

// Meta adds the Meta (Cmd on Mac, Win on Windows) modifier to the key.
func (k Key) Meta() Key {
	k.meta = true
	return k
}

// CtrlOrCmd adds a cross-platform modifier (Ctrl on Windows/Linux, Cmd on Mac).
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
// Example: OnKeydownKey(KeyS.Ctrl(), PostSSE("/api/save")).Prevent()
//   -> data-on:keydown__prevent="evt.ctrlKey && evt.key === 's' && @post('/api/save')"
func OnKeydownKey(key Key, action Action) EventHandler {
	return EventHandler{
		event:  "keydown",
		action: fmt.Sprintf("%s && %s", key.Condition(), action.String()),
	}
}

// OnKeyupKey creates a keyup handler that only fires for a specific key.
func OnKeyupKey(key Key, action Action) EventHandler {
	return EventHandler{
		event:  "keyup",
		action: fmt.Sprintf("%s && %s", key.Condition(), action.String()),
	}
}

// OnKeydownWindow creates a global keydown handler (listens on window).
// Useful for keyboard shortcuts that should work regardless of focus.
// Example: OnKeydownWindow(KeyS.CtrlOrCmd(), PostSSE("/api/save")).Prevent()
func OnKeydownWindow(key Key, action Action) EventHandler {
	return EventHandler{
		event:  "keydown",
		action: fmt.Sprintf("%s && %s", key.Condition(), action.String()),
		mods:   []string{"window"},
	}
}

// OnKeyupWindow creates a global keyup handler (listens on window).
func OnKeyupWindow(key Key, action Action) EventHandler {
	return EventHandler{
		event:  "keyup",
		action: fmt.Sprintf("%s && %s", key.Condition(), action.String()),
		mods:   []string{"window"},
	}
}

// ============================================================================
// Template Helper Functions
// ============================================================================
// These functions return strings for use directly in templ templates.

// Event creates an event name builder for use in templates.
// Example: Event("input").Debounce(300).String() -> "input__debounce.300ms"
func Event(name string) EventBuilder {
	return EventBuilder{event: name}
}

// EventBuilder builds event names with modifiers for template use.
type EventBuilder struct {
	event string
	mods  []string
}

// String returns the full event name with modifiers.
func (e EventBuilder) String() string {
	if len(e.mods) == 0 {
		return e.event
	}
	return e.event + "__" + strings.Join(e.mods, "__")
}

// Debounce adds a debounce modifier.
func (e EventBuilder) Debounce(ms int) EventBuilder {
	e.mods = append(e.mods, fmt.Sprintf("debounce.%dms", ms))
	return e
}

// Throttle adds a throttle modifier.
func (e EventBuilder) Throttle(ms int) EventBuilder {
	e.mods = append(e.mods, fmt.Sprintf("throttle.%dms", ms))
	return e
}

// Window adds the window modifier.
func (e EventBuilder) Window() EventBuilder {
	e.mods = append(e.mods, "window")
	return e
}

// Prevent adds the prevent modifier.
func (e EventBuilder) Prevent() EventBuilder {
	e.mods = append(e.mods, "prevent")
	return e
}

// Stop adds the stop modifier.
func (e EventBuilder) Stop() EventBuilder {
	e.mods = append(e.mods, "stop")
	return e
}

// Once adds the once modifier.
func (e EventBuilder) Once() EventBuilder {
	e.mods = append(e.mods, "once")
	return e
}

// Attr returns the full attribute name (data-on:event__mods).
func (e EventBuilder) Attr() string {
	return "data-on:" + e.String()
}
