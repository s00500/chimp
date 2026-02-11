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

func TestEventHandler_Modifiers(t *testing.T) {
	tests := []struct {
		name          string
		handler       EventHandler
		expectedEvent string
		expectedAction string
	}{
		{
			name:          "simple click",
			handler:       OnClick(PostSSE("/api/delete")),
			expectedEvent: "click",
			expectedAction: "@post('/api/delete')",
		},
		{
			name:          "click with debounce",
			handler:       OnClick(PostSSE("/api/save")).Debounce(300),
			expectedEvent: "click__debounce.300ms",
			expectedAction: "@post('/api/save')",
		},
		{
			name:          "input with debounce",
			handler:       OnInput(GetSSE("/api/search")).Debounce(300),
			expectedEvent: "input__debounce.300ms",
			expectedAction: "@get('/api/search')",
		},
		{
			name:          "submit with prevent",
			handler:       OnSubmit(PostSSE("/api/form")).Prevent(),
			expectedEvent: "submit__prevent",
			expectedAction: "@post('/api/form')",
		},
		{
			name:          "keydown with window",
			handler:       OnKeydown(RawAction("$foo = 'bar'")).Window(),
			expectedEvent: "keydown__window",
			expectedAction: "$foo = 'bar'",
		},
		{
			name:          "multiple modifiers",
			handler:       OnClick(PostSSE("/api")).Window().Debounce(500).Once(),
			expectedEvent: "click__window__debounce.500ms__once",
			expectedAction: "@post('/api')",
		},
		{
			name:          "scroll with throttle",
			handler:       OnScroll(RawAction("$scrollY = evt.target.scrollTop")).Throttle(100),
			expectedEvent: "scroll__throttle.100ms",
			expectedAction: "$scrollY = evt.target.scrollTop",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := tt.handler.toOption()
			if opt.event != tt.expectedEvent {
				t.Errorf("event = %q, want %q", opt.event, tt.expectedEvent)
			}
			if opt.action != tt.expectedAction {
				t.Errorf("action = %q, want %q", opt.action, tt.expectedAction)
			}
		})
	}
}

func TestKeyCondition(t *testing.T) {
	tests := []struct {
		name     string
		key      Key
		expected string
	}{
		{
			name:     "simple key",
			key:      KeyEnter,
			expected: "evt.key === 'Enter'",
		},
		{
			name:     "escape key",
			key:      KeyEscape,
			expected: "evt.key === 'Escape'",
		},
		{
			name:     "arrow key",
			key:      KeyArrowDown,
			expected: "evt.key === 'ArrowDown'",
		},
		{
			name:     "ctrl modifier",
			key:      KeyS.Ctrl(),
			expected: "evt.ctrlKey && evt.key === 's'",
		},
		{
			name:     "alt modifier",
			key:      KeyEnter.Alt(),
			expected: "evt.altKey && evt.key === 'Enter'",
		},
		{
			name:     "shift modifier",
			key:      KeyTab.Shift(),
			expected: "evt.shiftKey && evt.key === 'Tab'",
		},
		{
			name:     "meta modifier",
			key:      KeyS.Meta(),
			expected: "evt.metaKey && evt.key === 's'",
		},
		{
			name:     "ctrlOrCmd modifier",
			key:      KeyS.CtrlOrCmd(),
			expected: "(evt.ctrlKey || evt.metaKey) && evt.key === 's'",
		},
		{
			name:     "multiple modifiers",
			key:      KeyS.Ctrl().Shift(),
			expected: "evt.ctrlKey && evt.shiftKey && evt.key === 's'",
		},
		{
			name:     "ctrl+alt+delete",
			key:      KeyDelete.Ctrl().Alt(),
			expected: "evt.ctrlKey && evt.altKey && evt.key === 'Delete'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.key.Condition(); got != tt.expected {
				t.Errorf("Key.Condition() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestOnKeydownKey(t *testing.T) {
	tests := []struct {
		name   string
		handler EventHandler
		event  string
		action string
	}{
		{
			name:   "enter key",
			handler: OnKeydownKey(KeyEnter, PostSSE("/api/submit")),
			event:  "keydown",
			action: "evt.key === 'Enter' && @post('/api/submit')",
		},
		{
			name:   "escape key",
			handler: OnKeydownKey(KeyEscape, RawAction("$open = false")),
			event:  "keydown",
			action: "evt.key === 'Escape' && $open = false",
		},
		{
			name:   "ctrl+s with prevent",
			handler: OnKeydownKey(KeyS.Ctrl(), PostSSE("/api/save")).Prevent(),
			event:  "keydown__prevent",
			action: "evt.ctrlKey && evt.key === 's' && @post('/api/save')",
		},
		{
			name:   "ctrl+shift+s",
			handler: OnKeydownKey(KeyS.Ctrl().Shift(), PostSSE("/api/save-as")),
			event:  "keydown",
			action: "evt.ctrlKey && evt.shiftKey && evt.key === 's' && @post('/api/save-as')",
		},
		{
			name:   "arrow navigation",
			handler: OnKeydownKey(KeyArrowDown, RawAction("$selectedIndex++")),
			event:  "keydown",
			action: "evt.key === 'ArrowDown' && $selectedIndex++",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := tt.handler.toOption()
			if opt.event != tt.event {
				t.Errorf("event = %q, want %q", opt.event, tt.event)
			}
			if opt.action != tt.action {
				t.Errorf("action = %q, want %q", opt.action, tt.action)
			}
		})
	}
}

func TestOnKeydownWindow(t *testing.T) {
	handler := OnKeydownWindow(KeyS.CtrlOrCmd(), PostSSE("/api/save")).Prevent()
	opt := handler.toOption()

	expectedEvent := "keydown__window__prevent"
	if opt.event != expectedEvent {
		t.Errorf("event = %q, want %q", opt.event, expectedEvent)
	}

	expectedAction := "(evt.ctrlKey || evt.metaKey) && evt.key === 's' && @post('/api/save')"
	if opt.action != expectedAction {
		t.Errorf("action = %q, want %q", opt.action, expectedAction)
	}
}

func TestElementOptions(t *testing.T) {
	// Test that element options apply correctly
	config := applyElementOptions([]ElementOption{
		WithID("test-div"),
		WithClass("container"),
		OnClick(PostSSE("/api/click")),
		OnKeydownWindow(KeyEscape, RawAction("$open = false")),
		WithShow("$visible"),
	})

	if config.ID != "test-div" {
		t.Errorf("ID = %q, want %q", config.ID, "test-div")
	}
	if config.Class != "container" {
		t.Errorf("Class = %q, want %q", config.Class, "container")
	}
	if config.Datastar.Show != "$visible" {
		t.Errorf("Show = %q, want %q", config.Datastar.Show, "$visible")
	}

	// Check that event handlers are set
	attrs := config.CommonAttrs()
	if attrs["data-on:click"] != "@post('/api/click')" {
		t.Errorf("click handler = %q, want %q", attrs["data-on:click"], "@post('/api/click')")
	}
	expectedKeydown := "evt.key === 'Escape' && $open = false"
	if attrs["data-on:keydown__window"] != expectedKeydown {
		t.Errorf("keydown handler = %q, want %q", attrs["data-on:keydown__window"], expectedKeydown)
	}
}
