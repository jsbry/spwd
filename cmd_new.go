package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

var cmdNew = &Command{
	Run:       runNew,
	UsageLine: "new",
	Short:     "Register new password item interactively",
	Long: `Register new password item to data file interactively.
If input name already exists, you can update it.
`,
}

func runNew(ctx context, args []string) error {
	cfg, err := GetConfig()
	if err != nil {
		return err
	}
	Initialize(cfg)
	key, err := GetKey(cfg.KeyFile)
	if err != nil {
		return err
	}
	is, err := LoadItems(key, cfg.DataFile)
	if err != nil {
		return err
	}

	name, desc, pwd, err := scan()
	if err != nil {
		return err
	}
	nit := NewItem(name, desc, pwd)
	if err != nil {
		return err
	}
	if it := is.Find(name); it != nil {
		b, berr := scanBool(fmt.Sprintf("item '%s' already exists, update? [y/N]: ", name))
		if berr != nil {
			return berr
		}
		if !b {
			return nil
		}
		is = is.Update(nit)
	} else {
		is = append(is, nit)
	}

	err = is.Save(key, cfg.DataFile)
	if err != nil {
		return err
	}
	PrintSuccess(ctx.out, "password of '%s' is saved successfully", name)
	return nil
}

func scan() (name string, desc string, pwd string, err error) {
	name, err = scanText("Name: ")
	if name == "" {
		err = errors.New("name is required")
		return
	}
	desc, err = scanText("Description: ")
	pwd, err = scanPassword("Password: ")
	if pwd == "" {
		err = errors.New("password is required")
	}
	fmt.Println()
	return
}

func scanText(prompt string) (string, error) {
	if prompt != "" {
		fmt.Print(prompt)
	}
	in := bufio.NewScanner(os.Stdin)
	in.Scan()
	if err := in.Err(); err != nil {
		return "", err
	}
	return in.Text(), nil
}

func scanBool(prompt string) (bool, error) {
	if prompt != "" {
		fmt.Print(prompt)
	}
	in := bufio.NewScanner(os.Stdin)
	in.Scan()
	if err := in.Err(); err != nil {
		return false, err
	}
	inText := strings.ToLower(in.Text())

	return inText == "y" || inText == "yes", nil
}

func scanPassword(prompt string) (string, error) {
	if prompt != "" {
		fmt.Print(prompt)
	}
	p, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	return string(p), nil
}
