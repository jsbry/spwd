package main

import (
	"path/filepath"
	"testing"
)

func TestConfigScan(t *testing.T) {
	keyFile, dataFile, filteringCommand, unprotectiveCommands, err := configScan()
	if err != nil {
		t.Error(err)
	}
	if keyFile != "~/.ssh/id_rsa" {
		t.Error(err)
	}
	if dataFile != "~/.local/share/spwd/data.dat" {
		t.Error(err)
	}
	if filteringCommand != "peco" {
		t.Error(err)
	}
	if len(unprotectiveCommands) != 0 {
		t.Error(err)
	}
}

func TestCmdConfig(t *testing.T) {
	testConfig := Config{
		KeyFile:              filepath.Join("testdata", "key_file"),
		DataFile:             filepath.Join("testdata", "data.dat"),
		FilteringCommand:     "peco",
		UnprotectiveCommands: []string{"new", "copy"},
	}

	path := filepath.Join("testdata", "test_config.yml")
	err := testConfig.configSave(path)
	if err != nil {
		t.Error(err)
	}
}
