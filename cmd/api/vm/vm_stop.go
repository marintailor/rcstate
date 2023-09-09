package vm

// Stop will stop a virtual machine.
func (vm *VirtualMachine) Stop(name string) error {
	return vm.Instances.Stop(name)
}
