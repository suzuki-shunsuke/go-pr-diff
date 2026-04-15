package prdiff

import (
	"context"
	"os"
	"os/exec"
	"time"
)

const defaultWaitDelay = 1 * time.Minute

func setCancel(cmd *exec.Cmd) {
	cmd.Cancel = func() error {
		return cmd.Process.Signal(os.Interrupt)
	}
	cmd.WaitDelay = defaultWaitDelay
}

func command(ctx context.Context, name string, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, name, args...)
	setCancel(cmd)
	return cmd
}
