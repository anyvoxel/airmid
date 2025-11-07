package utils

import (
	"context"
	"log/slog"
	"runtime/debug"

	"github.com/anyvoxel/airmid/xerrors"
)

// SafeRun will execute cmd and recover any panic.
func SafeRun(cmd func()) {
	defer func() {
		if r := recover(); r != nil {
			var err error
			if ierr, ok := r.(error); ok {
				err = ierr
			} else {
				err = xerrors.Errorf("Recover from: '%v', stack: '%v'", r, string(debug.Stack()))
			}

			slog.ErrorContext(
				context.TODO(), "Panic when execute runnable",
				slog.Any("Error", err),
			)
		}
	}()

	cmd()
}
