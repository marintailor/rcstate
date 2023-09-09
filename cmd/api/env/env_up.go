package env

import "fmt"

// Up will turn up one or all environments based on provided config..
func (c *Config) Up() (string, error) {
	if c.Name != "" {
		return c.UpSingle()
	}

	return c.UpAll()
}

// UpSingle will turn up an environment.
func (c *Config) UpSingle() (string, error) {
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

// UpAll will turn up all environments.
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
