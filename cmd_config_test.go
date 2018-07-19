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

	if _, err = os.Stat(path); err != nil {
		t.Error(err)
	}
	p, err := ioutil.ReadFile(path)
	if err != nil {
		t.Error(err)
	}
	var cfg Config
	if err = yaml.Unmarshal(p, &cfg); err != nil {
		t.Error(err)
	}

	if cfg.KeyFile != filepath.Join("testdata", "key_file") {
		t.Errorf("KeyFileu Unmarshal failure: %#v", cfg.KeyFile)
	}
	if cfg.DataFile != filepath.Join("testdata", "data.dat") {
		t.Errorf("DataFile Unmarshal failure: %#v", cfg.DataFile)
	}
	if cfg.FilteringCommand != "peco" {
		t.Errorf("FilteringCommand Unmarshal failure: %#v", cfg.FilteringCommand)
	}
	if len(cfg.UnprotectiveCommands) != 2 {
		t.Errorf("UnprotectiveCommands Unmarshal failure: %#v", len(cfg.UnprotectiveCommands))
		return
	}
	if cfg.UnprotectiveCommands[0] != "new" {
		t.Errorf("UnprotectiveCommands Unmarshal failure: %#v", cfg.UnprotectiveCommands)
	}
	if cfg.UnprotectiveCommands[1] != "copy" {
		t.Errorf("UnprotectiveCommands Unmarshal failure: %#v", cfg.UnprotectiveCommands)
	}
}

func TestConfigTemp(t *testing.T) {
	err := configTemp()
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
