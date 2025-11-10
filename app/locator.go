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

package app

import (
	"io"
	"os"
	"path/filepath"
)

// Resource is the reader & name.
type Resource interface {
	io.Reader
	Name() string
}

// ResourceLocator will locate the filename.
type ResourceLocator interface {
	// Locate locate file by filename
	Locate(filename string) ([]Resource, error)
}

// localResourceLocator is the implementation of ResourceLocator.
type localResourceLocator struct {
	configDir []string `airmid:"value:${airmid.config.dir:=config/}"`
}

func (l *localResourceLocator) Locate(filename string) ([]Resource, error) {
	resources := make([]Resource, 0, 4)

	for _, dir := range l.configDir {
		filename, err := filepath.Abs(filepath.Join(dir, filename))
		if err != nil {
			return nil, err
		}

		file, err := os.Open(filename)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}

			return nil, err
		}
		resources = append(resources, file)
	}

	return resources, nil
}
