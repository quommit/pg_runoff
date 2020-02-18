/*
Copyright © 2020 José Tomás Navarro Carrión <jt.navarro@ua.es>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
  "fmt"
  "github.com/spf13/cobra"
  "github.com/spf13/viper"
  "os"
  homedir "github.com/mitchellh/go-homedir"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
  Use:   "runoff",
  Short: "pg_runoff CLI",
  Long: `Command Line Interface utilities
to assist users in managing and computing
runoff models using the PostgreSQL
pg_runoff extension. Please visit
https://github.com/siose-innova/pg_runoff
for pg_runoff CLI workflow examples.
To get help on runoff subcommands run:

runoff --help`,
  Version: os.Getenv("EXTENSION_VERSION"),
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
  rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.runoff.yaml)")
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

    // Search config in home directory with name ".runoff" (without extension).
    viper.AddConfigPath(home)
    viper.SetConfigName(".runoff")
  }

  viper.AutomaticEnv() // read in environment variables that match

  // If a config file is found, read it in.
  if err := viper.ReadInConfig(); err == nil {
    fmt.Println("Using config file:", viper.ConfigFileUsed())
  }
}
