package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	yaml "gopkg.in/yaml.v2"
)

func TestCmdList(t *testing.T) {
	teardown, err := setupTestData()
	if err != nil {
		if teardown != nil {
			teardown()
		}
	}
	defer teardown()
	out := &bytes.Buffer{}
	ctx := newContext(out)
	err = cmdList.Run(ctx, []string{})
	if err != nil {
		t.Error(err)
	}
	expected := "  NAME           DESCRIPTION           \n ------ ------------------------------- \n  foo   This is test password          \n  bar   This is another test password  \n"
	if got := out.String(); got != expected {
		t.Errorf("exptected: \n%s\n\nbut got: \n%s", expected, got)
		fmt.Println(len(expected), len(got))
	}
}

func setupTestData() (func(), error) {
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
	copyTestFile("data.yml", td)
	copyTestFile("identity_file", td)
	cfg := Config{
		IdentityFile: filepath.Join(td, "identity_file"),
		DataFile:     filepath.Join(td, "data.yml"),
	}
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

func copyTestFile(name string, dstDir string) error {
	src, err := os.Open(filepath.Join("testdata", name))
	if err != nil {
		return err
	}
	dst, err := os.Create(filepath.Join(dstDir, name))
	if err != nil {
		return err
	}
	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}
	return nil
}