package xapp

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
