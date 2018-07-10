package main

import (
	"fmt"
	"os"
	"runtime"
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
	err := configTemp()
	if err != nil {
		return err
	}
	keyFile, dataFile, filteringCommand, unprotectiveCommands, err := configScan()

	initCfg := Config{
		KeyFile:              keyFile,
		DataFile:             dataFile,
		FilteringCommand:     filteringCommand,
		UnprotectiveCommands: unprotectiveCommands,
	}
	path, err := app.ConfigFile("config.yml")
	if err != nil {
		return err
	}

	err = initCfg.configSave(path)
	if err != nil {
		return err
	}
	PrintSuccess(ctx.out, "config file saved as '%s' successfully", path)
	return nil
}

// Save config to file on given path.
func (cfg Config) configSave(path string) error {
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

// Temp config file create
func configTemp() error {
	dir, err := app.ConfigDir()
	if err != nil {
		return err
	}
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	path, err := app.ConfigFile("config.yml")
	if err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}

func configScan() (keyFile string, dataFile string, filteringCommand string, unprotectiveCommands []string, err error) {
	if keyFile, err = scanText("keyFile: "); err != nil {
		return
	}
	if keyFile == "" {
		if runtime.GOOS == "windows" {
			keyFile = os.Getenv("USERPROFILE") + "\\.ssh\\id_rsa"
		} else {
			keyFile = "~/.ssh/id_rsa"
		}
	}
	if dataFile, err = scanText("dataFile: "); err != nil {
		return
	}
	if dataFile == "" {
		if runtime.GOOS == "windows" {
			dataFile = os.Getenv("USERPROFILE") + "\\.local\\share\\spwd\\data.dat"
		} else {
			dataFile = "~/.local/share/spwd/data.dat"
		}
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
		uc := strings.Split(unprotectiveCommand, ",")
		for _, cmd := range uc {
			if str := strings.TrimSpace(cmd); str != "" {
				unprotectiveCommands = append(unprotectiveCommands, str)
			}
		}
	}
	fmt.Println()
	return
}
