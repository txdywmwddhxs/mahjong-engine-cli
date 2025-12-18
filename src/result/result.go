package result

type Result int

const (
	Default Result = iota
	Win
	Lose
	Equal
	Cancel
)
