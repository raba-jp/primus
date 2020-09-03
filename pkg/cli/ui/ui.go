package ui

import (
	"io"

	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
)

type UI interface {
	Printf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

func Printf(format string, args ...interface{}) {
	defaultUI.Printf(format, args...)
}

func Infof(format string, args ...interface{}) {
	defaultUI.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	defaultUI.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	defaultUI.Errorf(format, args...)
}

type CommandLine struct {
	Out    io.Writer
	Errout io.Writer
}

func (c *CommandLine) Printf(format string, args ...interface{}) {
	color.New(color.FgHiBlue).Fprintf(c.Out, format, args...)
}

func (c *CommandLine) Infof(format string, args ...interface{}) {
	color.New(color.FgGreen).Fprintf(c.Out, format, args...)
}

func (c *CommandLine) Warnf(format string, args ...interface{}) {
	color.New(color.FgYellow).Fprintf(c.Out, format, args...)
}

func (c *CommandLine) Errorf(format string, args ...interface{}) {
	color.New(color.FgRed).Fprintf(c.Errout, format, args...)
}

var defaultUI UI = &CommandLine{Out: colorable.NewColorableStdout(), Errout: colorable.NewColorableStderr()}

func SetDefaultUI(newUI UI) {
	defaultUI = newUI
}
