package cmd

import (
	"flag"
	"fmt"
	"strings"

	"github.com/marintailor/rcstate/api/gce"
	"github.com/marintailor/rcstate/api/record"
)

// VirtualMachine holds configuration and methods to manage virtual machine instances.
type VirtualMachine struct {
	Instances gce.Instances
	Opts      options
}

// options stores options from parsed flags.
type options struct {
	dns        dns
	externalIP bool
	ip         string
	ipList     []string
	name       string
	project    string
	zone       string
}

// dns stores DNS configuration.
type dns struct {
	domain     string
	recordName string
	recordType string
}

// vmRun executes the command 'vm'.
func vmRun(args []string) int {
	if len(args) < 2 {
		vmHelp()
		return 0
	}

	vm, err := NewVirtualMachine(args)
	if err != nil {
		fmt.Println("new virtual machine:", err)
		vmHelp()
		return 1
	}

	cmds := map[string]func() int{
		"list":   func() int { return vm.list() },
		"start":  func() int { return vm.start() },
		"status": func() int { return vm.status() },
		"stop":   func() int { return vm.stop() },
	}

	cmd, ok := cmds[args[0]]
	if !ok {
		fmt.Println("no such command: vm", args[0])
		help()
		return 1
	}

	return cmd()
}

// NewVirtualMachine returns a VirtualMachine struct.
func NewVirtualMachine(args []string) (*VirtualMachine, error) {
	var vm VirtualMachine

	if err := vm.getOptions(args); err != nil {
		return &vm, fmt.Errorf("parse flags: %w", err)
	}

	vm.Instances = *gce.NewInstances(vm.Opts.project, vm.Opts.zone)

	return &vm, nil
}

// getOptions will parse flags for options.
func (v *VirtualMachine) getOptions(args []string) error {
	f := flag.NewFlagSet(args[0], flag.ExitOnError)

	f.StringVar(&v.Opts.dns.domain, "domain", "", "Domain for DNS record")
	f.StringVar(&v.Opts.dns.domain, "d", "", "Domain for DNS record")

	f.BoolVar(&v.Opts.externalIP, "external-ip", false, "Get external IP address for DNS record")

	f.StringVar(&v.Opts.ip, "ip", "", "IP addresses for DNS record")

	f.StringVar(&v.Opts.name, "name", "", "Virtual Machine instance name")
	f.StringVar(&v.Opts.name, "n", "", "Virtual Machine instance name")

	f.StringVar(&v.Opts.project, "project", "", "Google Cloud Project ID")
	f.StringVar(&v.Opts.project, "p", "", "Google Cloud Project ID")

	f.StringVar(&v.Opts.dns.recordName, "dns-record-name", "", "DNS record name")

	f.StringVar(&v.Opts.dns.recordType, "dns-record-type", "", "Create the DNS record")

	f.StringVar(&v.Opts.zone, "zone", "", "Google Cloud Zone name")
	f.StringVar(&v.Opts.zone, "z", "", "Google Cloud Zone name")

	f.Usage = func() { vmHelp() }

	if err := f.Parse(args[1:]); err != nil {
		return fmt.Errorf("parse flags for command %q: %w", args[0], err)
	}

	return nil
}

// record will create a DNS record.
func (v *VirtualMachine) record(dnsRecord string) {
	if v.Opts.externalIP {
		externalIP, err := gce.GetInstanceExternalIP(v.Opts.name, v.Opts.project, v.Opts.zone)
		if err != nil {
			fmt.Printf("create record: get external IP address: %s", err)
			return
		}

		if externalIP == "" {
			fmt.Println("create record: instance does not have external IP address")
		}

		v.Opts.ipList = append(v.Opts.ipList, externalIP)
	}

	if v.Opts.ip != "" {
		ips := strings.Split(v.Opts.ip, ",")
		v.Opts.ipList = append(v.Opts.ipList, ips...)
	}

	if record.CheckRecordIP(dnsRecord, v.Opts.ipList) {
		fmt.Println("    DNS record is up-to-date")
		return
	}

	if err := record.NewRecord(v.Opts.ipList, v.Opts.dns.recordType, dnsRecord).Route53(); err != nil {
		fmt.Println(err)
		return
	}
}
