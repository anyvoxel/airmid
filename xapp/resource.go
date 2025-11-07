package xapp

import (
	"context"
	"io"
	"log/slog"

	"github.com/anyvoxel/airmid/logger"
	"github.com/anyvoxel/airmid/props"
	"github.com/anyvoxel/airmid/props/reader"
)

// Config holds the config resources.
type config struct {
	resourceLocator  ResourceLocator `airmid:"autowire:?"`
	ConfigExtensions []string        `airmid:"value:${airmid.config.extensions:=.yaml,.yml}"`
	ActiveProfiles   []string        `airmid:"value:${airmid.profiles.active:=}"`
}

func (*config) NewConfig() (*config, error) {
	return &config{}, nil
}

//nolint:revive,cyclop
func (c *config) loadProperty(p props.Properties) error {
	resources := []Resource{}

	logger.FromContext(context.TODO()).DebugContext(
		context.TODO(),
		"Configuration file extensions supported",
		slog.Any("ConfigExtensions", c.ConfigExtensions),
	)
	for _, ext := range c.ConfigExtensions {
		filename := "application" + ext
		ress, err := c.resourceLocator.Locate(filename)
		if err != nil {
			return err
		}

		resources = append(resources, ress...)
	}

	logger.FromContext(context.TODO()).DebugContext(
		context.TODO(),
		"Configuration file active profiles supported",
		slog.Any("ActiveProfiles", c.ActiveProfiles),
	)
	for _, profile := range c.ActiveProfiles {
		for _, ext := range c.ConfigExtensions {
			filename := "application-" + profile + ext
			ress, err := c.resourceLocator.Locate(filename)
			if err != nil {
				return err
			}

			resources = append(resources, ress...)
		}
	}

	for _, res := range resources {
		logger.FromContext(context.TODO()).DebugContext(
			context.TODO(),
			"Loading configuration properties",
			slog.String("FileName", res.Name()),
		)
		data, err := io.ReadAll(res)
		if err != nil {
			return err
		}

		objs, err := reader.Read(res.Name(), data)
		if err != nil {
			return err
		}

		for k, v := range objs {
			err := p.Set(k, v)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
