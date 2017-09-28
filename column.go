package xql

// buildColumn
// Build a Column object according to given field and tag, a tag should be:
// `xql:"Column('name',TYPE, ...)"`
// Available parameters:
// - Sequence()
// - ForeinKey()
// - nullable=True/False
// - unique=True/False
// - primary_key=True/False
// - index=True/False
// - default=TYPE|VALUE|FUNCTION
// -
func buildColumn(prop interface{}, tag string) Column {
	return Column{}
}
