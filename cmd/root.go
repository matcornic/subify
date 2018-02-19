package cmd

import (
	"fmt"

	"github.com/matcornic/subify/common/config"
	"github.com/matcornic/subify/common/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "subify",
	Short: "Tool to handle subtitles for your best TV Shows and movies",
	Long: `Tool to handle subtitles for your best TV Shows and movies
http://github.com/matcornic/subify`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Overwrite conf from config files
		config.Dev = viper.GetBool("root.dev")
		config.Verbose = viper.GetBool("root.verbose")
		utils.InitLoggingConf()
	},
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
		"Config file (default is $HOME/.subify.yaml|json|toml). Edit to change default behavior")
	RootCmd.PersistentFlags().BoolVarP(&config.Verbose, "verbose", "v", false,
		"Print more information while executing")
	RootCmd.PersistentFlags().BoolVarP(&config.Dev, "dev", "", false,
		"Instanciate development sandbox instead of production variables")
	viper.BindPFlag("root.dev", RootCmd.PersistentFlags().Lookup("dev"))
	viper.BindPFlag("root.verbose", RootCmd.PersistentFlags().Lookup("verbose"))
}

const (
	subifyConfigPath = "$HOME"
	subifyConfigFile = ".subify"
)

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if config.ConfigFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(config.ConfigFile)
	}

	viper.SetConfigName(subifyConfigFile) // name of config file (without extension)
	viper.AddConfigPath(subifyConfigPath) // adding home directory as first search path

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err == nil {
		fmt.Println("Using config file:" + viper.ConfigFileUsed())
	}
}
