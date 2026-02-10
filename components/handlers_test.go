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
			action:   Get("/api/users"),
			expected: "@get('/api/users')",
		},
		{
			name:     "get with debounce",
			action:   Get("/api/search").Debounce(300),
			expected: "__debounce_300ms: @get('/api/search')",
		},
		{
			name:     "post",
			action:   Post("/api/users"),
			expected: "@post('/api/users')",
		},
		{
			name:     "put with format args",
			action:   Put("/api/users/%d", 123),
			expected: "@put('/api/users/123')",
		},
		{
			name:     "delete",
			action:   Delete("/api/users/%d", 456),
			expected: "@delete('/api/users/456')",
		},
		{
			name:     "raw action",
			action:   Raw("$count++"),
			expected: "$count++",
		},
		{
			name:     "raw with debounce",
			action:   Raw("$search = $evt.target.value").Debounce(150),
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
			action:    Get("/api/search"),
			expected:  "$search.length >= 2 && @get('/api/search')",
		},
		{
			name:      "conditional with debounce",
			condition: "$query !== ''",
			action:    Get("/api/autocomplete").Debounce(300),
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
			actions:  []Action{Raw("$loading = true"), Post("/api/save")},
			expected: "$loading = true; @post('/api/save')",
		},
		{
			name:     "three actions",
			actions:  []Action{Raw("$loading = true"), Post("/api/save"), Raw("$loading = false")},
			expected: "$loading = true; @post('/api/save'); $loading = false",
		},
		{
			name:     "single action",
			actions:  []Action{Get("/api/data")},
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
	action := Raw("$loading = true").Then(Post("/api/save")).Then(Raw("$loading = false"))
	expected := "$loading = true; @post('/api/save'); $loading = false"

	if got := action.String(); got != expected {
		t.Errorf("Then().String() = %q, want %q", got, expected)
	}
}

func TestIfElse(t *testing.T) {
	action := IfElse("$active", Raw("$active = false"), Raw("$active = true"))
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
			option: OnClick(Post("/api/delete")),
			event:  "click",
			action: "@post('/api/delete')",
		},
		{
			name:   "OnChange",
			option: OnChange(Get("/api/validate")),
			event:  "change",
			action: "@get('/api/validate')",
		},
		{
			name:   "OnInput with debounce",
			option: OnInput(Get("/api/search").Debounce(300)),
			event:  "input",
			action: "__debounce_300ms: @get('/api/search')",
		},
		{
			name:   "OnSubmit",
			option: OnSubmit(Post("/api/form")),
			event:  "submit",
			action: "@post('/api/form')",
		},
		{
			name:   "OnLoad",
			option: OnLoad(Get("/api/init")),
			event:  "load",
			action: "@get('/api/init')",
		},
		{
			name:   "OnBlur",
			option: OnBlur(Get("/api/validate")),
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
