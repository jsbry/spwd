package main

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// Item is password item.
type Item struct {
	Name string `yaml:"name"`
	// Description is free description of password.
	Description string `yaml:"description"`
	// Password text.
	Password string `yaml:"password"`
}

// NewItem returns new item that initialized give values.
func NewItem(name string, desc string, enc string) Item {
	return Item{
		Name:        name,
		Description: desc,
		Password:    enc,
	}
}

// Items is item list.
type Items []Item

// Find item that has given keyword.
func (is Items) Find(name string) *Item {
	for _, i := range is {
		if name == i.Name {
			return &i
		}
	}
	return nil
}

// LoadItems load items from file on given path.
func LoadItems(key []byte, path string) (Items, error) {
	_, err := os.Stat(path)
	if err != nil {
		return Items{}, nil
	}

	p, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	dec, err := Decrypt(key, string(p))
	if err != nil {
		return nil, err
	}
	var is Items
	err = yaml.Unmarshal(dec, &is)
	if err != nil {
		return nil, err
	}
	return is, nil
}

// Save items to file on given path.
func (is Items) Save(key []byte, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	p, err := yaml.Marshal(is)
	if err != nil {
		return err
	}
	enc, err := Encrypt(key, p)
	if err != nil {
		return err
	}
	f.WriteString(enc)
	return nil
}

// Update same name item with given item.
func (is Items) Update(nit Item) Items {
	nis := Items(make([]Item, len(is)))
	for i, it := range is {
		if it.Name == nit.Name {
			nis[i] = nit
		} else {
			nis[i] = it
		}
	}
	return nis
}

// ToDataTable returns data for tablewriter.
func (is Items) ToDataTable() [][]string {
	data := make([][]string, len(is))
	for i, it := range is {
		data[i] = []string{it.Name, it.Description}
	}
	return data
}
