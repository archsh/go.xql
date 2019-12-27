package xql

type ConditionType uint

const (
	CONDITION_AND ConditionType = iota
	CONDITION_OR
)

type OrderType uint

const (
	ORDER_ASC OrderType = iota
	ORDER_DESC
)

type QueryFilter struct {
	Condition ConditionType // AND , OR
	Reversed  bool          // Reversed Column and Value if it is true
	Field     string
	Operator  string // Value will not used if empty.
	Function  string
	Value     interface{}
}

type QueryOrder struct {
	Type  OrderType
	Field string
}

type QueryColumn struct {
	FieldName string
	Function  string
	Alias     string
}

type UpdateColumn struct {
	Field    string
	Operator string
	Value    interface{}
}

type QueryExtra map[string]interface{}
