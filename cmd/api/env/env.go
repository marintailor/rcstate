package env

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/marintailor/rcstate/cmd/api/gce"
	"github.com/marintailor/rcstate/cmd/api/record"
	"github.com/marintailor/rcstate/cmd/api/ssh"
	"gopkg.in/yaml.v2"
)

// Environments stores environment details, and variables.
type Environments struct {
	Envs []Environment `yaml:"environment"`
	Vars Variable      `yaml:"variable"`
}

// Variable stores variables declared in environment file.
type Variable map[string]interface{}

// Environment stores details of an environment.
type Environment struct {
	Group []Group `yaml:"group"`
	Label string  `yaml:"label"`
	Name  string  `yaml:"name"`
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

// Config stores options from parsed flags.
type Config struct {
	Name   string       `json:"name"`
	Label  string       `json:"label"`
	All    bool         `json:"all"`
	Data   Environments `json:"data"`
	File   string
	Host   string
	Dry    bool
	Format string
}

// GetConfig will get configuration from JSON.
func (c *Config) GetConfig(b []byte) error {
	return json.Unmarshal(b, &c)
}

// ParseFlags will parse flags for options.
func (c *Config) ParseFlags(args []string) error {
	f := flag.NewFlagSet(args[0], flag.ExitOnError)

	f.BoolVar(&c.All, "all", false, "all environments")
	f.BoolVar(&c.All, "a", false, "all environments")

	f.BoolVar(&c.Dry, "dry", false, "Run the command without executing it")

	envFile := os.Getenv("RCSTATE_ENV_FILE")
	f.StringVar(&c.File, "env-file", envFile, "environment file")
	f.StringVar(&c.File, "e", envFile, "environment file")

	f.StringVar(&c.Format, "format", "", "Output format of API request")
	f.StringVar(&c.Format, "f", "", "Output format of API request")

	f.StringVar(&c.Host, "host", "", "Server host that will execute the commands")
	f.StringVar(&c.Host, "h", "", "Server host that will execute the commands")

	f.StringVar(&c.Label, "label", "", "environment label")
	f.StringVar(&c.Label, "l", "", "environment label")

	f.StringVar(&c.Name, "name", "", "environment name")
	f.StringVar(&c.Name, "n", "", "environment name")

	f.Usage = func() { fmt.Printf("missing or wrong option(s)\nfor usage information type:\n  rcstate env help\n\n") }

	if err := f.Parse(args[1:]); err != nil {
		return fmt.Errorf("parse flags: %w", err)
	}

	return nil
}

// ParseEnvironmentFile will parse an environment file for environments.
func (c *Config) ParseEnvironmentFile() error {
	if c.File == "" {
		return fmt.Errorf("no environment file was provided")
	}

	if _, err := os.Stat(c.File); err != nil {
		return fmt.Errorf("file stat %q: %w", c.File, err)
	}

	data, err := os.ReadFile(c.File)
	if err != nil {
		return fmt.Errorf("read file %q: %w", data, err)
	}

	if err := yaml.Unmarshal(data, &c.Data); err != nil {
		return fmt.Errorf("unmarshal template: %w", err)
	}

	tpl, err := environmentTemplate(c.Data, c.File)
	if err != nil {
		return fmt.Errorf("environment template: %w", err)
	}

	if err := yaml.Unmarshal(tpl.Bytes(), &c.Data); err != nil {
		return fmt.Errorf("unmarshal template: %w", err)
	}

	return nil
}

// environmentTemplate will fill placeholders inside the environment file.
func environmentTemplate(e Environments, ef string) (bytes.Buffer, error) {
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

// GetData return the environment data.
func (c *Config) GetData() ([]byte, error) {
	return json.Marshal(&c.Data)
}

// NewEnvironments returns an Environment struct.
func NewEnvironments(data string) (*Environments, error) {
	var e Environments

	if err := json.Unmarshal([]byte(data), &e); err != nil {
		return nil, fmt.Errorf("unmarshal env data: %w", err)
	}

	return &e, nil
}

// GetEnvironment return an Environment struct from the list of environments.
func (e *Environments) GetEnvironment(name string, label string) (Environment, error) {
	for _, env := range e.Envs {
		if env.Name == name && env.CheckLabel(label) {
			return env, nil
		}
	}

	return Environment{}, fmt.Errorf("environment %q with label %q not found", name, label)
}

// CheckLabel will check if a specific environment is labeled with provided labels.
func (env *Environment) CheckLabel(label string) bool {
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

// State manages the state of an environment.
func (env *Environment) State(state string) {
	for _, g := range env.Group {
		vm := gce.NewInstances(g.Project, g.Zone)
		for _, instance := range g.Resource.VM.Instance {
			switch state {
			case "up":
				groupStateUp(vm, g, instance)
			case "down":
				groupStateDown(vm, g, instance)
			}
		}
	}
}

// groupStateUp brings a group into Up state.
func groupStateUp(vm *gce.Instances, g Group, instance Instance) {
	if err := vm.Start(instance.Name); err != nil {
		fmt.Println("group state up: vm start:", err)
	}

	if instance.Record.Domain != "" {
		g.instanceRecord(instance)
	}

	host := getHost(instance, g.Project, g.Zone)

	if len(g.Resource.VM.Script.Up) > 0 {
		for _, cmd := range g.Resource.VM.Script.Up {
			cmd = strings.ReplaceAll(cmd, "&gt;", ">")
			g.Resource.VM.Script.execute(host, cmd)
		}
	}

	if len(instance.Script.Up) > 0 {
		for _, cmd := range instance.Script.Up {
			cmd = strings.ReplaceAll(cmd, "&gt;", ">")
			instance.Script.execute(host, cmd)
		}
	}
}

// groupStateUp brings a group into Down state.
func groupStateDown(vm *gce.Instances, g Group, instance Instance) {
	host := getHost(instance, g.Project, g.Zone)

	if len(instance.Script.Down) > 0 {
		for _, cmd := range instance.Script.Down {
			cmd = strings.ReplaceAll(cmd, "&gt;", ">")
			g.Resource.VM.Script.execute(host, cmd)
		}
	}

	if len(g.Resource.VM.Script.Down) > 0 {
		for _, cmd := range g.Resource.VM.Script.Down {
			cmd = strings.ReplaceAll(cmd, "&gt;", ">")
			instance.Script.execute(host, cmd)
		}
	}

	if err := vm.Stop(instance.Name); err != nil {
		fmt.Println("group state down: vm stop:", err)
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
		return
	}

	if err := record.NewRecord(inst.Record.IP, inst.Record.Type, inst.Record.Zone, inst.Record.Domain).Route53(); err != nil {
		fmt.Println("instance record: new record:", err)
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
