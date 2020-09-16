// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package cmd

import (
	"github.com/spf13/cobra"
)

// Injectors from wire.go:

func Initialize() *cobra.Command {
	planCommand := NewPlanCommand()
	applyCommand := NewApplyCommand()
	versionCommand := NewVersionCommand()
	command := NewCommand(planCommand, applyCommand, versionCommand)
	return command
}