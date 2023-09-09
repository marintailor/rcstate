package vm

// Status returns the status of the virtual machine.
func (vm *VirtualMachine) Status(name string) (string, error) {
	return vm.Instances.Status(name)
}
