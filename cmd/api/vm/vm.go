package vm

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	"github.com/marintailor/rcstate/cmd/api/gce"
	"github.com/marintailor/rcstate/cmd/api/record"
	"github.com/marintailor/rcstate/cmd/api/ssh"
)

// VirtualMachine holds configuration and methods to manage virtual machine instances.
type VirtualMachine struct {
	Instances gce.Instances
	Project   string
	Zone      string
}

// Config stores options from parsed flags.
type Config struct {
	DNS        DNS
	Dry        bool
	ExternalIP bool
	Format     string
	Host       string
	Ip         string
	IpList     []string
	Name       string
	Project    string
	Script     VMScript
	Zone       string
}

// DNS stores DNS configuration.
type DNS struct {
	Domain     string
	RecordName string
	RecordType string
}

// VMScript stores shell commands.
type VMScript struct {
	CMD string
	SSH ssh.SSH
}

// NewVirtualMachine returns a VirtualMachine struct.
func NewVirtualMachine(project string, zone string) (*VirtualMachine, error) {
	var vm VirtualMachine

	vm.Instances = *gce.NewInstances(project, zone)

	return &vm, nil
}

// GetConfig will get configuration from JSON.
func (c *Config) GetConfig(b []byte) error {
	return json.Unmarshal(b, &c)
}

// ParseFlags will parse flags for options.
func (c *Config) ParseFlags(args []string) error {
	f := flag.NewFlagSet(args[0], flag.ExitOnError)

	f.StringVar(&c.DNS.Domain, "domain", "", "Domain for DNS record")
	f.StringVar(&c.DNS.Domain, "d", "", "Domain for DNS record")

	f.BoolVar(&c.Dry, "dry", false, "Run the command without executing it")

	f.BoolVar(&c.ExternalIP, "external-ip", false, "Get external IP address for DNS record")

	f.StringVar(&c.Format, "format", "", "Output format of API request")
	f.StringVar(&c.Format, "f", "", "Output format of API request")

	f.StringVar(&c.Host, "host", "", "Server host that will execute the commands")
	f.StringVar(&c.Host, "h", "", "Server host that will execute the commands")

	f.StringVar(&c.Ip, "ip", "", "IP addresses for DNS record")

	f.StringVar(&c.Name, "name", "", "Virtual Machine instance name")
	f.StringVar(&c.Name, "n", "", "Virtual Machine instance name")

	f.StringVar(&c.Project, "project", "", "Google Cloud Project ID")
	f.StringVar(&c.Project, "p", "", "Google Cloud Project ID")

	f.StringVar(&c.DNS.RecordName, "dns-record-name", "", "DNS record name")

	f.StringVar(&c.DNS.RecordType, "dns-record-type", "", "Create the DNS record")

	f.StringVar(&c.Script.CMD, "script", "", "run shell command on remote host")
	f.StringVar(&c.Script.CMD, "s", "", "run shell command on remote host")

	f.StringVar(&c.Script.SSH.Key, "ssh-key", "", "path to the SSH private key")

	f.StringVar(&c.Script.SSH.Port, "ssh-port", "", "SSH port number")

	f.StringVar(&c.Script.SSH.User, "ssh-user", "", "SSH username")

	f.StringVar(&c.Zone, "zone", "", "Google Cloud Zone name")
	f.StringVar(&c.Zone, "z", "", "Google Cloud Zone name")

	f.Usage = func() { fmt.Printf("missing or wrong option(s)\nfor usage information type:\n  rcstate vm help\n\n") }

	if err := f.Parse(args[1:]); err != nil {
		return fmt.Errorf("parse flags for command %q: %w", args[0], err)
	}

	return nil
}

// Record will create a DNS record.
func (c *Config) Record(dnsRecord string) {
	if c.ExternalIP {
		externalIP, err := gce.GetInstanceExternalIP(c.Name, c.Project, c.Zone)
		if err != nil {
			fmt.Printf("create record: get external IP address: %s", err)
			return
		}

		if externalIP == "" {
			fmt.Println("create record: instance does not have external IP address")
		}

		c.IpList = append(c.IpList, externalIP)
	}

	if c.Ip != "" {
		ips := strings.Split(c.Ip, ",")
		c.IpList = append(c.IpList, ips...)
	}

	if record.CheckRecordIP(dnsRecord, c.IpList) {
		return
	}

	if err := record.NewRecord(c.IpList, c.DNS.RecordType, dnsRecord, c.DNS.Domain).Route53(); err != nil {
		fmt.Println("new record:", err)
		return
	}
}

// ExecuteScript will execute the shell commands.
func (c *Config) ExecuteScript() {
	host := c.getHost()

	script, err := ssh.NewSSH(host, c.Script.SSH.Port, c.Script.SSH.User, c.Script.SSH.Key)
	if err != nil {
		fmt.Println("vm new ssh:", err)
		return
	}

	if err := script.CMD(c.Script.CMD); err != nil {
		fmt.Println("vm script cmd:", err)
	}
}

// getHost return a valid host address.
func (c *Config) getHost() string {
	var host string
	if c.DNS.RecordName != "" && c.DNS.Domain != "" {
		host = fmt.Sprintf("%s.%s", c.DNS.RecordName, c.DNS.Domain)
	}

	if host == "" && c.ExternalIP {
		var err error
		host, err = gce.GetInstanceExternalIP(c.Name, c.Project, c.Zone)
		if err != nil {
			fmt.Printf("create record: get external IP address: %s", err)
		}
	}

	if host == "" && len(c.Ip) > 0 {
		ipList := strings.Split(c.Ip, ",")
		if len(ipList) > 0 {
			host = ipList[0]
		}
	}

	return host
}
