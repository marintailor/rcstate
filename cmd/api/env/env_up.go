package env

import "fmt"

func (c *Config) Up() (string, error) {
	if c.Name != "" {
		return c.UpName()
	}

	return c.UpAll()
}

func (c *Config) UpName() (string, error) {
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

	env.State("up")

	return "{ \"status\": \"success\" }", nil
}

func (c *Config) UpAll() (string, error) {
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
			env.State("up")
		}
	}

	return "{ \"status\": \"success\" }", nil
}
