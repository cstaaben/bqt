// Package formatter provides functionality for formatting data into strings.
package formatter

var SupportedFormats = map[string]struct{}{
	"table": {},
	"csv":   {},
	"json":  {},
}

// Formatter is an interface for formatters that can format data into a string.
type Formatter interface {
	// Format takes a map of data and returns a string.
	Format(data map[string]any) (string, error)
}

// Table formats data into a table.
type Table struct{}

// NewTable creates a new Table formatter.
func NewTable() *Table {
	return &Table{}
}

func (t *Table) Format(data map[string]any) (string, error) {
	panic("not implemented")
}

// CSV formats data into a CSV.
type CSV struct{}

// NewCSV creates a new CSV formatter.
func NewCSV() *CSV {
	return &CSV{}
}

func (t *CSV) Format(data map[string]any) (string, error) {
	panic("not implemented")
}

// JSON formats data into a JSON.
type JSON struct{}

// NewJSON creates a new JSON formatter.
func NewJSON() *JSON {
	return &JSON{}
}

func (t *JSON) Format(data map[string]any) (string, error) {
	panic("not implemented")
}
