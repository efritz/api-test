package jsonconfig

import "github.com/efritz/api-test/config"

type BasicAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *BasicAuth) Translate() (*config.BasicAuth, error) {
	if c == nil {
		return nil, nil
	}

	return &config.BasicAuth{
		Username: c.Username, // TODO - should compile
		Password: c.Password,
	}, nil
}
