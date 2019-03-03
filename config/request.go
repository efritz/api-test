package config

type (
	Request struct {
		URI     string
		Method  string
		Headers map[string]string
		Auth    *BasicAuth
		Body    string
	}

	BasicAuth struct {
		Username string
		Password string
	}
)
