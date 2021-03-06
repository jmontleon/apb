//
// Copyright (c) 2018 Red Hat, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package cmd

import (
	"os"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Verbose controls the logging level, when enabled will set level to debug
var Verbose bool

// CfgFile is the configuration file
var CfgFile string

var rootCmd = &cobra.Command{
	Use:   "apb",
	Short: "Tool for working with Ansible Playbook Bundles",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if Verbose {
			log.SetLevel(log.DebugLevel)
		}
	},
}

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{DisableTimestamp: true})

	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVar(&CfgFile, "config", "", "configuration file (default is $HOME/.apb)")
}

func initConfig() {
	viper.SetConfigType("json")
	if CfgFile != "" {
		viper.SetConfigFile(CfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(".apb")
		filePath := home + "/.apb.json"
		if err := viper.ReadInConfig(); err != nil {
			log.Warning("Didn't find config file, creating one.")
			file, err := os.Create(filePath)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
			file.WriteString("{}")
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Error("Can't read config: ", err)
		os.Exit(1)
	}
}

// Execute invokes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
