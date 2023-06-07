package cmd

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/marintailor/rcstate/api/gce"
	"github.com/marintailor/rcstate/api/record"
	"github.com/marintailor/rcstate/api/ssh"

	"gopkg.in/yaml.v3"
)

// Environments stores all environments, provided Options, and related information.
type Environments struct {
	Envs     []Environment `yaml:"environment"`
	Provider string        `yaml:"provider"`
	Vars     Variable      `yaml:"variable"`
	all      bool
	file     string
	label    string
	name     string
}

// Variable stores variables declared in environment file.
type Variable map[string]interface{}

// Environment stores details of an environment.
type Environment struct {
	Group []Group `yaml:"group"`
	Label string  `yaml:"label"`
	Name  string  `yaml:"name"`
	State string
}

// Group stores details of a group.
type Group struct {
	Name     string   `yaml:"name"`
	Project  string   `yaml:"project"`
	Resource Resource `yaml:"resource"`
	Zone     string   `yaml:"zone"`
}

// Resource stores declared resources in a group.
type Resource struct {
	VM VM `yaml:"vm"`
}

// VM stores details about Virtual Machine resource.
type VM struct {
	Instance []Instance `yaml:"instance"`
	Script   EnvScript  `yaml:"script"`
}

// Instance stores details of an instance in Virtual Machine resource.
type Instance struct {
	Name   string    `yaml:"name"`
	Record Record    `yaml:"record"`
	Script EnvScript `yaml:"script"`
}

// EnvScript stores shell commands.
type EnvScript struct {
	Down []string `yaml:"down"`
	SSH  SSH      `yaml:"ssh"`
	Up   []string `yaml:"up"`
}

// SSH stores configuration for SSH connection.
type SSH struct {
	Key  string `yaml:"key"`
	Port string `yaml:"port"`
	User string `yaml:"user"`
}

// Record stores details to create a DNS record.
type Record struct {
	Domain     string   `yaml:"domain"`
	ExternalIP bool     `yaml:"external_ip"`
	IP         []string `yaml:"ip"`
	Type       string   `yaml:"type"`
	Zone       string   `yaml:"zone"`
}

// envRun will run logic to manage environments.
func envRun(args []string) int {
	envs, err := NewEnvironments(args)
	if err != nil {
		fmt.Println("new environments:", err)
		envHelp()
		return 1
	}

	if envs.name == "" && !envs.all {
		fmt.Println("environment not specified")
		return 1
	}

	commands := map[string]func() int{
		"down": func() int { return envs.down() },
		"help": func() int { return envHelp() },
		"show": func() int { return envs.show() },
		"up":   func() int { return envs.up() },
	}

	command, ok := commands[args[0]]
	if !ok {
		fmt.Println("No such command: env", args[0])
		help()
		return 1
	}

	return command()
}

// NewEnvironments return a Environments struct.
func NewEnvironments(args []string) (*Environments, error) {
	var e Environments

	if len(args) < 2 && args[0] != "help" {
		return nil, fmt.Errorf("not enough options")
	}

	if err := e.getFlags(args); err != nil {
		return nil, fmt.Errorf("get flags: %w", err)
	}

	if err := e.parseEnvironmentFile(); err != nil {
		return nil, fmt.Errorf("get flags: %w", err)
	}

	return &e, nil
}

// getFlags will parse flags for environment options.
func (e *Environments) getFlags(args []string) error {
	f := flag.NewFlagSet(args[0], flag.ExitOnError)

	f.BoolVar(&e.all, "all", false, "all environments")
	f.BoolVar(&e.all, "a", false, "all environments")

	f.StringVar(&e.file, "env-file", "", "environment file")
	f.StringVar(&e.file, "e", "", "environment file")

	f.StringVar(&e.label, "label", "", "environment label")
	f.StringVar(&e.label, "l", "", "environment label")

	f.StringVar(&e.name, "name", "", "environment name")
	f.StringVar(&e.name, "n", "", "environment name")

	if err := f.Parse(args[1:]); err != nil {
		return fmt.Errorf("parse flags: %w", err)
	}

	return nil
}

// parseEnvironmentFile will parse an environment file for environments.
func (e *Environments) parseEnvironmentFile() error {
	if e.file == "" {
		return fmt.Errorf("no environment file was provided")
	}

	if _, err := os.Stat(e.file); err != nil {
		return fmt.Errorf("file stat %q: %w", e.file, err)
	}

	data, err := os.ReadFile(e.file)
	if err != nil {
		return fmt.Errorf("read file %q: %w", data, err)
	}

	if err := yaml.Unmarshal(data, &e); err != nil {
		return fmt.Errorf("unmarshal template: %w", err)
	}

	tpl, err := e.environmentTemplate(e.file)
	if err != nil {
		return fmt.Errorf("environment template: %w", err)
	}

	if err := yaml.Unmarshal(tpl.Bytes(), &e); err != nil {
		return fmt.Errorf("unmarshal template: %w", err)
	}

	return nil
}

// environmentTemplate will fill placeholders inside the environment file.
func (e *Environments) environmentTemplate(ef string) (bytes.Buffer, error) {
	var d bytes.Buffer

	if _, err := os.Stat(ef); err != nil {
		return d, fmt.Errorf("file stat %q: %w", ef, err)
	}

	f, err := os.ReadFile(ef)
	if err != nil {
		return d, fmt.Errorf("read file %q: %w", ef, err)
	}

	tpl, err := template.New("").Parse(string(f))
	if err != nil {
		return d, fmt.Errorf("payload template: %s", err)
	}

	if err = tpl.Execute(&d, e.Vars); err != nil {
		return d, fmt.Errorf("write template: %s", err)
	}

	return d, nil
}

// getEnvironment will return a specific environment.
func (e *Environments) getEnvironment(name string) (Environment, error) {
	for _, env := range e.Envs {
		if env.Name == name {
			return env, nil
		}
	}

	return Environment{}, fmt.Errorf("environment %q present", name)
}

// checkLabel will check if a specific environment is labeled with provided labels.
func (env *Environment) checkLabel(label string) bool {
	if label == "" {
		return true
	}

	if env.Label == "" && label != "" {
		return false
	}

	labelList := strings.Split(label, ",")
	labelNumber := len(labelList)

	var count int
	for _, l := range labelList {
		envLabel := strings.Split(env.Label, ",")
		for _, el := range envLabel {
			if el == l {
				count++
			}
		}
	}

	return labelNumber == count
}

// groupState will bring a group to a desired state.
func groupState(group []Group, state string) {
	for i, g := range group {
		vm := gce.NewInstances(g.Project, g.Zone)

		if i > 0 {
			fmt.Println(strings.Repeat("-", 40))
		}

		fmt.Printf("\nGROUP: %s\nPROJECT: %s\nZONE: %s\n\n", g.Name, g.Project, g.Zone)

		for _, instance := range g.Resource.VM.Instance {
			switch state {
			case "up":
				groupStateUp(vm, g, instance)
			case "down":
				groupStateDown(vm, g, instance)
			}
			fmt.Printf("=== Done\n\n")
		}
	}
}

// groupStateUp brings a group into Up state.
func groupStateUp(vm *gce.Instances, g Group, instance Instance) {
	fmt.Printf("=== Start instance %q\n", instance.Name)
	if err := vm.Start(instance.Name); err != nil {
		fmt.Println(err)
	}

	if instance.Record.Domain != "" {
		fmt.Println("=== Create DNS record:", instance.Record.Zone)
		g.instanceRecord(instance)
	}

	host := getHost(instance, g.Project, g.Zone)

	if len(g.Resource.VM.Script.Up) > 0 {
		fmt.Println("=== Execute VM script")
		for _, cmd := range g.Resource.VM.Script.Up {
			g.Resource.VM.Script.execute(host, cmd)
		}
	}

	if len(instance.Script.Up) > 0 {
		fmt.Println("=== Execute Instance script")
		for _, cmd := range instance.Script.Up {
			instance.Script.execute(host, cmd)
		}
	}
}

// groupStateUp brings a group into Down state.
func groupStateDown(vm *gce.Instances, g Group, instance Instance) {
	fmt.Printf("=== Stop instance %q\n", instance.Name)

	host := getHost(instance, g.Project, g.Zone)

	if len(instance.Script.Down) > 0 {
		fmt.Println("=== Execute Instance script")
		for _, cmd := range instance.Script.Down {
			g.Resource.VM.Script.execute(host, cmd)
		}
	}

	if len(g.Resource.VM.Script.Down) > 0 {
		fmt.Println("=== Execute VM script")
		for _, cmd := range g.Resource.VM.Script.Down {
			instance.Script.execute(host, cmd)
		}
	}

	if err := vm.Stop(instance.Name); err != nil {
		fmt.Println(err)
	}
}

// instanceRecord creates the DNS record for an instance.
func (g *Group) instanceRecord(inst Instance) {
	if inst.Record.ExternalIP {
		externalIP, err := gce.GetInstanceExternalIP(inst.Name, g.Project, g.Zone)
		if err != nil {
			fmt.Printf("create record: get external IP address: %s", err)
			return
		}

		if externalIP == "<nil>" {
			fmt.Println("create record: instance does not have external IP address")
		}

		if externalIP != "<nil>" {
			inst.Record.IP = append(inst.Record.IP, externalIP)
		}
	}

	if record.CheckRecordIP(inst.Record.Zone, inst.Record.IP) {
		fmt.Println("    DNS record is up-to-date")
		return
	}

	if err := record.NewRecord(inst.Record.IP, inst.Record.Type, inst.Record.Zone, inst.Record.Domain).Route53(); err != nil {
		fmt.Println(err)
		return
	}
}

// getHost return a valid host address.
func getHost(inst Instance, p string, z string) string {
	if inst.Record.Zone != "" {
		return inst.Record.Zone
	}

	if len(inst.Record.IP) != 0 {
		return inst.Record.IP[0]
	}

	ip, err := gce.GetInstanceExternalIP(inst.Name, p, z)
	if err != nil {
		fmt.Println("host external IP:", err)
	}

	return ip
}

// execute will execute a shell command.
func (s *EnvScript) execute(host string, cmd string) {
	script, err := ssh.NewSSH(host, s.SSH.Port, s.SSH.User, s.SSH.Key)
	if err != nil {
		fmt.Println("env script:", err)
		return
	}

	if err := script.CMD(cmd); err != nil {
		fmt.Println("env script cmd:", err)
	}
}
