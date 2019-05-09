package config

type Test struct {
	Name     string
	Enabled  bool
	Disabled bool
	Request  *Request
	Response *Response
}
