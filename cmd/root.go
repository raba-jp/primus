package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "starlark_iac",
	Short: "configuration management tool",
	Long:  "configuration management tool",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.starlark_iac.yaml)")

	rootCmd.Flags().BoolP("loglevel", "l", false, "Log Level")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			zap.L().Error("Failed to get home directory", zap.Error(err))
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".starlark_iac")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		zap.L().Info("Using config file", zap.String("path", viper.ConfigFileUsed()))
	}
}
