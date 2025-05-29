package main

import "errors"

type command struct {
	Name string
	Args []string
}

type commands struct {
	registeredCmds map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCmds[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	f, exists := c.registeredCmds[cmd.Name]
	if !exists {
		return errors.New("command does not exist")
	}
	return f(s, cmd)
}
