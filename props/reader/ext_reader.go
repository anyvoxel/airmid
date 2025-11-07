package reader

import (
	"strings"

	"github.com/anyvoxel/airmid/xerrors"
)

type extReader struct {
	exts []string

	fn func(data []byte) (map[string]any, error)
}

func (r *extReader) Read(data []byte) (map[string]any, error) {
	return r.fn(data)
}

func (r *extReader) Match(filename string) error {
	for _, ext := range r.exts {
		if strings.HasSuffix(filename, ext) {
			return nil
		}
	}

	return xerrors.Errorf("'%v' cannot support filename '%v'", r.Name(), filename)
}

func (r *extReader) Name() string {
	return strings.Join(r.exts, ",")
}
