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

// envDown will stop all resources in environments.
func envDown(args []string) int {
	cfg := env.Config{}

	if err := cfg.ParseFlags(args); err != nil {
		fmt.Println("get config:", err)
	}

	if err := cfg.ParseEnvironmentFile(); err != nil {
		fmt.Println("parse env file:", err)
		return 1
	}

	if cfg.Host != "" {
		return downRemote(&cfg)
	}

	return downLocal(&cfg)
}

// downLocal will stop resources by executing the logic locally.
func downLocal(cfg *env.Config) int {
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
			fmt.Println("down local: get env:", err)
			return 1
		}

		fmt.Printf("\n%s\nENVIRONMENT: %s\nLABEL: %s\n%s\n", strings.Repeat("=", 40), env.Name, env.Label, strings.Repeat("=", 40))
		env.State("down")

		return 0
	case cfg.All:
		var count int

		for _, env := range e.Envs {
			if !env.CheckLabel(cfg.Label) {
				continue
			}

			fmt.Printf("\n%s\nENVIRONMENT: %s\nLABEL: %s\n%s\n", strings.Repeat("=", 40), env.Name, env.Label, strings.Repeat("=", 40))
			env.State("down")
			count++
		}

		if len(e.Envs) > 0 && count == 0 {
			fmt.Printf("no environment is labeled with %q\n", cfg.Label)
		}
	}

	return 0
}

// downRemote will stop resources by sending a request to remote server.
func downRemote(c *env.Config) int {
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
		_, err = client.Down(buff.String(), c.Host)
		if err != nil {
			fmt.Println("client env down:", err)
			return 1
		}
	}

	return 0
}
