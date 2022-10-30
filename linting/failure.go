package linting

type Failure struct {
	RuleName string
	Message  string
	Position int
}
