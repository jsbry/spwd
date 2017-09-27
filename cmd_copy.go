package main

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/atotto/clipboard"
)

var cmdCopy = &Command{
	Run:       runCopy,
	UsageLine: "copy",
	Short:     "Copy password to clipboard",
	Long:      ``,
}

func runCopy(ctx context, args []string) error {
	if len(args) == 0 {
		return errors.New("item name required")
	}
	cfg, err := GetConfig()
	if err != nil {
		return err
	}
	Initialize(cfg)
	is, err := LoadItems(cfg.DataFile)
	if err != nil {
		return err
	}
	it := is.Find(args[0])
	if it == nil {
		return fmt.Errorf("item not found: %s", args[0])
	}

	keySrc, err := ioutil.ReadFile(cfg.IdentityFile)
	if err != nil {
		return err
	}
	key := GenKey(keySrc)
	dec, err := Decode(it.Encrypted)
	if err != nil {
		return err
	}
	pwd, err := Decrypt(key, dec)
	if err != nil {
		return err
	}
	clipboard.WriteAll(pwd)
	fmt.Fprintln(ctx.out, fmt.Sprintf("password of '%s' copy to clipboard", it.Name))
	return nil
}
