package components

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
	Value    string
	Label    string
	Disabled bool
}

// Column represents a column definition for DataTable
type Column struct {
	Field    string
	Header   string
	Sortable bool
	Type     ColumnType
	Width    string // optional width (e.g., "100px", "20%")
}

// ColumnType defines special column types
type ColumnType string

const (
	ColumnText    ColumnType = ""
	ColumnActions ColumnType = "actions"
	ColumnStatus  ColumnType = "status"
	ColumnDate    ColumnType = "date"
)
