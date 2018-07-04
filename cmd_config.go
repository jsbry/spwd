package main

import (
	"fmt"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

var cmdConfig = &Command{
	Run:       runConfig,
	UsageLine: "config",
	Short:     "Create config.yml",
	Long: `Create a configuration file.
If config.yml exists, use it.
`,
}

func runConfig(ctx context, args []string) error {
	keyFile, dataFile, filteringCommand, unprotectiveCommands, err := configScan()

	initCfg := Config{
		KeyFile:              keyFile,
		DataFile:             dataFile,
		FilteringCommand:     filteringCommand,
		UnprotectiveCommands: unprotectiveCommands,
	}
	err = initCfg.configSave()
	if err != nil {
		return err
	}

	return nil
}

// Save config to file on path.
func (cfg Config) configSave() error {
	path, err := app.ConfigFile("config.yml")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	p, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	f.WriteString(string(p))
	return nil
}

func configScan() (keyFile string, dataFile string, filteringCommand string, unprotectiveCommands []string, err error) {
	if keyFile, err = scanText("keyFile: "); err != nil {
		return
	}
	if keyFile == "" {
		keyFile = "~/.ssh/id_rsa"
	}
	if dataFile, err = scanText("dataFile: "); err != nil {
		return
	}
	if dataFile == "" {
		dataFile = "~/.local/share/spwd/data.dat"
	}
	if filteringCommand, err = scanText("filteringCommand: "); err != nil {
		return
	}
	if filteringCommand == "" {
		filteringCommand = "peco"
	}
	var unprotectiveCommand string
	if unprotectiveCommand, err = scanText("unprotectiveCommands: "); err != nil {
		return
	}
	if unprotectiveCommand != "" {
		unprotectiveCommands = strings.Split(unprotectiveCommand, ",")
	}
	fmt.Println()
	return
}
