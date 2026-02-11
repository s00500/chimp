package components

import (
	"strings"

	"github.com/a-h/templ"
)

// DatastarAttrs holds common Datastar attributes that can be applied to components
type DatastarAttrs struct {
	Bind    string            // data-bind (two-way binding for form inputs)
	On      map[string]string // data-on:event -> action
	BindMap map[string]string // data-bind:attr -> expression (for attribute bindings)
	Signals map[string]string // data-signals
	Attrs   map[string]string // data-attr:name -> expression
	Show    string            // data-show expression
	Text    string            // data-text expression
}

// ToAttrs converts DatastarAttrs to templ.Attributes for use in templates
func (d *DatastarAttrs) ToAttrs() templ.Attributes {
	attrs := templ.Attributes{}

	if d.Bind != "" {
		attrs["data-bind"] = d.Bind
	}

	for event, action := range d.On {
		attrs["data-on:"+event] = action
	}

	for attr, expr := range d.BindMap {
		attrs["data-bind:"+attr] = expr
	}

	if len(d.Signals) > 0 {
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

// ============================================================================
// Common Config (embedded by all component configs)
// ============================================================================

// CommonConfig holds attributes shared by most components
type CommonConfig struct {
	ID       string
	Class    string
	Datastar DatastarAttrs
	Attrs    templ.Attributes
}

// CommonAttrs returns all common attributes including datastar attrs
func (c *CommonConfig) CommonAttrs() templ.Attributes {
	attrs := templ.Attributes{}

	if c.ID != "" {
		attrs["id"] = c.ID
	}

	for event, action := range c.Datastar.On {
		attrs["data-on:"+event] = action
	}
	for attr, expr := range c.Datastar.BindMap {
		attrs["data-bind:"+attr] = expr
	}
	for name, expr := range c.Datastar.Attrs {
		attrs["data-attr:"+name] = expr
	}
	if c.Datastar.Show != "" {
		attrs["data-show"] = c.Datastar.Show
	}

	for k, v := range c.Attrs {
		attrs[k] = v
	}

	return attrs
}

// ============================================================================
// Option Interfaces (each component defines what options it accepts)
// ============================================================================

type FormOption interface{ applyToForm(*FormConfig) }
type ButtonOption interface{ applyToButton(*ButtonConfig) }
type AutocompleteOption interface{ applyToAutocomplete(*AutocompleteConfig) }
type DataTableOption interface{ applyToDataTable(*DataTableConfig) }
type StackOption interface{ applyToStack(*StackConfig) }
type FormGroupOption interface{ applyToFormGroup(*FormGroupConfig) }
type CardOption interface{ applyToCard(*CardConfig) }
type SectionOption interface{ applyToSection(*SectionConfig) }

// ============================================================================
// Common Options (implement multiple interfaces - work on many components)
// ============================================================================

// idOption sets the id attribute
type idOption string

func (o idOption) applyToForm(c *FormConfig)               { c.ID = string(o) }
func (o idOption) applyToButton(c *ButtonConfig)           { c.ID = string(o) }
func (o idOption) applyToAutocomplete(c *AutocompleteConfig) { c.ID = string(o) }
func (o idOption) applyToDataTable(c *DataTableConfig)     { c.ID = string(o) }
func (o idOption) applyToStack(c *StackConfig)             { c.ID = string(o) }
func (o idOption) applyToFormGroup(c *FormGroupConfig)     { c.ID = string(o) }
func (o idOption) applyToCard(c *CardConfig)               { c.ID = string(o) }
func (o idOption) applyToSection(c *SectionConfig)         { c.ID = string(o) }

// WithID sets the id attribute (works on any component)
func WithID(id string) idOption { return idOption(id) }

// classOption adds CSS classes
type classOption string

func (o classOption) apply(c *CommonConfig) {
	if c.Class == "" {
		c.Class = string(o)
	} else {
		c.Class += " " + string(o)
	}
}
func (o classOption) applyToForm(c *FormConfig)               { o.apply(&c.CommonConfig) }
func (o classOption) applyToButton(c *ButtonConfig)           { o.apply(&c.CommonConfig) }
func (o classOption) applyToAutocomplete(c *AutocompleteConfig) { o.apply(&c.CommonConfig) }
func (o classOption) applyToDataTable(c *DataTableConfig)     { o.apply(&c.CommonConfig) }
func (o classOption) applyToStack(c *StackConfig)             { o.apply(&c.CommonConfig) }
func (o classOption) applyToFormGroup(c *FormGroupConfig)     { o.apply(&c.CommonConfig) }
func (o classOption) applyToCard(c *CardConfig)               { o.apply(&c.CommonConfig) }
func (o classOption) applyToSection(c *SectionConfig)         { o.apply(&c.CommonConfig) }

// WithClass adds CSS classes (works on any component)
func WithClass(class string) classOption { return classOption(class) }

// onOption adds a data-on:event handler
type onOption struct {
	event  string
	action string
}

func (o onOption) apply(c *CommonConfig) {
	if c.Datastar.On == nil {
		c.Datastar.On = make(map[string]string)
	}
	c.Datastar.On[o.event] = o.action
}
func (o onOption) applyToForm(c *FormConfig)               { o.apply(&c.CommonConfig) }
func (o onOption) applyToButton(c *ButtonConfig)           { o.apply(&c.CommonConfig) }
func (o onOption) applyToAutocomplete(c *AutocompleteConfig) { o.apply(&c.CommonConfig) }
func (o onOption) applyToDataTable(c *DataTableConfig)     { o.apply(&c.CommonConfig) }
func (o onOption) applyToStack(c *StackConfig)             { o.apply(&c.CommonConfig) }
func (o onOption) applyToFormGroup(c *FormGroupConfig)     { o.apply(&c.CommonConfig) }
func (o onOption) applyToCard(c *CardConfig)               { o.apply(&c.CommonConfig) }
func (o onOption) applyToSection(c *SectionConfig)         { o.apply(&c.CommonConfig) }

// WithOn adds a data-on:event handler (works on any component)
func WithOn(event, action string) onOption { return onOption{event, action} }

// bindAttrOption adds a data-bind:attr binding
type bindAttrOption struct {
	attr string
	expr string
}

func (o bindAttrOption) apply(c *CommonConfig) {
	if c.Datastar.BindMap == nil {
		c.Datastar.BindMap = make(map[string]string)
	}
	c.Datastar.BindMap[o.attr] = o.expr
}
func (o bindAttrOption) applyToForm(c *FormConfig)               { o.apply(&c.CommonConfig) }
func (o bindAttrOption) applyToButton(c *ButtonConfig)           { o.apply(&c.CommonConfig) }
func (o bindAttrOption) applyToAutocomplete(c *AutocompleteConfig) { o.apply(&c.CommonConfig) }
func (o bindAttrOption) applyToDataTable(c *DataTableConfig)     { o.apply(&c.CommonConfig) }
func (o bindAttrOption) applyToStack(c *StackConfig)             { o.apply(&c.CommonConfig) }
func (o bindAttrOption) applyToFormGroup(c *FormGroupConfig)     { o.apply(&c.CommonConfig) }
func (o bindAttrOption) applyToCard(c *CardConfig)               { o.apply(&c.CommonConfig) }
func (o bindAttrOption) applyToSection(c *SectionConfig)         { o.apply(&c.CommonConfig) }

// WithBindAttr adds a data-bind:attr binding for attribute bindings (works on any component)
// Example: WithBindAttr("class", "$active ? 'selected' : ''")
func WithBindAttr(attr, expr string) bindAttrOption { return bindAttrOption{attr, expr} }

// showOption sets the data-show expression
type showOption string

func (o showOption) apply(c *CommonConfig) { c.Datastar.Show = string(o) }
func (o showOption) applyToForm(c *FormConfig)               { o.apply(&c.CommonConfig) }
func (o showOption) applyToButton(c *ButtonConfig)           { o.apply(&c.CommonConfig) }
func (o showOption) applyToAutocomplete(c *AutocompleteConfig) { o.apply(&c.CommonConfig) }
func (o showOption) applyToDataTable(c *DataTableConfig)     { o.apply(&c.CommonConfig) }
func (o showOption) applyToStack(c *StackConfig)             { o.apply(&c.CommonConfig) }
func (o showOption) applyToFormGroup(c *FormGroupConfig)     { o.apply(&c.CommonConfig) }
func (o showOption) applyToCard(c *CardConfig)               { o.apply(&c.CommonConfig) }
func (o showOption) applyToSection(c *SectionConfig)         { o.apply(&c.CommonConfig) }

// WithShow sets the data-show expression (works on any component)
func WithShow(expr string) showOption { return showOption(expr) }

// attrsOption adds extra HTML attributes
type attrsOption templ.Attributes

func (o attrsOption) apply(c *CommonConfig) {
	if c.Attrs == nil {
		c.Attrs = make(templ.Attributes)
	}
	for k, v := range o {
		c.Attrs[k] = v
	}
}
func (o attrsOption) applyToForm(c *FormConfig)               { o.apply(&c.CommonConfig) }
func (o attrsOption) applyToButton(c *ButtonConfig)           { o.apply(&c.CommonConfig) }
func (o attrsOption) applyToAutocomplete(c *AutocompleteConfig) { o.apply(&c.CommonConfig) }
func (o attrsOption) applyToDataTable(c *DataTableConfig)     { o.apply(&c.CommonConfig) }
func (o attrsOption) applyToStack(c *StackConfig)             { o.apply(&c.CommonConfig) }
func (o attrsOption) applyToFormGroup(c *FormGroupConfig)     { o.apply(&c.CommonConfig) }
func (o attrsOption) applyToCard(c *CardConfig)               { o.apply(&c.CommonConfig) }
func (o attrsOption) applyToSection(c *SectionConfig)         { o.apply(&c.CommonConfig) }

// WithAttrs adds extra HTML attributes (works on any component)
func WithAttrs(attrs templ.Attributes) attrsOption { return attrsOption(attrs) }

// ============================================================================
// Form Config & Options
// ============================================================================

// FormConfig holds configuration for form input components
type FormConfig struct {
	CommonConfig // embedded

	Type        string // input type: text, email, password, number, etc.
	Placeholder string
	Required    bool
	Disabled    bool
	Readonly    bool
	Autofocus   bool

	Min       string // min value for number/date
	Max       string // max value for number/date
	MinLength int
	MaxLength int
	Pattern   string // regex pattern
	Step      string // step for number inputs

	Rows int // textarea rows
	Cols int // textarea cols

	Options     []SelectOption // select/radio options
	Multiple    bool
	EmptyOption string // placeholder option text

	Error string // data-show expression for error
}

// Form-specific option types
type formTypeOption string
func (o formTypeOption) applyToForm(c *FormConfig) { c.Type = string(o) }

type formPlaceholderOption string
func (o formPlaceholderOption) applyToForm(c *FormConfig) { c.Placeholder = string(o) }

type formRequiredOption struct{}
func (o formRequiredOption) applyToForm(c *FormConfig) { c.Required = true }

type formDisabledOption struct{}
func (o formDisabledOption) applyToForm(c *FormConfig) { c.Disabled = true }

type formReadonlyOption struct{}
func (o formReadonlyOption) applyToForm(c *FormConfig) { c.Readonly = true }

type formAutofocusOption struct{}
func (o formAutofocusOption) applyToForm(c *FormConfig) { c.Autofocus = true }

type formMinOption string
func (o formMinOption) applyToForm(c *FormConfig) { c.Min = string(o) }

type formMaxOption string
func (o formMaxOption) applyToForm(c *FormConfig) { c.Max = string(o) }

type formMinLengthOption int
func (o formMinLengthOption) applyToForm(c *FormConfig) { c.MinLength = int(o) }

type formMaxLengthOption int
func (o formMaxLengthOption) applyToForm(c *FormConfig) { c.MaxLength = int(o) }

type formPatternOption string
func (o formPatternOption) applyToForm(c *FormConfig) { c.Pattern = string(o) }

type formStepOption string
func (o formStepOption) applyToForm(c *FormConfig) { c.Step = string(o) }

type formRowsOption int
func (o formRowsOption) applyToForm(c *FormConfig) { c.Rows = int(o) }

type formColsOption int
func (o formColsOption) applyToForm(c *FormConfig) { c.Cols = int(o) }

type formOptionsOption []SelectOption
func (o formOptionsOption) applyToForm(c *FormConfig) { c.Options = []SelectOption(o) }

type formMultipleOption struct{}
func (o formMultipleOption) applyToForm(c *FormConfig) { c.Multiple = true }

type formEmptyOptionOption string
func (o formEmptyOptionOption) applyToForm(c *FormConfig) { c.EmptyOption = string(o) }

type formErrorOption string
func (o formErrorOption) applyToForm(c *FormConfig) { c.Error = string(o) }

type formBindOption string
func (o formBindOption) applyToForm(c *FormConfig) { c.Datastar.Bind = string(o) }

type formSignalOption struct {
	name  string
	value string
}
func (o formSignalOption) applyToForm(c *FormConfig) {
	if c.Datastar.Signals == nil {
		c.Datastar.Signals = make(map[string]string)
	}
	c.Datastar.Signals[o.name] = o.value
}

type formDataAttrOption struct {
	name string
	expr string
}
func (o formDataAttrOption) applyToForm(c *FormConfig) {
	if c.Datastar.Attrs == nil {
		c.Datastar.Attrs = make(map[string]string)
	}
	c.Datastar.Attrs[o.name] = o.expr
}

// Form option constructors
func WithType(t string) formTypeOption             { return formTypeOption(t) }
func WithPlaceholder(p string) formPlaceholderOption { return formPlaceholderOption(p) }
func WithRequired() formRequiredOption             { return formRequiredOption{} }
func WithDisabled() formDisabledOption             { return formDisabledOption{} }
func WithReadonly() formReadonlyOption             { return formReadonlyOption{} }
func WithAutofocus() formAutofocusOption           { return formAutofocusOption{} }
func WithMin(min string) formMinOption             { return formMinOption(min) }
func WithMax(max string) formMaxOption             { return formMaxOption(max) }
func WithMinLength(n int) formMinLengthOption      { return formMinLengthOption(n) }
func WithMaxLength(n int) formMaxLengthOption      { return formMaxLengthOption(n) }
func WithPattern(pattern string) formPatternOption { return formPatternOption(pattern) }
func WithStep(step string) formStepOption          { return formStepOption(step) }
func WithRows(rows int) formRowsOption             { return formRowsOption(rows) }
func WithCols(cols int) formColsOption             { return formColsOption(cols) }
func WithOptions(options []SelectOption) formOptionsOption { return formOptionsOption(options) }
func WithMultiple() formMultipleOption             { return formMultipleOption{} }
func WithEmptyOption(text string) formEmptyOptionOption { return formEmptyOptionOption(text) }
func WithError(expr string) formErrorOption        { return formErrorOption(expr) }
func WithBind(expr string) formBindOption          { return formBindOption(expr) }
func WithSignal(name, value string) formSignalOption { return formSignalOption{name, value} }
func WithDataAttr(name, expr string) formDataAttrOption { return formDataAttrOption{name, expr} }

// applyFormOptions applies all options to a FormConfig
func applyFormOptions(options []FormOption) *FormConfig {
	config := &FormConfig{
		Type: "text",
		Rows: 3,
	}
	for _, opt := range options {
		opt.applyToForm(config)
	}
	return config
}

// InputAttrs returns all attributes for an input element
func (c *FormConfig) InputAttrs() templ.Attributes {
	attrs := c.CommonAttrs()

	if c.Datastar.Bind != "" {
		attrs["data-bind"] = c.Datastar.Bind
	}

	return attrs
}

// ============================================================================
// Button Config & Options
// ============================================================================

// ButtonConfig holds configuration for button components
type ButtonConfig struct {
	CommonConfig // embedded

	Variant  Variant
	Size     Size
	Type     string // button, submit, reset
	Disabled bool
	Loading  string // data-show expression for loading state
}

// Button-specific option types
type buttonVariantOption Variant
func (o buttonVariantOption) applyToButton(c *ButtonConfig) { c.Variant = Variant(o) }

type buttonSizeOption Size
func (o buttonSizeOption) applyToButton(c *ButtonConfig) { c.Size = Size(o) }

type buttonTypeOption string
func (o buttonTypeOption) applyToButton(c *ButtonConfig) { c.Type = string(o) }

type buttonDisabledOption struct{}
func (o buttonDisabledOption) applyToButton(c *ButtonConfig) { c.Disabled = true }

type buttonLoadingOption string
func (o buttonLoadingOption) applyToButton(c *ButtonConfig) { c.Loading = string(o) }

// Button option constructors
func WithVariant(v Variant) buttonVariantOption     { return buttonVariantOption(v) }
func WithSize(s Size) buttonSizeOption              { return buttonSizeOption(s) }
func WithButtonType(t string) buttonTypeOption      { return buttonTypeOption(t) }
func WithButtonDisabled() buttonDisabledOption      { return buttonDisabledOption{} }
func WithLoading(expr string) buttonLoadingOption   { return buttonLoadingOption(expr) }

// applyButtonOptions applies all options to a ButtonConfig
func applyButtonOptions(options []ButtonOption) *ButtonConfig {
	config := &ButtonConfig{
		Variant: VariantPrimary,
		Type:    "button",
	}
	for _, opt := range options {
		opt.applyToButton(config)
	}
	return config
}

// ButtonAttrs returns all attributes for a button element
func (c *ButtonConfig) ButtonAttrs() templ.Attributes {
	return c.CommonAttrs()
}

// ============================================================================
// Autocomplete Config & Options
// ============================================================================

// AutocompleteConfig holds configuration for autocomplete component
type AutocompleteConfig struct {
	CommonConfig // embedded

	SearchEndpoint string
	DisplayField   string
	ValueField     string
	MinChars       int
	Debounce       int

	Placeholder string
	Required    bool
	Disabled    bool
	Error       string
}

// Autocomplete-specific option types
type acSearchEndpointOption string
func (o acSearchEndpointOption) applyToAutocomplete(c *AutocompleteConfig) { c.SearchEndpoint = string(o) }

type acDisplayFieldOption string
func (o acDisplayFieldOption) applyToAutocomplete(c *AutocompleteConfig) { c.DisplayField = string(o) }

type acValueFieldOption string
func (o acValueFieldOption) applyToAutocomplete(c *AutocompleteConfig) { c.ValueField = string(o) }

type acMinCharsOption int
func (o acMinCharsOption) applyToAutocomplete(c *AutocompleteConfig) { c.MinChars = int(o) }

type acDebounceOption int
func (o acDebounceOption) applyToAutocomplete(c *AutocompleteConfig) { c.Debounce = int(o) }

type acPlaceholderOption string
func (o acPlaceholderOption) applyToAutocomplete(c *AutocompleteConfig) { c.Placeholder = string(o) }

type acBindOption string
func (o acBindOption) applyToAutocomplete(c *AutocompleteConfig) { c.Datastar.Bind = string(o) }

type acRequiredOption struct{}
func (o acRequiredOption) applyToAutocomplete(c *AutocompleteConfig) { c.Required = true }

type acDisabledOption struct{}
func (o acDisabledOption) applyToAutocomplete(c *AutocompleteConfig) { c.Disabled = true }

type acErrorOption string
func (o acErrorOption) applyToAutocomplete(c *AutocompleteConfig) { c.Error = string(o) }

// Autocomplete option constructors
func WithSearchEndpoint(endpoint string) acSearchEndpointOption { return acSearchEndpointOption(endpoint) }
func WithDisplayField(field string) acDisplayFieldOption       { return acDisplayFieldOption(field) }
func WithValueField(field string) acValueFieldOption           { return acValueFieldOption(field) }
func WithMinChars(n int) acMinCharsOption                      { return acMinCharsOption(n) }
func WithDebounce(ms int) acDebounceOption                     { return acDebounceOption(ms) }
func WithAutocompletePlaceholder(p string) acPlaceholderOption { return acPlaceholderOption(p) }
func WithAutocompleteBind(expr string) acBindOption            { return acBindOption(expr) }
func WithAutocompleteRequired() acRequiredOption               { return acRequiredOption{} }
func WithAutocompleteDisabled() acDisabledOption               { return acDisabledOption{} }
func WithAutocompleteError(expr string) acErrorOption          { return acErrorOption(expr) }

// applyAutocompleteOptions applies all options and returns config
func applyAutocompleteOptions(options []AutocompleteOption) *AutocompleteConfig {
	config := &AutocompleteConfig{
		DisplayField: "name",
		ValueField:   "id",
		MinChars:     2,
		Debounce:     300,
	}
	for _, opt := range options {
		opt.applyToAutocomplete(config)
	}
	return config
}

// ============================================================================
// DataTable Config & Options
// ============================================================================

// DataTableConfig holds configuration for data table component
type DataTableConfig struct {
	CommonConfig // embedded

	DataEndpoint string
	Columns      []Column
	PageSize     int
	Selectable   bool
	Sortable     bool
	SignalPrefix string
	RowActions   func(rowIndex int) templ.Component
}

// DataTable-specific option types
type dtDataEndpointOption string
func (o dtDataEndpointOption) applyToDataTable(c *DataTableConfig) { c.DataEndpoint = string(o) }

type dtColumnsOption []Column
func (o dtColumnsOption) applyToDataTable(c *DataTableConfig) { c.Columns = []Column(o) }

type dtPageSizeOption int
func (o dtPageSizeOption) applyToDataTable(c *DataTableConfig) { c.PageSize = int(o) }

type dtSelectableOption struct{}
func (o dtSelectableOption) applyToDataTable(c *DataTableConfig) { c.Selectable = true }

type dtSignalPrefixOption string
func (o dtSignalPrefixOption) applyToDataTable(c *DataTableConfig) { c.SignalPrefix = string(o) }

type dtRowActionsOption func(rowIndex int) templ.Component
func (o dtRowActionsOption) applyToDataTable(c *DataTableConfig) { c.RowActions = o }

// DataTable option constructors
func WithDataEndpoint(endpoint string) dtDataEndpointOption         { return dtDataEndpointOption(endpoint) }
func WithColumns(columns []Column) dtColumnsOption                  { return dtColumnsOption(columns) }
func WithPageSize(size int) dtPageSizeOption                        { return dtPageSizeOption(size) }
func WithSelectable() dtSelectableOption                            { return dtSelectableOption{} }
func WithSignalPrefix(prefix string) dtSignalPrefixOption           { return dtSignalPrefixOption(prefix) }
func WithRowActions(fn func(rowIndex int) templ.Component) dtRowActionsOption { return dtRowActionsOption(fn) }

// applyDataTableOptions applies all options and returns config
func applyDataTableOptions(options []DataTableOption) *DataTableConfig {
	config := &DataTableConfig{
		PageSize: 25,
		Sortable: true,
	}
	for _, opt := range options {
		opt.applyToDataTable(config)
	}
	return config
}

// ============================================================================
// Stack Config & Options
// ============================================================================

// StackConfig holds configuration for stack layout component
type StackConfig struct {
	CommonConfig // embedded

	Gap       string
	Direction string
}

// Stack-specific option types
type stackGapOption string
func (o stackGapOption) applyToStack(c *StackConfig) { c.Gap = string(o) }

type stackDirectionOption string
func (o stackDirectionOption) applyToStack(c *StackConfig) { c.Direction = string(o) }

// Stack option constructors
func WithGap(gap string) stackGapOption           { return stackGapOption(gap) }
func WithDirection(dir string) stackDirectionOption { return stackDirectionOption(dir) }

// applyStackOptions applies all options and returns config
func applyStackOptions(options []StackOption) *StackConfig {
	config := &StackConfig{
		Gap:       "md",
		Direction: "col",
	}
	for _, opt := range options {
		opt.applyToStack(config)
	}
	return config
}

// ============================================================================
// FormGroup Config & Options
// ============================================================================

// FormGroupConfig holds configuration for form group wrapper
type FormGroupConfig struct {
	CommonConfig // embedded

	OnSubmit string
}

// FormGroup-specific option types
type fgSubmitOption string
func (o fgSubmitOption) applyToFormGroup(c *FormGroupConfig) { c.OnSubmit = string(o) }

type fgSignalsOption map[string]string
func (o fgSignalsOption) applyToFormGroup(c *FormGroupConfig) { c.Datastar.Signals = o }

// FormGroup option constructors
func WithFormSubmit(action string) fgSubmitOption              { return fgSubmitOption(action) }
func WithFormSignals(signals map[string]string) fgSignalsOption { return fgSignalsOption(signals) }

// applyFormGroupOptions applies all options and returns config
func applyFormGroupOptions(options []FormGroupOption) *FormGroupConfig {
	config := &FormGroupConfig{}
	for _, opt := range options {
		opt.applyToFormGroup(config)
	}
	return config
}

// ============================================================================
// Card Config & Options
// ============================================================================

// CardConfig holds configuration for card components
type CardConfig struct {
	CommonConfig // embedded

	Title   string
	Padding string
}

// Card-specific option types
type cardTitleOption string
func (o cardTitleOption) applyToCard(c *CardConfig) { c.Title = string(o) }

type cardPaddingOption string
func (o cardPaddingOption) applyToCard(c *CardConfig) { c.Padding = string(o) }

// Card option constructors
func WithCardTitle(title string) cardTitleOption { return cardTitleOption(title) }
func WithPadding(size string) cardPaddingOption  { return cardPaddingOption(size) }

// applyCardOptions applies all options and returns config
func applyCardOptions(options []CardOption) *CardConfig {
	config := &CardConfig{
		Padding: "md",
	}
	for _, opt := range options {
		opt.applyToCard(config)
	}
	return config
}

// ============================================================================
// Section Config & Options
// ============================================================================

// SectionConfig holds configuration for section components
type SectionConfig struct {
	CommonConfig // embedded

	Title   string
	Padding string
}

// Section-specific option types
type sectionTitleOption string
func (o sectionTitleOption) applyToSection(c *SectionConfig) { c.Title = string(o) }

type sectionPaddingOption string
func (o sectionPaddingOption) applyToSection(c *SectionConfig) { c.Padding = string(o) }

// Section option constructors
func WithSectionTitle(title string) sectionTitleOption   { return sectionTitleOption(title) }
func WithSectionPadding(size string) sectionPaddingOption { return sectionPaddingOption(size) }

// applySectionOptions applies all options and returns config
func applySectionOptions(options []SectionOption) *SectionConfig {
	config := &SectionConfig{
		Padding: "md",
	}
	for _, opt := range options {
		opt.applyToSection(config)
	}
	return config
}

// ============================================================================
// Element Config & Options (for Div, Span, etc.)
// ============================================================================

// ElementOption is implemented by options that can be applied to generic elements
type ElementOption interface{ applyToElement(*ElementConfig) }

// ElementConfig holds configuration for generic element components (Div, Span, etc.)
type ElementConfig struct {
	CommonConfig // embedded
}

// Common options for Element
func (o idOption) applyToElement(c *ElementConfig)    { c.ID = string(o) }
func (o classOption) applyToElement(c *ElementConfig) { o.apply(&c.CommonConfig) }
func (o onOption) applyToElement(c *ElementConfig)    { o.apply(&c.CommonConfig) }
func (o bindAttrOption) applyToElement(c *ElementConfig)  { o.apply(&c.CommonConfig) }
func (o showOption) applyToElement(c *ElementConfig)  { o.apply(&c.CommonConfig) }
func (o attrsOption) applyToElement(c *ElementConfig) { o.apply(&c.CommonConfig) }

// Element-specific options
type elementBindOption string

func (o elementBindOption) applyToElement(c *ElementConfig) { c.Datastar.Bind = string(o) }

type elementSignalsOption map[string]string

func (o elementSignalsOption) applyToElement(c *ElementConfig) { c.Datastar.Signals = o }

type elementTextOption string

func (o elementTextOption) applyToElement(c *ElementConfig) { c.Datastar.Text = string(o) }

// Element option constructors
func WithElementBind(expr string) elementBindOption               { return elementBindOption(expr) }
func WithSignals(signals map[string]string) elementSignalsOption  { return elementSignalsOption(signals) }
func WithText(expr string) elementTextOption                      { return elementTextOption(expr) }

// applyElementOptions applies all options and returns config
func applyElementOptions(options []ElementOption) *ElementConfig {
	config := &ElementConfig{}
	for _, opt := range options {
		opt.applyToElement(config)
	}
	return config
}
