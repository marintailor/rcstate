package env

import (
	"encoding/json"
	"fmt"

	"github.com/marintailor/rcstate/cmd/api/gce"
)

type ShowEnvironment struct {
	Name  string      `json:"name"`
	Label string      `json:"label"`
	Group []ShowGroup `json:"group"`
}

type ShowGroup struct {
	Name     string       `json:"name"`
	Project  string       `json:"project"`
	Zone     string       `json:"zone"`
	Resource ShowResource `json:"resource"`
}

type ShowResource struct {
	VM []gce.Instance `json:"vm"`
}

func (se *ShowEnvironment) GetDetails(e Environment) {
	se.Name = e.Name
	se.GetDetailsGroup(e.Group)
}

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

func GetDetailsVM(insts []Instance, project string, zone string) []gce.Instance {
	var list []gce.Instance

	instances := *gce.NewInstances(project, zone)
	if err := instances.GetInstancesList(); err != nil {
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

func (c *Config) Show() (string, error) {
	// log.Println(c)
	if c.Name != "" {
		return c.ShowName()
	}

	return c.ShowAll()
}

func (c *Config) ShowName() (string, error) {
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
	out.GetDetails(env)

	json, err := json.Marshal(out)
	if err != nil {
		return "", fmt.Errorf("marshal env json: %w", err)
	}

	return string(json), nil
}

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
			item.GetDetails(env)
			list = append(list, item)
		}
	}

	json, err := json.Marshal(list)
	if err != nil {
		return "", fmt.Errorf("marshal env list json: %w", err)
	}

	return string(json), nil
}
