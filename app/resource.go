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
	"context"
	"io"
	"log/slog"

	"github.com/anyvoxel/airmid/app/reader"
	"github.com/anyvoxel/airmid/ioc/props"
	slogctx "github.com/veqryn/slog-context"
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

	slogctx.FromCtx(context.TODO()).DebugContext(
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

	slogctx.FromCtx(context.TODO()).DebugContext(
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
		slogctx.FromCtx(context.TODO()).DebugContext(
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
