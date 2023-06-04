package cmd

// vmRun executes the command 'vm'.
func vmRun(args []string) int {
	if len(args) == 0 {
		help()
		return 0
	}

	return vmHelp()
}
