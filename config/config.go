package config

type (
	Config struct {
		Scenarios map[string]*Scenario
		Options   *Options `json:"options,omitempty"`
	}

	Options struct {
		ForceSequential bool
	}
)
