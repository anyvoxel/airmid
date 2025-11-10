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
	"os"
	"strings"

	"github.com/anyvoxel/airmid/ioc/props"
	"github.com/anyvoxel/airmid/xapp/env"
)

// PropertiesLoader is a loader for loading properties.
type PropertiesLoader interface {
	// LoadProperties loads properties to p
	LoadProperties(p props.Properties) error
}

// PropertiesLoaders holds a set of loader.
type PropertiesLoaders interface {
	PropertiesLoader
	Add(p PropertiesLoader) PropertiesLoaders
}

// PropertiesLoader load convert/load env to properties.
type envPropertiesLoader struct {
	envLoader    env.Loader
	keyConvertFn func(envKey string) string
}

// NewEnvPropertiesLoader return a instance of envPropertiesLoader which implements Properties.
func NewEnvPropertiesLoader(prefix string, keyConvertFn func(string) string,
) PropertiesLoader {
	envIncludePattern := os.Getenv("AIRMID_INCLUDE_ENV_PATTERNS")
	envExcludePattern := os.Getenv("AIRMID_EXCLUDE_ENV_PATTERNS")

	return &envPropertiesLoader{
		keyConvertFn: keyConvertFn,
		envLoader: env.NewEnvLoader(
			env.WithPrefixOption(prefix),
			env.WithEnvIncludePattern(envIncludePattern),
			env.WithEnvExcludePattern(envExcludePattern),
		),
	}
}

// LoadProperties loads properties from environment variables.
func (l *envPropertiesLoader) LoadProperties(p props.Properties) error {
	for k, v := range l.envLoader.Load() {
		key := l.keyConvertFn(k)
		// if the env value is array, set to properties as key[0]=value0, key[1]=value1.
		if values := strings.Split(v, ","); len(values) > 1 {
			if err := p.Set(key, values); err != nil {
				return err
			}
		} else {
			// if the env value is a string or an array which only has one element, set to properties as key=value.
			if err := p.Set(key, v); err != nil {
				return err
			}
		}
	}
	return nil
}

// DefaultEnvKeyConvertFunc convert the env key to prop key by following rules:
// 1. to lowercase
// 2. replace '_' to '.'.
var DefaultEnvKeyConvertFunc = func(envKey string) string {
	return strings.ReplaceAll(strings.ToLower(envKey), "_", ".")
}

// OptionArgsPropertiesLoader is used for loading args properties.
type optionArgsPropertiesLoader struct{}

// NewOptionArgsPropertiesLoader returns a new instance of optionArgsPropertiesLoader.
func NewOptionArgsPropertiesLoader() PropertiesLoader {
	return &optionArgsPropertiesLoader{}
}

// LoadProperties loads properties from command line args.
//
//nolint:revive
func (o *optionArgsPropertiesLoader) LoadProperties(p props.Properties) error {
	for _, arg := range os.Args[1:] {
		// TODO: support non-optional args, such as "--airmid.application.port 8080 -flag false"
		parts := strings.SplitN(arg, "=", 2)
		switch len(parts) {
		case 1:
			key := strings.TrimLeft(parts[0], "-")
			if len(key) == 0 {
				continue
			}
			if err := p.Set(key, true); err != nil {
				return err
			}
		case 2:
			key := strings.TrimLeft(parts[0], "-")
			if len(key) == 0 {
				continue
			}
			rawValue := parts[1]
			// if the flag value is array, set to properties as key[0]=value0, key[1]=value1.
			if values := strings.Split(rawValue, ","); len(values) > 1 {
				if err := p.Set(key, values); err != nil {
					return err
				}
			} else {
				// if the flag value is a string or an array which only has one element, set to properties as key=value.
				if err := p.Set(key, rawValue); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

type propertiesLoaders struct {
	loaders []PropertiesLoader
}

// NewPropertiesLoaders holds a set of property loads to load properties in the order that the loads have been set.
func NewPropertiesLoaders() PropertiesLoaders {
	return &propertiesLoaders{
		loaders: make([]PropertiesLoader, 0),
	}
}

// LoadProperties load properties in following order:
// 1. flags
// 2. environment variables
// 3. configurations.
func (l *propertiesLoaders) LoadProperties(p props.Properties) error {
	for _, loader := range l.loaders {
		if err := loader.LoadProperties(p); err != nil {
			return err
		}
	}
	return nil
}

func (l *propertiesLoaders) Add(p PropertiesLoader) PropertiesLoaders {
	l.loaders = append(l.loaders, p)
	return l
}
