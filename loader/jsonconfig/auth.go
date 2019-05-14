package jsonconfig

import (
	"fmt"

	"github.com/efritz/api-test/config"
)

type BasicAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *BasicAuth) Translate() (*config.BasicAuth, error) {
	if a == nil {
		return nil, nil
	}

	usernameTemplate, err := compile(a.Username)
	if err != nil {
		return nil, fmt.Errorf("illegal username template (%s)", err.Error())
	}

	passwordTemplate, err := compile(a.Password)
	if err != nil {
		return nil, fmt.Errorf("illegal password template (%s)", err.Error())
	}

	return &config.BasicAuth{
		Username: usernameTemplate,
		Password: passwordTemplate,
	}, nil
}
