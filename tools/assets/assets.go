package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"

	"golang.org/x/xerrors"

	"github.com/mmcloughlin/ec3/internal/gocode"
)

type Asset struct {
	Name string
	Data []byte
}

func LoadAssets(filenames []string) ([]Asset, error) {
	var assets []Asset
	for _, filename := range filenames {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		assets = append(assets, Asset{
			Name: filename,
			Data: data,
		})
	}
	return assets, nil
}

type Config struct {
	GeneratedBy string
	PackageName string
	Function    string
}

func (cfg *Config) Validate() error {
	for _, field := range []struct {
		Value, Description string
	}{
		{cfg.GeneratedBy, "generated by"},
		{cfg.PackageName, "package name"},
		{cfg.Function, "function name"},
	} {
		if field.Value == "" {
			return xerrors.Errorf("missing %s", field.Description)
		}
	}
	return nil
}

func Generate(cfg *Config, assets []Asset) ([]byte, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	p := gocode.NewGenerator()

	// Header.
	p.CodeGenerationWarning(cfg.GeneratedBy)
	p.Package(cfg.PackageName)

	p.Import("fmt")

	// Function.
	p.Printf("func %s(name string) ([]byte, error)", cfg.Function)
	p.EnterBlock()
	p.Printf("switch name")
	p.EnterBlock()

	for _, asset := range assets {
		p.Linef("case %q:", asset.Name)
		p.Linef("return []byte(%s), nil", Quote(asset.Data))
		p.NL()
	}

	p.Linef("default:")
	p.Linef("return nil, fmt.Errorf(\"unknown asset %%s\", name)")

	p.LeaveBlock()
	p.LeaveBlock()

	return p.Formatted()
}

func Quote(b []byte) string {
	// Use a raw string if possible.
	const backtick = '`'
	if !bytes.ContainsRune(b, backtick) {
		return fmt.Sprintf("%c%s%c", backtick, b, backtick)
	}

	// Fallback to a double-quoted string.
	return strconv.Quote(string(b))
}
