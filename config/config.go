package config

type (
	Config struct {
		BaseURL              string
		GlobalRequestHeaders map[string]string
		Tests                []*Test
	}
)
