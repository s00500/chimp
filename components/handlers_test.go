package components

import (
	"testing"
)

func TestAction_String(t *testing.T) {
	tests := []struct {
		name     string
		action   Action
		expected string
	}{
		{
			name:     "simple get",
			action:   GetSSE("/api/users"),
			expected: "@get('/api/users')",
		},
		{
			name:     "get with debounce",
			action:   GetSSE("/api/search").Debounce(300),
			expected: "__debounce_300ms: @get('/api/search')",
		},
		{
			name:     "post",
			action:   PostSSE("/api/users"),
			expected: "@post('/api/users')",
		},
		{
			name:     "put with format args",
			action:   PutSSE("/api/users/%d", 123),
			expected: "@put('/api/users/123')",
		},
		{
			name:     "delete",
			action:   DeleteSSE("/api/users/%d", 456),
			expected: "@delete('/api/users/456')",
		},
		{
			name:     "raw action",
			action:   RawAction("$count++"),
			expected: "$count++",
		},
		{
			name:     "raw with debounce",
			action:   RawAction("$search = $evt.target.value").Debounce(150),
			expected: "__debounce_150ms: $search = $evt.target.value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.action.String(); got != tt.expected {
				t.Errorf("Action.String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestWhen(t *testing.T) {
	tests := []struct {
		name      string
		condition string
		action    Action
		expected  string
	}{
		{
			name:      "conditional get",
			condition: "$search.length >= 2",
			action:    GetSSE("/api/search"),
			expected:  "$search.length >= 2 && @get('/api/search')",
		},
		{
			name:      "conditional with debounce",
			condition: "$query !== ''",
			action:    GetSSE("/api/autocomplete").Debounce(300),
			expected:  "__debounce_300ms: $query !== '' && @get('/api/autocomplete')",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := When(tt.condition, tt.action).String()
			if got != tt.expected {
				t.Errorf("When().String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestChain(t *testing.T) {
	tests := []struct {
		name     string
		actions  []Action
		expected string
	}{
		{
			name:     "two actions",
			actions:  []Action{RawAction("$loading = true"), PostSSE("/api/save")},
			expected: "$loading = true; @post('/api/save')",
		},
		{
			name:     "three actions",
			actions:  []Action{RawAction("$loading = true"), PostSSE("/api/save"), RawAction("$loading = false")},
			expected: "$loading = true; @post('/api/save'); $loading = false",
		},
		{
			name:     "single action",
			actions:  []Action{GetSSE("/api/data")},
			expected: "@get('/api/data')",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Chain(tt.actions...).String()
			if got != tt.expected {
				t.Errorf("Chain().String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestThen(t *testing.T) {
	action := RawAction("$loading = true").Then(PostSSE("/api/save")).Then(RawAction("$loading = false"))
	expected := "$loading = true; @post('/api/save'); $loading = false"

	if got := action.String(); got != expected {
		t.Errorf("Then().String() = %q, want %q", got, expected)
	}
}

func TestIfElse(t *testing.T) {
	action := IfElse("$active", RawAction("$active = false"), RawAction("$active = true"))
	expected := "$active ? $active = false : $active = true"

	if got := action.String(); got != expected {
		t.Errorf("IfElse().String() = %q, want %q", got, expected)
	}
}

func TestEventHandlers(t *testing.T) {
	tests := []struct {
		name     string
		option   onOption
		event    string
		action   string
	}{
		{
			name:   "OnClick",
			option: OnClick(PostSSE("/api/delete")),
			event:  "click",
			action: "@post('/api/delete')",
		},
		{
			name:   "OnChange",
			option: OnChange(GetSSE("/api/validate")),
			event:  "change",
			action: "@get('/api/validate')",
		},
		{
			name:   "OnInput with debounce",
			option: OnInput(GetSSE("/api/search").Debounce(300)),
			event:  "input",
			action: "__debounce_300ms: @get('/api/search')",
		},
		{
			name:   "OnSubmit",
			option: OnSubmit(PostSSE("/api/form")),
			event:  "submit",
			action: "@post('/api/form')",
		},
		{
			name:   "OnLoad",
			option: OnLoad(GetSSE("/api/init")),
			event:  "load",
			action: "@get('/api/init')",
		},
		{
			name:   "OnBlur",
			option: OnBlur(GetSSE("/api/validate")),
			event:  "blur",
			action: "@get('/api/validate')",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.option.event != tt.event {
				t.Errorf("event = %q, want %q", tt.option.event, tt.event)
			}
			if tt.option.action != tt.action {
				t.Errorf("action = %q, want %q", tt.option.action, tt.action)
			}
		})
	}
}

func TestKeyString(t *testing.T) {
	tests := []struct {
		name     string
		key      Key
		expected string
	}{
		{
			name:     "simple key",
			key:      KeyEnter,
			expected: "enter",
		},
		{
			name:     "escape key",
			key:      KeyEscape,
			expected: "escape",
		},
		{
			name:     "arrow key",
			key:      KeyArrowDown,
			expected: "arrowdown",
		},
		{
			name:     "ctrl modifier",
			key:      KeyS.Ctrl(),
			expected: "ctrl.s",
		},
		{
			name:     "alt modifier",
			key:      KeyEnter.Alt(),
			expected: "alt.enter",
		},
		{
			name:     "shift modifier",
			key:      KeyTab.Shift(),
			expected: "shift.tab",
		},
		{
			name:     "meta modifier",
			key:      KeyS.Meta(),
			expected: "meta.s",
		},
		{
			name:     "ctrlormeta modifier",
			key:      KeyS.CtrlOrMeta(),
			expected: "ctrlormeta.s",
		},
		{
			name:     "multiple modifiers",
			key:      KeyS.Ctrl().Shift(),
			expected: "ctrl.shift.s",
		},
		{
			name:     "ctrl+alt+delete",
			key:      KeyDelete.Ctrl().Alt(),
			expected: "ctrl.alt.delete",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.key.String(); got != tt.expected {
				t.Errorf("Key.String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestOnKeydownKey(t *testing.T) {
	tests := []struct {
		name   string
		option onOption
		event  string
		action string
	}{
		{
			name:   "enter key",
			option: OnKeydownKey(KeyEnter, PostSSE("/api/submit")),
			event:  "keydown.enter",
			action: "@post('/api/submit')",
		},
		{
			name:   "escape key",
			option: OnKeydownKey(KeyEscape, RawAction("$open = false")),
			event:  "keydown.escape",
			action: "$open = false",
		},
		{
			name:   "ctrl+s",
			option: OnKeydownKey(KeyS.Ctrl(), PostSSE("/api/save")),
			event:  "keydown.ctrl.s",
			action: "@post('/api/save')",
		},
		{
			name:   "ctrl+shift+s",
			option: OnKeydownKey(KeyS.Ctrl().Shift(), PostSSE("/api/save-as")),
			event:  "keydown.ctrl.shift.s",
			action: "@post('/api/save-as')",
		},
		{
			name:   "arrow navigation",
			option: OnKeydownKey(KeyArrowDown, RawAction("$selectedIndex++")),
			event:  "keydown.arrowdown",
			action: "$selectedIndex++",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.option.event != tt.event {
				t.Errorf("event = %q, want %q", tt.option.event, tt.event)
			}
			if tt.option.action != tt.action {
				t.Errorf("action = %q, want %q", tt.option.action, tt.action)
			}
		})
	}
}

func TestOnKeyupKey(t *testing.T) {
	option := OnKeyupKey(KeyEscape, RawAction("$modal = false"))

	if option.event != "keyup.escape" {
		t.Errorf("event = %q, want %q", option.event, "keyup.escape")
	}
	if option.action != "$modal = false" {
		t.Errorf("action = %q, want %q", option.action, "$modal = false")
	}
}
