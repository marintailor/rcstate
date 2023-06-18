package cli

import (
	"fmt"
	"strings"

	"github.com/marintailor/rcstate/cmd/api/env"
)

// up will stop all resources in environments.
func envUp(args []string) int {
	cfg := env.Config{}

	if err := cfg.ParseFlags(args); err != nil {
		fmt.Println("get config:", err)
	}

	return upLocal(&cfg)
}

func upLocal(cfg *env.Config) int {
	if err := cfg.ParseEnvironmentFile(); err != nil {
		fmt.Println("parse env file:", err)
		return 1
	}

	data, err := cfg.GetData()
	if err != nil {
		fmt.Println("get config data:", err)
	}

	e, err := env.NewEnvironments(string(data))
	if err != nil {
		fmt.Println("new environment:", err)
		return 1
	}

	switch {
	case cfg.Name != "":
		env, err := e.GetEnvironment(cfg.Name, cfg.Label)
		if err != nil {
			fmt.Println(err)
			return 1
		}

		fmt.Printf("\n%s\nENVIRONMENT: %s\nLABEL: %s\n%s\n", strings.Repeat("=", 40), env.Name, env.Label, strings.Repeat("=", 40))
		env.State("up")

		return 0
	case cfg.All:
		var count int

		for _, env := range e.Envs {
			if !env.CheckLabel(cfg.Label) {
				continue
			}

			fmt.Printf("\n%s\nENVIRONMENT: %s\nLABEL: %s\n%s\n", strings.Repeat("=", 40), env.Name, env.Label, strings.Repeat("=", 40))
			env.State("up")
			count++
		}

		if len(e.Envs) > 0 && count == 0 {
			fmt.Printf("no environment is labeled with %q\n", cfg.Label)
		}
	}

	return 0
}
