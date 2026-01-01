package health

type Report struct {
	Validator string   `json:"validator"`
	Healthy   bool     `json:"healthy"`
	Issues    []string `json:"issues"`

	Jailed     bool `json:"jailed"`
	Bonded     bool `json:"bonded"`
	HasTokens  bool `json:"has_tokens"`
	CatchingUp bool `json:"catching_up"`
}
