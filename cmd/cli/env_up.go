package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	client "github.com/marintailor/rcstate/client/env"
	"github.com/marintailor/rcstate/cmd/api/env"
)

// envUp will stop all resources in environments.
func envUp(args []string) int {
	cfg := env.Config{}

	if err := cfg.ParseFlags(args); err != nil {
		fmt.Println("get config:", err)
	}

	if cfg.Host != "" {
		return upRemote(&cfg)
	}

	return upLocal(&cfg)
}

// upLocal will stop resources by executing the logic locally.
func upLocal(c *env.Config) int {
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
		env, err := e.GetEnvironment(c.Name, c.Label)
		if err != nil {
			fmt.Println("up local: get env:", err)
			return 1
		}

		fmt.Printf("\n%s\nENVIRONMENT: %s\nLABEL: %s\n%s\n", strings.Repeat("=", 40), env.Name, env.Label, strings.Repeat("=", 40))
		env.State("up")

		return 0
	case c.All:
		var count int

		for _, env := range e.Envs {
			if !env.CheckLabel(c.Label) {
				continue
			}

			fmt.Printf("\n%s\nENVIRONMENT: %s\nLABEL: %s\n%s\n", strings.Repeat("=", 40), env.Name, env.Label, strings.Repeat("=", 40))
			env.State("up")
			count++
		}

		if len(e.Envs) > 0 && count == 0 {
			fmt.Printf("no environment is labeled with %q\n", c.Label)
		}
	}

	return 0
}

// upRemote will stop resources by sending a request to remote server.
func upRemote(c *env.Config) int {
	if err := c.ParseEnvironmentFile(); err != nil {
		fmt.Println("parse env file:", err)
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

	var buff bytes.Buffer

	enc := json.NewEncoder(&buff)
	enc.SetEscapeHTML(false)

	err = enc.Encode(c)
	if err != nil {
		log.Fatal(err)
	}

	if c.Format == "json" {
		fmt.Println(buff.String())
	}

	if !c.Dry {
		_, err = client.Up(buff.String(), c.Host)
		if err != nil {
			fmt.Println("client env up:", err)
			return 1
		}
	}

	return 0
}
