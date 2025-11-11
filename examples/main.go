package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	_ "github.com/anyvoxel/airmid/examples/runner1"
	_ "github.com/anyvoxel/airmid/examples/runner2"
	_ "github.com/anyvoxel/airmid/examples/runner3"

	"github.com/anyvoxel/airmid/app"
	slogctx "github.com/veqryn/slog-context"
)

// CurrentProjectPath get the project root path
func CurrentProjectPath() string {
	path := currentFilePath()

	ppath, err := filepath.Abs(filepath.Join(filepath.Dir(path), "../"))
	if err != nil {
		panic(fmt.Errorf("build current project path with %s failed, %w", path, err))
	}

	f, err := os.Stat(ppath)
	if err != nil {
		panic(fmt.Errorf("stat project path %s failed, %w", ppath, err))
	}

	if f.Mode()&os.ModeSymlink != 0 {
		fpath, err := os.Readlink(ppath)
		if err != nil {
			panic(fmt.Errorf("readlink from path %s failed, %w", fpath, err))
		}
		ppath = fpath
	}

	return ppath
}

func currentFilePath() string {
	_, file, _, _ := runtime.Caller(1)
	return file
}

func main() {
	configDirs := []string{
		CurrentProjectPath() + "/examples/config/",
		CurrentProjectPath() + "/examples/config/profiles",
	}
	os.Setenv("AIRMID_AIRMID_CONFIG_DIR", strings.Join(configDirs, ","))
	os.Setenv("AIRMID_AIRMID_PROFILES_ACTIVE", "default")

	err := app.Run(context.Background(), app.WithAttributes(
		app.Attribute{
			Key:   "module",
			Value: "1",
		},
		app.Attribute{
			Key:   "component",
			Value: "2",
		},
		app.Attribute{
			Key:   "version",
			Value: "3",
		},
		app.Attribute{
			Key:   "go_version",
			Value: "4",
		},
		app.Attribute{
			Key:   "branch",
			Value: "5",
		},
		app.Attribute{
			Key:   "rivision",
			Value: "6",
		},
		app.Attribute{
			Key:   "build_date",
			Value: "7",
		},
	))
	if err != nil {
		slogctx.FromCtx(context.TODO()).ErrorContext(context.Background(), "Run application failed", slog.Any("Error", err))
		panic(err)
	}
}
