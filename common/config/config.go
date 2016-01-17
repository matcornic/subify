package config

var (
	// Verbose conditions the quantity of output of this tool
	Verbose bool

	// ConfigFile is $HOME/.subify(.json) by default
	ConfigFile string

	// Dev mode to use sandbox parameters
	Dev bool
)
