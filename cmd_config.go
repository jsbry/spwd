package main

var cmdConfig = &Command{
	Run:       runConfig,
	UsageLine: "config",
	Short:     "Create config.yml",
	Long: `Create a configuration file.
If config.yml exists, use it.
`,
}

func runConfig(ctx context, args []string) error {

	return nil
}
