package jsonconfig

import "github.com/efritz/api-test/config"

type BasicAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *BasicAuth) Translate() (*config.BasicAuth, error) {
	if a == nil {
		return nil, nil
	}

	return &config.BasicAuth{
		Username: a.Username, // TODO - should compile
		Password: a.Password,
	}, nil
}
