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

// Package env provide the implement to read property from ENV
package env

import (
	"context"
	"log/slog"
	"os"
	"regexp"
	"strings"

	"github.com/anyvoxel/airmid/anvil/logger"
)

// Loader is a loader for loading environment variables.
type Loader interface {
	Load() map[string]string
}

// loaderImpl is the implementation of Loader.
type envLoader struct {
	// prefix is used to filter airmid corresponding env
	prefix string
	// includePattern is a regex pattern for involving env vars
	includePattern string
	// excludePattern is a regex pattern for excluding env vars
	excludePattern string

	envs map[string]string
}

// EnvLoaderOptions ...
// nolint
type EnvLoaderOptions struct {
	// prefix is used to filter airmid corresponding env
	prefix string
	// includePattern is a regex pattern for involving env vars
	includePattern string
	// excludePattern is a regex pattern for excluding env vars
	excludePattern string
}

// EnvLoaderOption ...
// nolint
type EnvLoaderOption func(options *EnvLoaderOptions)

var defaultEnvLoaderOptions = EnvLoaderOptions{}

// NewEnvLoader returns a new loaderImpl instance which implements Loader.
func NewEnvLoader(opts ...EnvLoaderOption) Loader {
	options := defaultEnvLoaderOptions
	for _, o := range opts {
		o(&options)
	}
	return &envLoader{
		prefix:         options.prefix,
		includePattern: options.includePattern,
		excludePattern: options.excludePattern,
		envs:           make(map[string]string),
	}
}

// WithPrefixOption is an option to set env prefix.
func WithPrefixOption(prefix string) EnvLoaderOption {
	return func(options *EnvLoaderOptions) {
		options.prefix = prefix
	}
}

// WithEnvExcludePattern is an option to set env exclude pattern.
func WithEnvExcludePattern(excludePattern string) EnvLoaderOption {
	return func(options *EnvLoaderOptions) {
		options.excludePattern = excludePattern
	}
}

// WithEnvIncludePattern is an option to set env include pattern.
func WithEnvIncludePattern(includePattern string) EnvLoaderOption {
	return func(options *EnvLoaderOptions) {
		options.includePattern = includePattern
	}
}

func (l *envLoader) Load() map[string]string {
	envStrs := os.Environ()
	for _, envStr := range envStrs {
		kvArr := strings.SplitN(envStr, "=", 2)
		if len(kvArr) != 2 {
			continue
		}

		key := kvArr[0]
		if l.prefix != "" {
			if !strings.HasPrefix(key, l.prefix) {
				continue
			}
			key = strings.TrimPrefix(key, l.prefix)
		}

		if l.filterEnv(key) {
			continue
		}

		l.envs[key] = kvArr[1]
	}

	logger.FromContext(context.TODO()).DebugContext(
		context.TODO(),
		"loading airmid environments variables",
		slog.String("Prefix", l.prefix),
		slog.Any("envs", l.envs),
	)
	return l.envs
}

// filterEnv filter the env by following rules:
// 1. not start with specified prefix
// 2. exclude returns true
// 3. include returns false.
func (l *envLoader) filterEnv(key string) bool {
	includeFn := func(s string) bool {
		if l.includePattern == "" {
			return true
		}
		matched, err := regexp.MatchString(l.includePattern, s)
		if err != nil {
			panic(err)
		}
		return matched
	}
	excludeFn := func(s string) bool {
		if l.excludePattern == "" {
			return false
		}
		matched, err := regexp.MatchString(l.excludePattern, s)
		if err != nil {
			panic(err)
		}
		return matched
	}
	return excludeFn(key) || !includeFn(key)
}
