package env

import "fmt"

func (c *Config) Down() (string, error) {
	if c.Name != "" {
		return c.DownName()
	}

	return c.DownAll()
}

func (c *Config) DownName() (string, error) {
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

	env.State("down")

	return "{ \"status\": \"success\" }", nil
}

func (c *Config) DownAll() (string, error) {
	data, err := c.GetData()
	if err != nil {
		return "", fmt.Errorf("marshal env: %w", err)
	}

	e, err := NewEnvironments(string(data))
	if err != nil {
		return "", fmt.Errorf("new environment: %w", err)
	}

	for _, env := range e.Envs {
		if env.CheckLabel(c.Label) {
			env.State("down")
		}
	}

	return "{ \"status\": \"success\" }", nil
}
