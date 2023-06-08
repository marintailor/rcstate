package cmd

import (
	"fmt"
	"strings"
)

// up will start all resources in environments.
func (e *Environments) up() int {
	switch {
	case e.name != "":
		env, err := e.getEnvironment(e.name)
		if err != nil {
			fmt.Println(err)
			return 1
		}

		if !env.checkLabel(e.label) {
			fmt.Printf("environment %q is not labeled with %q\n", e.name, e.label)
			return 1
		}

		fmt.Printf("\n%s\nENVIRONMENT: %s\nLABEL: %s\n%s\n", strings.Repeat("=", 40), env.Name, env.Label, strings.Repeat("=", 40))
		groupState(env.Group, "up")

		return 0

	case e.all:
		var count int

		for _, env := range e.Envs {
			if !env.checkLabel(e.label) {
				continue
			}

			fmt.Printf("\n%s\nENVIRONMENT: %s\nLABEL: %s\n%s\n", strings.Repeat("=", 40), env.Name, env.Label, strings.Repeat("=", 40))
			groupState(env.Group, "up")
			count++
		}

		if len(e.Envs) > 0 && count == 0 {
			fmt.Printf("no environment is labeled with %q\n", e.label)
		}
	}

	return 0
}
