package linting

// Failure is the linting error model
type Failure struct {
	RuleName string
	Message  string
	Position int // line number in the changelog file
}
