package handlers

import (
	"bytes"
	"context"

	command "github.com/raba-jp/primus/pkg/operations/command/handlers"
)

func cmdArgs(ctx context.Context, executable command.ExecutableHandler, ct cmdType, cmds []string) (string, []string) {
	cmd := "pacman"
	var (
		options []string
	)
	yay := usableYay(ctx, executable)
	if yay {
		cmd = "yay"
	}

	switch ct {
	case install:
		powerpill := usablePowerpill(ctx, executable)
		if powerpill && yay {
			opts := []string{"--pacman", "powerpill", "-S", "--noconfirm"}
			options = make([]string, 0, len(opts)+len(cmds))
			options = append(options, opts...)
		} else {
			opts := []string{"-S", "--noconfirm"}
			options = make([]string, 0, len(opts)+len(cmds))
			options = append(options, opts...)
		}
		options = append(options, cmds...)
	case uninstall:
		options = []string{"-R", "--noconfirm"}
		options = append(options, cmds...)
	}

	return cmd, options
}

func sprintCmd(cmd string, options []string) string {
	buf := new(bytes.Buffer)
	buf.WriteString(cmd)
	for _, op := range options {
		buf.WriteString(" " + op)
	}
	buf.WriteString("\n")
	return buf.String()
}

func usableYay(ctx context.Context, executable command.ExecutableHandler) bool {
	return executable.Run(ctx, "yay")
}

func usablePowerpill(ctx context.Context, executable command.ExecutableHandler) bool {
	return executable.Run(ctx, "powerpill")
}
