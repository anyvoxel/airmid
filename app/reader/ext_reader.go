// Copyright (c) 2025 The anyvoxel Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package reader

import (
	"strings"

	"github.com/anyvoxel/airmid/anvil/xerrors"
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
