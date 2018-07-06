package main

import (
	"bytes"
	"path/filepath"
	"testing"
)

func TestRunConfig(t *testing.T) {
	out := &bytes.Buffer{}
	ctx := newContext(out, "config")
	err := cmdConfig.Run(ctx, []string{})
	if err != nil {
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
