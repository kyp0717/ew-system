package controllers

type SearchBarArgs struct {
	FromMenu  string
	SKUs      []string
	Names     []string
	Category  []string // Added Category field
	KeyValue  []float64
	KeyString []string
}
