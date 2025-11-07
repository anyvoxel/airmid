// Package reader provide the function to read resource content
package reader

import (
	"github.com/anyvoxel/airmid/xerrors"
)

// Reader is the props-config file reader.
type Reader interface {
	// Read will unmarshal the data
	Read(data []byte) (map[string]any, error)

	// Match will check the Reader can process the file
	Match(filename string) error

	// Name return the reader name
	Name() string
}

var (
	readers = []Reader{}
)

// RegisterReader register the reader to slice.
func RegisterReader(r Reader) error {
	for _, ri := range readers {
		if ri == r {
			return xerrors.WrapDuplicate(
				"RegisterReader duplicate: reader '%v' is already exists with '%v'", r.Name(), ri.Name())
		}
	}

	readers = append(readers, r)
	return nil
}

// RegisterExtFileReader register the ext file reader to slice.
func RegisterExtFileReader(fn func(data []byte) (map[string]any, error), exts ...string) error {
	if len(exts) == 0 {
		return xerrors.Errorf("RegisterExtFileReader: exts cannot be empty")
	}

	return RegisterReader(&extReader{
		exts: exts,
		fn:   fn,
	})
}

// Read will unmarshal the data.
func Read(filename string, data []byte) (map[string]any, error) {
	for _, r := range readers {
		if err := r.Match(filename); err != nil {
			continue
		}

		return r.Read(data)
	}

	return nil, xerrors.Errorf("Cannot found reader for '%v'", filename)
}
