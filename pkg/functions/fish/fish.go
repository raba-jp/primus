package fish

type VariableScope int

const (
	UniversalScope VariableScope = iota + 1
	GlobalScope
	LocalScope
)
