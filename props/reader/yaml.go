package reader

import (
	"gopkg.in/yaml.v2"

	"github.com/anyvoxel/airmid/xerrors"
)

func yamlRead(data []byte) (map[string]any, error) {
	m := make(map[string]any)
	err := yaml.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func init() {
	err := RegisterExtFileReader(yamlRead, ".yaml", ".yml")
	if err != nil {
		panic(xerrors.Wrapf(err, "Register yaml reader"))
	}
}
