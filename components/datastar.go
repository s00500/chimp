package components

import (
	"strings"

	"github.com/a-h/templ"
)

// DatastarAttrs holds common Datastar attributes that can be applied to components
type DatastarAttrs struct {
	Model   string            // data-model
	On      map[string]string // data-on:event -> action
	Bind    map[string]string // data-bind:attr -> expression
	Signals map[string]string // data-signals
	Attrs   map[string]string // data-attr:name -> expression
	Show    string            // data-show expression
	Text    string            // data-text expression
}

// ToAttrs converts DatastarAttrs to templ.Attributes for use in templates
func (d *DatastarAttrs) ToAttrs() templ.Attributes {
	attrs := templ.Attributes{}

	if d.Model != "" {
		attrs["data-model"] = d.Model
	}

	for event, action := range d.On {
		attrs["data-on:"+event] = action
	}

	for attr, expr := range d.Bind {
		attrs["data-bind:"+attr] = expr
	}

	if len(d.Signals) > 0 {
		// Format as JSON object
		var sb strings.Builder
		sb.WriteString("{")
		first := true
		for name, value := range d.Signals {
			if !first {
				sb.WriteString(", ")
			}
			sb.WriteString(name)
			sb.WriteString(": ")
			sb.WriteString(value)
			first = false
		}
		sb.WriteString("}")
		attrs["data-signals"] = sb.String()
	}

	for name, expr := range d.Attrs {
		attrs["data-attr:"+name] = expr
	}

	if d.Show != "" {
		attrs["data-show"] = d.Show
	}

	if d.Text != "" {
		attrs["data-text"] = d.Text
	}

	return attrs
}

// FormConfig holds configuration for form input components
type FormConfig struct {
	// Basic attributes
	Type        string // input type: text, email, password, number, etc.
	Placeholder string
	Required    bool
	Disabled    bool
	Readonly    bool
	Autofocus   bool

	// Validation
	Min       string // min value for number/date
	Max       string // max value for number/date
	MinLength int
	MaxLength int
	Pattern   string // regex pattern
	Step      string // step for number inputs

	// Textarea specific
	Rows int
	Cols int

	// Select specific
	Options     []SelectOption
	Multiple    bool
	EmptyOption string // placeholder option text

	// Error handling
	Error string // data-show expression for error, or static error text

	// Datastar attributes
	Datastar DatastarAttrs

	// Extra attributes
	Class      string
	ExtraAttrs templ.Attributes
}

// FormOption is a function that modifies FormConfig
type FormOption func(*FormConfig)

// WithType sets the input type (text, email, password, number, etc.)
func WithType(t string) FormOption {
	return func(c *FormConfig) {
		c.Type = t
	}
}

// WithPlaceholder sets the placeholder text
func WithPlaceholder(p string) FormOption {
	return func(c *FormConfig) {
		c.Placeholder = p
	}
}

// WithRequired marks the field as required
func WithRequired() FormOption {
	return func(c *FormConfig) {
		c.Required = true
	}
}

// WithDisabled marks the field as disabled
func WithDisabled() FormOption {
	return func(c *FormConfig) {
		c.Disabled = true
	}
}

// WithReadonly marks the field as readonly
func WithReadonly() FormOption {
	return func(c *FormConfig) {
		c.Readonly = true
	}
}

// WithAutofocus sets autofocus on the field
func WithAutofocus() FormOption {
	return func(c *FormConfig) {
		c.Autofocus = true
	}
}

// WithMin sets the minimum value
func WithMin(min string) FormOption {
	return func(c *FormConfig) {
		c.Min = min
	}
}

// WithMax sets the maximum value
func WithMax(max string) FormOption {
	return func(c *FormConfig) {
		c.Max = max
	}
}

// WithMinLength sets the minimum length
func WithMinLength(n int) FormOption {
	return func(c *FormConfig) {
		c.MinLength = n
	}
}

// WithMaxLength sets the maximum length
func WithMaxLength(n int) FormOption {
	return func(c *FormConfig) {
		c.MaxLength = n
	}
}

// WithPattern sets a regex validation pattern
func WithPattern(pattern string) FormOption {
	return func(c *FormConfig) {
		c.Pattern = pattern
	}
}

// WithStep sets the step value for number inputs
func WithStep(step string) FormOption {
	return func(c *FormConfig) {
		c.Step = step
	}
}

// WithRows sets the number of rows for textarea
func WithRows(rows int) FormOption {
	return func(c *FormConfig) {
		c.Rows = rows
	}
}

// WithCols sets the number of columns for textarea
func WithCols(cols int) FormOption {
	return func(c *FormConfig) {
		c.Cols = cols
	}
}

// WithOptions sets the options for select/radio components
func WithOptions(options []SelectOption) FormOption {
	return func(c *FormConfig) {
		c.Options = options
	}
}

// WithMultiple allows multiple selection
func WithMultiple() FormOption {
	return func(c *FormConfig) {
		c.Multiple = true
	}
}

// WithEmptyOption adds a placeholder option to select
func WithEmptyOption(text string) FormOption {
	return func(c *FormConfig) {
		c.EmptyOption = text
	}
}

// WithError sets the error expression or text
func WithError(expr string) FormOption {
	return func(c *FormConfig) {
		c.Error = expr
	}
}

// WithClass adds additional CSS classes
func WithClass(class string) FormOption {
	return func(c *FormConfig) {
		c.Class = class
	}
}

// WithAttrs adds extra HTML attributes
func WithAttrs(attrs templ.Attributes) FormOption {
	return func(c *FormConfig) {
		c.ExtraAttrs = attrs
	}
}

// Datastar-specific options

// WithModel sets the data-model attribute for two-way binding
func WithModel(expr string) FormOption {
	return func(c *FormConfig) {
		c.Datastar.Model = expr
	}
}

// WithOn adds a data-on:event handler
func WithOn(event, action string) FormOption {
	return func(c *FormConfig) {
		if c.Datastar.On == nil {
			c.Datastar.On = make(map[string]string)
		}
		c.Datastar.On[event] = action
	}
}

// WithBind adds a data-bind:attr binding
func WithBind(attr, expr string) FormOption {
	return func(c *FormConfig) {
		if c.Datastar.Bind == nil {
			c.Datastar.Bind = make(map[string]string)
		}
		c.Datastar.Bind[attr] = expr
	}
}

// WithSignal adds a signal to data-signals
func WithSignal(name, value string) FormOption {
	return func(c *FormConfig) {
		if c.Datastar.Signals == nil {
			c.Datastar.Signals = make(map[string]string)
		}
		c.Datastar.Signals[name] = value
	}
}

// WithDataAttr adds a data-attr:name binding
func WithDataAttr(name, expr string) FormOption {
	return func(c *FormConfig) {
		if c.Datastar.Attrs == nil {
			c.Datastar.Attrs = make(map[string]string)
		}
		c.Datastar.Attrs[name] = expr
	}
}

// WithShow sets the data-show expression
func WithShow(expr string) FormOption {
	return func(c *FormConfig) {
		c.Datastar.Show = expr
	}
}

// ApplyOptions applies all options to a FormConfig and returns it
func ApplyFormOptions(options []FormOption) *FormConfig {
	config := &FormConfig{
		Type: "text", // default type
		Rows: 3,      // default textarea rows
	}
	for _, opt := range options {
		opt(config)
	}
	return config
}

// InputAttrs returns all attributes for an input element including datastar attrs
func (c *FormConfig) InputAttrs() templ.Attributes {
	attrs := templ.Attributes{}

	// Add datastar attributes
	if c.Datastar.Model != "" {
		attrs["data-model"] = c.Datastar.Model
	}
	for event, action := range c.Datastar.On {
		attrs["data-on:"+event] = action
	}
	for attr, expr := range c.Datastar.Bind {
		attrs["data-bind:"+attr] = expr
	}
	for name, expr := range c.Datastar.Attrs {
		attrs["data-attr:"+name] = expr
	}

	// Add extra attrs
	for k, v := range c.ExtraAttrs {
		attrs[k] = v
	}

	return attrs
}

// ButtonConfig holds configuration for button components
type ButtonConfig struct {
	Variant  Variant
	Size     Size
	Type     string // button, submit, reset
	Disabled bool
	Loading  string // data-show expression for loading state

	// Datastar attributes
	Datastar DatastarAttrs

	// Extra attributes
	Class      string
	ExtraAttrs templ.Attributes
}

// ButtonOption is a function that modifies ButtonConfig
type ButtonOption func(*ButtonConfig)

// WithVariant sets the button variant
func WithVariant(v Variant) ButtonOption {
	return func(c *ButtonConfig) {
		c.Variant = v
	}
}

// WithSize sets the component size
func WithSize(s Size) ButtonOption {
	return func(c *ButtonConfig) {
		c.Size = s
	}
}

// WithButtonType sets the button type (button, submit, reset)
func WithButtonType(t string) ButtonOption {
	return func(c *ButtonConfig) {
		c.Type = t
	}
}

// WithButtonDisabled marks the button as disabled
func WithButtonDisabled() ButtonOption {
	return func(c *ButtonConfig) {
		c.Disabled = true
	}
}

// WithLoading sets the loading state expression
func WithLoading(expr string) ButtonOption {
	return func(c *ButtonConfig) {
		c.Loading = expr
	}
}

// WithButtonClass adds additional CSS classes to button
func WithButtonClass(class string) ButtonOption {
	return func(c *ButtonConfig) {
		c.Class = class
	}
}

// WithButtonAttrs adds extra HTML attributes to button
func WithButtonAttrs(attrs templ.Attributes) ButtonOption {
	return func(c *ButtonConfig) {
		c.ExtraAttrs = attrs
	}
}

// WithButtonOn adds a data-on:event handler to button
func WithButtonOn(event, action string) ButtonOption {
	return func(c *ButtonConfig) {
		if c.Datastar.On == nil {
			c.Datastar.On = make(map[string]string)
		}
		c.Datastar.On[event] = action
	}
}

// WithButtonBind adds a data-bind:attr binding to button
func WithButtonBind(attr, expr string) ButtonOption {
	return func(c *ButtonConfig) {
		if c.Datastar.Bind == nil {
			c.Datastar.Bind = make(map[string]string)
		}
		c.Datastar.Bind[attr] = expr
	}
}

// ApplyButtonOptions applies all options to a ButtonConfig and returns it
func ApplyButtonOptions(options []ButtonOption) *ButtonConfig {
	config := &ButtonConfig{
		Variant: VariantPrimary,
		Type:    "button",
	}
	for _, opt := range options {
		opt(config)
	}
	return config
}

// ButtonAttrs returns all attributes for a button element including datastar attrs
func (c *ButtonConfig) ButtonAttrs() templ.Attributes {
	attrs := templ.Attributes{}

	// Add datastar attributes
	for event, action := range c.Datastar.On {
		attrs["data-on:"+event] = action
	}
	for attr, expr := range c.Datastar.Bind {
		attrs["data-bind:"+attr] = expr
	}

	// Add extra attrs
	for k, v := range c.ExtraAttrs {
		attrs[k] = v
	}

	return attrs
}

// AutocompleteConfig holds configuration for autocomplete component
type AutocompleteConfig struct {
	SearchEndpoint string // SSE endpoint for search
	DisplayField   string // field to show in input (default: "name")
	ValueField     string // field to use as value (default: "id")
	MinChars       int    // minimum chars before search (default: 2)
	Debounce       int    // debounce time in ms (default: 300)

	// Inherit from FormConfig
	FormConfig
}

// AutocompleteOption is a function that modifies AutocompleteConfig
type AutocompleteOption func(*AutocompleteConfig)

// WithSearchEndpoint sets the SSE search endpoint
func WithSearchEndpoint(endpoint string) AutocompleteOption {
	return func(c *AutocompleteConfig) {
		c.SearchEndpoint = endpoint
	}
}

// WithDisplayField sets which field to display in the input
func WithDisplayField(field string) AutocompleteOption {
	return func(c *AutocompleteConfig) {
		c.DisplayField = field
	}
}

// WithValueField sets which field to use as the value
func WithValueField(field string) AutocompleteOption {
	return func(c *AutocompleteConfig) {
		c.ValueField = field
	}
}

// WithMinChars sets minimum characters before search triggers
func WithMinChars(n int) AutocompleteOption {
	return func(c *AutocompleteConfig) {
		c.MinChars = n
	}
}

// WithDebounce sets the debounce time in milliseconds
func WithDebounce(ms int) AutocompleteOption {
	return func(c *AutocompleteConfig) {
		c.Debounce = ms
	}
}

// WithAutocompletePlaceholder sets the placeholder for autocomplete
func WithAutocompletePlaceholder(p string) AutocompleteOption {
	return func(c *AutocompleteConfig) {
		c.FormConfig.Placeholder = p
	}
}

// WithAutocompleteModel sets the data-model for the hidden value field
func WithAutocompleteModel(expr string) AutocompleteOption {
	return func(c *AutocompleteConfig) {
		c.FormConfig.Datastar.Model = expr
	}
}

// WithAutocompleteRequired marks the autocomplete as required
func WithAutocompleteRequired() AutocompleteOption {
	return func(c *AutocompleteConfig) {
		c.FormConfig.Required = true
	}
}

// WithAutocompleteDisabled marks the autocomplete as disabled
func WithAutocompleteDisabled() AutocompleteOption {
	return func(c *AutocompleteConfig) {
		c.FormConfig.Disabled = true
	}
}

// WithAutocompleteError sets the error expression for autocomplete
func WithAutocompleteError(expr string) AutocompleteOption {
	return func(c *AutocompleteConfig) {
		c.FormConfig.Error = expr
	}
}

// ApplyAutocompleteOptions applies all options and returns config
func ApplyAutocompleteOptions(options []AutocompleteOption) *AutocompleteConfig {
	config := &AutocompleteConfig{
		DisplayField: "name",
		ValueField:   "id",
		MinChars:     2,
		Debounce:     300,
	}
	for _, opt := range options {
		opt(config)
	}
	return config
}

// DataTableConfig holds configuration for data table component
type DataTableConfig struct {
	DataEndpoint string   // SSE endpoint for data
	Columns      []Column // column definitions
	PageSize     int      // items per page (default: 25)
	Selectable   bool     // show checkbox column
	Sortable     bool     // enable sorting (default: true based on column config)

	// Datastar signals prefix (for namespacing, e.g., "users" -> $users.page)
	SignalPrefix string

	// Custom row actions renderer
	RowActions func(rowIndex int) templ.Component
}

// DataTableOption is a function that modifies DataTableConfig
type DataTableOption func(*DataTableConfig)

// WithDataEndpoint sets the SSE data endpoint
func WithDataEndpoint(endpoint string) DataTableOption {
	return func(c *DataTableConfig) {
		c.DataEndpoint = endpoint
	}
}

// WithColumns sets the table columns
func WithColumns(columns []Column) DataTableOption {
	return func(c *DataTableConfig) {
		c.Columns = columns
	}
}

// WithPageSize sets items per page
func WithPageSize(size int) DataTableOption {
	return func(c *DataTableConfig) {
		c.PageSize = size
	}
}

// WithSelectable enables row selection checkboxes
func WithSelectable() DataTableOption {
	return func(c *DataTableConfig) {
		c.Selectable = true
	}
}

// WithSignalPrefix sets a prefix for datastar signals
func WithSignalPrefix(prefix string) DataTableOption {
	return func(c *DataTableConfig) {
		c.SignalPrefix = prefix
	}
}

// WithRowActions sets a custom row actions renderer
func WithRowActions(fn func(rowIndex int) templ.Component) DataTableOption {
	return func(c *DataTableConfig) {
		c.RowActions = fn
	}
}

// ApplyDataTableOptions applies all options and returns config
func ApplyDataTableOptions(options []DataTableOption) *DataTableConfig {
	config := &DataTableConfig{
		PageSize: 25,
		Sortable: true,
	}
	for _, opt := range options {
		opt(config)
	}
	return config
}

// StackConfig holds configuration for stack layout component
type StackConfig struct {
	Gap       string // gap size: xs, sm, md, lg, xl or custom value
	Direction string // flex direction: col (default), row
	Class     string // additional classes
}

// StackOption is a function that modifies StackConfig
type StackOption func(*StackConfig)

// WithGap sets the gap between stack items
func WithGap(gap string) StackOption {
	return func(c *StackConfig) {
		c.Gap = gap
	}
}

// WithDirection sets the stack direction
func WithDirection(dir string) StackOption {
	return func(c *StackConfig) {
		c.Direction = dir
	}
}

// WithStackClass adds additional CSS classes
func WithStackClass(class string) StackOption {
	return func(c *StackConfig) {
		c.Class = class
	}
}

// ApplyStackOptions applies all options and returns config
func ApplyStackOptions(options []StackOption) *StackConfig {
	config := &StackConfig{
		Gap:       "md",
		Direction: "col",
	}
	for _, opt := range options {
		opt(config)
	}
	return config
}

// FormGroupConfig holds configuration for form group wrapper
type FormGroupConfig struct {
	OnSubmit string // data-on:submit action
	Class    string // additional classes

	// Datastar attributes
	Datastar DatastarAttrs
}

// FormGroupOption is a function that modifies FormGroupConfig
type FormGroupOption func(*FormGroupConfig)

// WithFormSubmit sets the form submit handler
func WithFormSubmit(action string) FormGroupOption {
	return func(c *FormGroupConfig) {
		c.OnSubmit = action
	}
}

// WithFormClass adds additional CSS classes to form
func WithFormClass(class string) FormGroupOption {
	return func(c *FormGroupConfig) {
		c.Class = class
	}
}

// WithFormSignals adds signals to the form
func WithFormSignals(signals map[string]string) FormGroupOption {
	return func(c *FormGroupConfig) {
		c.Datastar.Signals = signals
	}
}

// ApplyFormGroupOptions applies all options and returns config
func ApplyFormGroupOptions(options []FormGroupOption) *FormGroupConfig {
	config := &FormGroupConfig{}
	for _, opt := range options {
		opt(config)
	}
	return config
}
