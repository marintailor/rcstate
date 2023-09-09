package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	client "github.com/marintailor/rcstate/client/env"
	"github.com/marintailor/rcstate/cmd/api/env"
	"github.com/marintailor/rcstate/cmd/api/gce"
)

// envShow returns the information about the environment(s).
func envShow(args []string) int {
	cfg := env.Config{}

	if err := cfg.ParseFlags(args); err != nil {
		fmt.Println("get config:", err)
	}

	if cfg.Host != "" {
		return showRemote(&cfg)
	}

	return showLocal(&cfg)
}

// showLocal returns the information about the environment(s) by executing the logic locally.
func showLocal(c *env.Config) int {
	if err := c.ParseEnvironmentFile(); err != nil {
		fmt.Println("parse env file:", err)
		return 1
	}

	data, err := c.GetData()
	if err != nil {
		fmt.Println("get config data:", err)
	}

	e, err := env.NewEnvironments(string(data))
	if err != nil {
		fmt.Println("new environment:", err)
		return 1
	}

	switch {
	case c.Name != "":
		environment, err := e.GetEnvironment(c.Name, c.Label)
		if err != nil {
			fmt.Println("show local: get env:", err)
			return 1
		}

		showEnvironment(environment)

		return 0
	case c.All:
		var count int

		for _, environment := range e.Envs {
			if !environment.CheckLabel(c.Label) {
				continue
			}

			showEnvironment(environment)

			count++
		}

		if len(e.Envs) > 0 && count == 0 {
			fmt.Printf("no environment is labeled with %q\n", c.Label)
		}
	}

	return 0
}

// showRemote returns the information about the environment(s) by sending a request to remote server.
func showRemote(c *env.Config) int {
	if err := c.ParseEnvironmentFile(); err != nil {
		fmt.Println("parse env file:", err)
		return 1
	}

	j, err := json.Marshal(c)
	if err != nil {
		fmt.Println("marshal config:", err)
		return 1
	}

	if c.Format == "json" {
		fmt.Println(string(j))
	}

	if !c.Dry {
		data, err := client.Show(string(j), c.Host)
		if err != nil {
			fmt.Println("client env show:", err)
			return 1
		}

		label, err := remoteCheckLabel(c)
		if err != nil {
			fmt.Println("remote check label:", err)
		}

		if !label {
			fmt.Println("no environment was found with label", c.Label)
			return 1
		}

		switch {
		case c.Name != "":
			var se env.ShowEnvironment

			if err := json.Unmarshal([]byte(data), &se); err != nil {
				fmt.Println("unmarshal env show:", err)
				return 1
			}

			showEnvironmentRemote(se)

			return 0
		case c.All:
			var se []env.ShowEnvironment

			if err := json.Unmarshal([]byte(data), &se); err != nil {
				fmt.Println("unmarshal env show:", err)
				return 1
			}

			for _, e := range se {
				showEnvironmentRemote(e)
			}
		}
	}

	return 0
}

// CheckLabel will check if a specific environment is labeled with provided labels.
func remoteCheckLabel(c *env.Config) (bool, error) {
	data, err := c.GetData()
	if err != nil {
		fmt.Println("get config data:", err)
	}

	e, err := env.NewEnvironments(string(data))
	if err != nil {
		fmt.Println("new environment:", err)
		return false, err
	}

	var count int

	for _, environment := range e.Envs {
		if !environment.CheckLabel(c.Label) {
			continue
		}

		count++
	}

	if len(e.Envs) > 0 && count == 0 {
		return false, nil
	}

	return true, nil
}

// showEnvironment will show all resources in specific environment.
func showEnvironment(env env.Environment) {
	fmt.Printf("\n%s\nENVIRONMENT: %s\nLABEL: %s\n%s\n", strings.Repeat("=", 40), env.Name, env.Label, strings.Repeat("=", 40))
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

// showEnvironment will show all resources in specific environment.
func showEnvironmentRemote(env env.ShowEnvironment) {
	fmt.Printf("\n%s\nENVIRONMENT: %s\nLABEL: %s\n%s\n", strings.Repeat("=", 40), env.Name, env.Label, strings.Repeat("=", 40))
	for i, g := range env.Group {
		if i > 0 {
			fmt.Println(strings.Repeat("-", 40))
		}

		fmt.Printf("\nGROUP: %s\nPROJECT: %s\nZONE: %s\n\n", g.Name, g.Project, g.Zone)
		fmt.Printf("VIRTUAL MACHINES\n")

		pw := padWidthRemote(g.Resource.VM)

		list, err := gce.NewInstances(g.Project, g.Zone).GetList()
		if err != nil {
			fmt.Printf("new instances: %v", err)
		}

		for j, instance := range g.Resource.VM {
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

// padWidth returns width of the pad.
func padWidthRemote(instances []gce.Instance) int {
	pw := 16
	for _, instance := range instances {
		if len(instance.Name) >= pw {
			pw = len(instance.Name) + 2
		}
	}

	return pw
}
