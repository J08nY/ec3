package db

import (
	"bufio"
	"io"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/xerrors"

	"github.com/mmcloughlin/ec3/efd"
	"github.com/mmcloughlin/ec3/efd/op3/parse"
)

type Key struct {
	Path           string
	Class          string
	Section        string
	Shape          string
	Representation string
	Operation      string
	Name           string
	Ext            string
}

func KeyFromFilename(filename string) Key {
	k := Key{Path: filename}
	k.Ext = filepath.Ext(filename)
	path := strings.TrimSuffix(filename, k.Ext)
	parts := strings.Split(path, "/")
	n := len(parts)
	dst := []*string{&k.Class, &k.Section, &k.Shape, &k.Representation, &k.Operation}
	for i := 0; i < n-1 && i < len(dst); i++ {
		*dst[i] = parts[i]
	}
	k.Name = parts[n-1]
	return k
}

func (k Key) IsOP3() bool { return k.Ext == ".op3" }

func (k Key) IsFormula() bool { return k.Operation != "" }

func (k Key) IsShape() bool { return k.Shape != "" && k.Name == "coordinates" }

func (k Key) IsRepresentation() bool { return k.Representation != "" && k.Name == "variables" }

func (k Key) ShapeID() string { return path.Join(k.Class, k.Shape) }

func (k Key) RepresentationID() string { return path.Join(k.ShapeID(), k.Representation) }

func (k Key) OperationID() string { return path.Join(k.RepresentationID(), k.Operation) }

func (k Key) FormulaID() string { return path.Join(k.OperationID(), k.Name) }

type Database struct {
	Shapes          map[string]*efd.Shape
	Representations map[string]*efd.Representation
	Formulae        map[string]*efd.Formula
}

func New() *Database {
	return &Database{
		Shapes:          map[string]*efd.Shape{},
		Representations: map[string]*efd.Representation{},
		Formulae:        map[string]*efd.Formula{},
	}
}

func (d Database) shape(k string) *efd.Shape {
	if _, ok := d.Shapes[k]; !ok {
		d.Shapes[k] = &efd.Shape{
			ID: k,
		}
	}
	return d.Shapes[k]
}

func (d Database) representation(k string) *efd.Representation {
	if _, ok := d.Representations[k]; !ok {
		d.Representations[k] = &efd.Representation{
			ID: k,
		}
	}
	return d.Representations[k]
}

func (d Database) formula(k string) *efd.Formula {
	if _, ok := d.Formulae[k]; !ok {
		d.Formulae[k] = &efd.Formula{
			ID: k,
		}
	}
	return d.Formulae[k]
}

func Read(archive string) (*Database, error) {
	p := parser{
		DB: New(),
	}
	if err := Walk(archive, p); err != nil {
		return nil, err
	}
	return p.DB, nil
}

type parser struct {
	DB *Database
}

func (p parser) Visit(filename string, r io.Reader) error {
	k := KeyFromFilename(filename)
	switch {
	case k.IsOP3():
		return p.op3(k, r)
	case k.IsShape():
		return p.shape(k, r)
	case k.IsRepresentation():
		return p.representation(k, r)
	case k.IsFormula():
		return p.formula(k, r)
	case k.Name == "README":
		// pass
	default:
		return xerrors.Errorf("unknown file: %s", filename)
	}
	return nil
}

func (p parser) formula(k Key, r io.Reader) error {
	f := p.DB.formula(k.FormulaID())
	f.Tag = k.Name
	f.Class = k.Class
	f.Shape = p.DB.shape(k.ShapeID())
	f.Representation = p.DB.representation(k.RepresentationID())
	f.Operation = k.Operation

	props, err := ParseProperties(r)
	if err != nil {
		return err
	}

	for prop, vs := range props {
		switch prop {
		case "source":
			f.Source, err = atmostone(prop, vs)
		case "compute":
			f.Compute = vs
		case "assume":
			f.Assume = vs
		case "parameter":
			f.Parameters = vs
		case "appliesto":
			f.AppliesTo, err = atmostone(prop, vs)
		default:
			return xerrors.Errorf("unknown property %q", prop)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (p parser) representation(k Key, r io.Reader) error {
	repr := p.DB.representation(k.RepresentationID())
	repr.Tag = k.Representation
	repr.Class = k.Class
	repr.Shape = p.DB.shape(k.ShapeID())

	props, err := ParseProperties(r)
	if err != nil {
		return err
	}

	for prop, vs := range props {
		switch prop {
		case "name":
			repr.Name, err = exactlyone(prop, vs)
		case "assume":
			repr.Assume = vs
		case "parameter":
			repr.Parameters = vs
		case "variable":
			repr.Variables = vs
		case "satisfying":
			repr.Satisfying = vs
		default:
			return xerrors.Errorf("unknown property %q", prop)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (p parser) shape(k Key, r io.Reader) error {
	s := p.DB.shape(k.ShapeID())
	s.Tag = k.Shape
	s.Class = k.Class

	props, err := ParseProperties(r)
	if err != nil {
		return err
	}

	for prop, vs := range props {
		switch prop {
		case "name":
			s.Name, err = exactlyone(prop, vs)
		case "coordinate":
			s.Coordinates = vs
		case "a0", "a1", "a2", "a3", "a4", "a6":
			i := int(prop[1] - '0')
			s.A[i], err = exactlyone(prop, vs)
		case "satisfying":
			s.Satisfying = vs
		case "parameter":
			s.Parameters = vs
		case "addition":
			s.Addition = vs
		case "doubling":
			s.Doubling = vs
		case "negation":
			s.Negation = vs
		case "neutral":
			s.Neutral = vs
		case "fromweierstrass":
			s.FromWeierstrass = vs
		case "toweierstrass":
			s.ToWeierstrass = vs
		default:
			return xerrors.Errorf("unknown property %q", prop)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (p parser) op3(k Key, r io.Reader) error {
	prog, err := parse.Reader(k.Path, r)
	// Note we expect some files to fail parsing, so we supress errors here.
	if err != nil {
		return nil
	}

	f := p.DB.formula(k.FormulaID())
	f.Program = prog
	return nil
}

var whitespace = regexp.MustCompile("[[:space:]]+")

func ParseProperties(r io.Reader) (map[string][]string, error) {
	properties := map[string][]string{}
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		parts := whitespace.Split(line, 2)
		if len(parts) != 2 || parts[0] == "" {
			continue
		}
		key, value := parts[0], parts[1]
		properties[key] = append(properties[key], value)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return properties, nil
}

func exactlyone(prop string, vs []string) (string, error) {
	if len(vs) != 1 {
		return "", xerrors.Errorf("expected exactly one value for %q", prop)
	}
	return vs[0], nil
}

func atmostone(prop string, vs []string) (string, error) {
	if len(vs) > 1 {
		return "", xerrors.Errorf("expected at most one value for %q", prop)
	}
	if len(vs) == 1 {
		return vs[0], nil
	}
	return "", nil
}