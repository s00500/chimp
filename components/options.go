package components

import "strings"

// Size represents component size variants
type Size string

const (
	SizeXS Size = "xs"
	SizeSM Size = "sm"
	SizeMD Size = "md"
	SizeLG Size = "lg"
	SizeXL Size = "xl"
)

// Variant represents button/component style variants
type Variant string

const (
	VariantPrimary     Variant = "primary"
	VariantSecondary   Variant = "secondary"
	VariantOutline     Variant = "outline"
	VariantGhost       Variant = "ghost"
	VariantDestructive Variant = "destructive"
)

// SelectOption represents an option in a select, radio group, or autocomplete
type SelectOption struct {
	Value    string `json:"value"`
	Label    string `json:"label"`
	Disabled bool   `json:"disabled,omitempty"`
}

// Column represents a column definition for DataTable
type Column struct {
	Field       string
	Header      string
	Sortable    bool
	Type        ColumnType
	Width       string   // optional width (e.g., "100px", "20%")
	EditType    EditType // how this column is edited (empty = read-only in edit mode)
	EditOptions []string // options for EditSelect
}

// ColumnType defines special column types
type ColumnType string

const (
	ColumnText    ColumnType = ""
	ColumnActions ColumnType = "actions"
	ColumnStatus  ColumnType = "status"
	ColumnDate    ColumnType = "date"
)

// EditType defines how a column can be edited in edit mode
type EditType string

const (
	EditNone   EditType = ""       // not editable (default)
	EditText   EditType = "text"   // <input type="text">
	EditSelect EditType = "select" // <select> with options
)

// EditSignalName returns the Datastar signal name for an editable column field.
// e.g. "name" → "editName", "email" → "editEmail"
func EditSignalName(field string) string {
	if field == "" {
		return ""
	}
	return "edit" + strings.ToUpper(field[:1]) + field[1:]
}

// editSignalsInit builds a data-signals value string for initializing edit signals.
func editSignalsInit(values map[string]string, columns []Column) string {
	var sb strings.Builder
	sb.WriteString("{")
	first := true
	for _, col := range columns {
		if col.EditType == "" || col.EditType == EditNone {
			continue
		}
		if !first {
			sb.WriteString(", ")
		}
		val := values[col.Field]
		val = strings.ReplaceAll(val, `\`, `\\`)
		val = strings.ReplaceAll(val, `'`, `\'`)
		sb.WriteString(EditSignalName(col.Field))
		sb.WriteString(": '")
		sb.WriteString(val)
		sb.WriteString("'")
		first = false
	}
	sb.WriteString("}")
	return sb.String()
}
