package vm

// Start will stop a virtual machine.
func (vm *VirtualMachine) Start(name string) error {
	return vm.Instances.Start(name)
}
