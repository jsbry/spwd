package main

import (
	"path/filepath"
	"testing"
)

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
