package cmd

import (
	"github.com/vincentdaniel/subify/common/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	logger "github.com/spf13/jwalterweatherman"
	"github.com/vincentdaniel/subify/common/utils"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "subify",
	Short: "Tool to handle subtitles for your best TV Shows and movies",
	Long: `Tool to handle subtitles for your best TV Shows and movies
http://github.com/vincentdaniel/subify`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		utils.ExitPrintError(err, "An error occurred while running subify")
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	// Default configuration can be overridden
	RootCmd.PersistentFlags().StringVar(&config.ConfigFile, "config", "",
		"Config file (default is $HOME/.subify.json). Edit to change default behaviour")
	RootCmd.PersistentFlags().BoolVarP(&config.Verbose, "verbose", "v", false,
		"Print more information while executing")
	RootCmd.PersistentFlags().BoolVarP(&config.Dev, "dev", "", false,
		"Instanciate development sandbox instead of production variables")

}


const (
	SubifyConfigPath = "$HOME"
	SubifyConfigFile = ".subify"
)

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if config.ConfigFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(config.ConfigFile)
	}

	viper.SetConfigName(SubifyConfigFile) // name of config file (without extension)
	viper.AddConfigPath(SubifyConfigPath)   // adding home directory as first search path
	viper.AutomaticEnv()           // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err == nil {
		if config.Verbose {
			logger.INFO.Println("Using config file:", viper.ConfigFileUsed())
		}
	} else {
		logger.WARN.Println("Could not read config file (", viper.ConfigFileUsed(), "): ", err)
	}
}
