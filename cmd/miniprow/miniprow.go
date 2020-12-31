/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"fmt"
	"net/http"
	"os"

	boskosServer "github.com/cheld/miniprow/pkg/boskos/server"
	"github.com/cheld/miniprow/pkg/common/info"
	commonServer "github.com/cheld/miniprow/pkg/common/server"
	"github.com/cheld/miniprow/pkg/common/util"
	piperServer "github.com/cheld/miniprow/pkg/piper/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	if err := command().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "miniprow",
	Short: "Mini Prow helps to automate tasks for typical development projects.",
	Long: `Mini Prow helps to automate tasks for typical development projects. 
Tools are integrated using webhooks, flexible rules and triggers. Examples
for automation tasks are job execution, GitHub/Gitlab policy enforcement,
chat-ops via /foo style commands and Slack notifications.`,

	Run: func(cmd *cobra.Command, args []string) {

		// read cli flags
		piperFileName, _ := cmd.Flags().GetString("config-piper")
		boskosFileName, _ := cmd.Flags().GetString("config-boskos")
		secret, _ := cmd.Flags().GetString("secret")
		bindaddr, _ := cmd.Flags().GetString("bind-addr")
		port, _ := cmd.Flags().GetInt("port")
		logSetting, _ := cmd.Flags().GetString("log-level")
		setLogLevel(logSetting)

		// read environment flags
		util.Environment.Value("PORT").Update(&port)

		// read config files
		piperCfgFile := util.FindExistingFile(util.DefaultConfigLocations(piperFileName))
		piperCfg, err := util.ReadConfiguration(piperCfgFile, "PIPER_CONFIG")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		boskosCfgFile := util.FindExistingFile(util.DefaultConfigLocations(boskosFileName))
		boskosCfg, err := util.ReadConfiguration(boskosCfgFile, "BOSKOS_CONFIG")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// print version info
		logrus.Infof("Version: %s, commit %s\n", info.Version, info.Commit)

		// Register http endpoints
		mux := http.NewServeMux()
		mux.Handle("/piper/", piperServer.NewPiper(piperCfg, secret))
		mux.Handle("/boskos/", boskosServer.NewBoskos(boskosCfg))
		mux.Handle("/common/", commonServer.NewCommon())

		// Start server
		addr := fmt.Sprintf("%s:%d", bindaddr, port)
		logrus.Infof("Starting server with %s\n", addr)
		server := &http.Server{
			Handler: mux,
			Addr:    addr,
		}
		err = server.ListenAndServe()
		fmt.Println(err)
	},
}

func command() *cobra.Command {
	rootCmd.Flags().IntP("port", "p", 3000, "Port for the HTTP endpoint")
	rootCmd.Flags().StringP("bind-addr", "", "0.0.0.0", "the bind addr of the server")
	rootCmd.Flags().StringP("secret", "s", "", "Protect access to the webhook")
	rootCmd.Flags().StringP("config-boskos", "", "boskos.yaml", "config file for boskos")
	rootCmd.Flags().StringP("config-piper", "", "piper.yaml", "config file for piper")
	rootCmd.PersistentFlags().StringP("log-level", "l", "DEBUG", "set debug log level")
	return rootCmd
}

func setLogLevel(logSetting string) {
	logLevel, err := logrus.ParseLevel(logSetting)
	if err != nil {
		fmt.Printf("Log level %s cannot be set\n", logSetting)
		os.Exit(1)
	}
	logrus.SetLevel(logLevel)
}
