package config

type Scenario struct {
	Name         string
	Dependencies []string
	Tests        []*Test
}
