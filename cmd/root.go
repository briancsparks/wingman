package cmd

/* Copyright © 2022 Brian C Sparks <briancsparks@gmail.com> -- MIT (see LICENSE file) */

import (
  "fmt"
  "github.com/spf13/cobra"
  "github.com/spf13/pflag"
  "os"
  "strings"

  homedir "github.com/mitchellh/go-homedir"
  "github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
  Use:   "wingman",
  Short: "A brief description of your application",
  Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
  // Uncomment the following line if your bare application
  // has an action associated with it:
  //	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func init() {
  cobra.OnInitialize(initConfig)

  // Here you will define your flags and configuration settings.
  // Cobra supports persistent flags, which, if defined here,
  // will be global for your application.

  rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.wingman.yaml)")

  // Cobra also supports local flags, which will only run
  // when this action is called directly.
  //rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
  if cfgFile != "" {
    // Use config file from the flag.
    viper.SetConfigFile(cfgFile)
  } else {
    // Find home directory.
    home, err := homedir.Dir()
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }

    // Search config in home directory with name ".wingman" (without extension).
    viper.AddConfigPath(home)
    viper.SetConfigName(".wingman")
  }

  viper.SetEnvPrefix("WINGMAN")
  viper.AutomaticEnv() // read in environment variables that match

  bindFlags(rootCmd)

  // If a config file is found, read it in.
  if err := viper.ReadInConfig(); err == nil {
    fmt.Println("Using config file:", viper.ConfigFileUsed())
  }
}

// Stolen from: https://github.com/carolynvs/stingoftheviper/blob/main/main.go
// Bind each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command /*, v *viper.Viper*/) {
  cmd.Flags().VisitAll(func(f *pflag.Flag) {

    // Environment variables can't have dashes in them, so bind them to their equivalent
    // keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
    if strings.Contains(f.Name, "-") {
      envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
      viper.BindEnv(f.Name, fmt.Sprintf("%s_%s", "KRONK", envVarSuffix))
    }

    // Apply the viper config value to the flag when the flag is not set and viper has a value
    if !f.Changed && viper.IsSet(f.Name) {
      val := viper.Get(f.Name)
      cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
    }
  })
}
