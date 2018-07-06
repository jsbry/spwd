package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	yaml "gopkg.in/yaml.v2"
)

func TestCmdConfig(t *testing.T) {
	teardown, err := setupTestConfig()
	if err != nil {
		if teardown != nil {
			teardown()
		}
	}
	defer teardown()
	out := &bytes.Buffer{}
	ctx := newContext(out, "config")
	err = cmdConfig.Run(ctx, []string{})
	if err != nil {
		t.Error(err)
	}
}

func TestSaveConfig(t *testing.T) {
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

func setupTestConfig() (func(), error) {
	td, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, err
	}
	f := func() {
		if _, err = os.Stat(td); err == nil {
			os.Remove(td)
		}
	}

	os.Setenv("XDG_CONFIG_HOME", td)
	os.Mkdir(filepath.Join(td, "spwd"), 0755)
	cfg := Config{}
	cfp, err := os.Create(filepath.Join(td, "spwd", "config.yml"))
	defer cfp.Close()
	if err != nil {
		return f, err
	}
	p, err := yaml.Marshal(cfg)
	if err != nil {
		return f, err
	}
	cfp.Write(p)
	return f, nil
}
