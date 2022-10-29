package linting

type Config struct {
	Rules     []Rule
	RuleConfs map[string][]any
}
