package env

import (
	"encoding/json"
	"fmt"

	"github.com/marintailor/rcstate/cmd/api/gce"
)

// ShowEnvironment stores environment information for show command.
type ShowEnvironment struct {
	Name  string      `json:"name"`
	Label string      `json:"label"`
	Group []ShowGroup `json:"group"`
}

// ShowEnvironment stores group information for show command.
type ShowGroup struct {
	Name     string       `json:"name"`
	Project  string       `json:"project"`
	Zone     string       `json:"zone"`
	Resource ShowResource `json:"resource"`
}

// ShowEnvironment stores resource information for show command.
type ShowResource struct {
	VM []gce.Instance `json:"vm"`
}

// GetDetailsEnv will get details about the environment for show command.
func (se *ShowEnvironment) GetDetailsEnv(e Environment) {
	se.Name = e.Name
	se.GetDetailsGroup(e.Group)
}

// GetDetailsEnv will get details about the group for show command.
func (se *ShowEnvironment) GetDetailsGroup(groups []Group) {
	for _, g := range groups {
		group := ShowGroup{}
		group.Name = g.Name
		group.Project = g.Project
		group.Zone = g.Zone
		group.Resource.VM = GetDetailsVM(g.Resource.VM.Instance, g.Project, g.Zone)

		se.Group = append(se.Group, group)
	}

}

// GetDetailsEnv will get details about virtual machines for show command.
func GetDetailsVM(insts []Instance, project string, zone string) []gce.Instance {
	var list []gce.Instance

	instances := *gce.NewInstances(project, zone)
	if _, err := instances.GetInstancesList(); err != nil {
		fmt.Println("list instances:", err)
	}

	for _, instance := range instances.List {
		for _, inst := range insts {
			if instance.Name == inst.Name {
				list = append(list, instance)
			}
		}
	}

	return list
}

// Show returns the information about the environment(s).
func (c *Config) Show() (string, error) {
	if c.Name != "" {
		return c.ShowSingle()
	}

	return c.ShowAll()
}

// ShowSingle returns the information about an environment.
func (c *Config) ShowSingle() (string, error) {
	data, err := c.GetData()
	if err != nil {
		return "", fmt.Errorf("marshal env: %w", err)
	}

	e, err := NewEnvironments(string(data))
	if err != nil {
		return "", fmt.Errorf("new environment: %w", err)
	}

	env, err := e.GetEnvironment(c.Name, c.Label)
	if err != nil {
		return "", fmt.Errorf("get environment: %w", err)
	}

	if env.Name == "" {
		return fmt.Sprintf("{ \"error\": \"environment %q with label %q not found\"}", c.Name, c.Label), nil
	}

	var out ShowEnvironment
	out.GetDetailsEnv(env)

	json, err := json.Marshal(out)
	if err != nil {
		return "", fmt.Errorf("marshal env json: %w", err)
	}

	return string(json), nil
}

// ShowAll returns the information about all environments.
func (c *Config) ShowAll() (string, error) {
	data, err := c.GetData()
	if err != nil {
		return "", fmt.Errorf("marshal env: %w", err)
	}

	e, err := NewEnvironments(string(data))
	if err != nil {
		return "", fmt.Errorf("new environment: %w", err)
	}

	var list []ShowEnvironment
	for _, env := range e.Envs {
		var item ShowEnvironment
		if env.CheckLabel(c.Label) {
			item.GetDetailsEnv(env)
			list = append(list, item)
		}
	}

	json, err := json.Marshal(list)
	if err != nil {
		return "", fmt.Errorf("marshal env list json: %w", err)
	}

	return string(json), nil
}
