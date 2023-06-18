package cli

import (
	"fmt"
	"strings"

	"github.com/marintailor/rcstate/cmd/api/env"
	"github.com/marintailor/rcstate/cmd/api/gce"
)

func envShow(args []string) int {
	cfg := env.Config{}

	if err := cfg.ParseFlags(args); err != nil {
		fmt.Println("get config:", err)
	}

	return showLocal(&cfg)
}

func showLocal(cfg *env.Config) int {
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
		environment, err := e.GetEnvironment(cfg.Name, cfg.Label)
		if err != nil {
			fmt.Println(err)
			return 1
		}

		showEnvironment(environment)

		return 0
	case cfg.All:
		var count int

		for _, environment := range e.Envs {
			if !environment.CheckLabel(cfg.Label) {
				continue
			}

			showEnvironment(environment)

			count++
		}

		if len(e.Envs) > 0 && count == 0 {
			fmt.Printf("no environment is labeled with %q\n", cfg.Label)
		}
	}

	return 0
}

// showEnvironment will show all resources in specific environment.
func showEnvironment(env env.Environment) {
	fmt.Printf("\n%s\nENVIRONMENT LOCAL: %s\nLABEL: %s\n%s\n", strings.Repeat("=", 40), env.Name, env.Label, strings.Repeat("=", 40))
	for i, g := range env.Group {
		if i > 0 {
			fmt.Println(strings.Repeat("-", 40))
		}

		fmt.Printf("\nGROUP: %s\nPROJECT: %s\nZONE: %s\n\n", g.Name, g.Project, g.Zone)
		fmt.Printf("VIRTUAL MACHINES\n")

		pw := padWidth(g.Resource.VM.Instance)

		list, err := gce.NewInstances(g.Project, g.Zone).GetList()
		if err != nil {
			fmt.Printf("new instances: %v", err)
		}

		for j, instance := range g.Resource.VM.Instance {
			var status string
			for _, inst := range list {
				if inst.Name == instance.Name {
					status = inst.Status
				}
			}
			padding := strings.Repeat(" ", pw-len(instance.Name))
			fmt.Printf("%d. %s%sStatus: %s\n", j+1, instance.Name, padding, status)
		}
		fmt.Println()
	}
}

// padWidth returns width of the pad.
func padWidth(instances []env.Instance) int {
	pw := 16
	for _, instance := range instances {
		if len(instance.Name) >= pw {
			pw = len(instance.Name) + 2
		}
	}

	return pw
}
