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
	"github.com/cheld/miniprow/pkg/common/config"
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
		piperCfg, _ := cmd.Flags().GetString("config-piper")
		boskosCfg, _ := cmd.Flags().GetString("config-boskos")
		secret, _ := cmd.Flags().GetString("secret")
		bindaddr, _ := cmd.Flags().GetString("bind-addr")
		port, _ := cmd.Flags().GetInt("port")
		envSettings, _ := cmd.Flags().GetStringToString("env")
		logSetting, _ := cmd.Flags().GetString("log-level")
		setLogLevel(logSetting)

		// print version info
		logrus.Infof("Version: %s, commit %s\n", info.Version, info.Commit)

		// read environment flags
		util.Environment.Value("PORT").Update(&port)

		// find config files
		piperCfg = config.FindFile(piperCfg, "piper.yaml")
		logrus.Infof("Piper config found at path %s\n", piperCfg)
		boskosCfg = config.FindFile(boskosCfg, "boskos.yaml")
		logrus.Infof("Boskos config found at path %s\n", boskosCfg)

		// Register http endpoints
		mux := http.NewServeMux()
		boskosServer.Register(mux, boskosCfg)
		piperServer.Register(mux, piperCfg, envSettings, secret)
		commonServer.Register(mux)

		// Start server
		addr := fmt.Sprintf("%s:%d", bindaddr, port)
		logrus.Infof("Starting server with %s\n", addr)
		server := &http.Server{
			Handler: mux,
			Addr:    addr,
		}
		err := server.ListenAndServe()
		fmt.Println(err)
	},
}

func command() *cobra.Command {
	rootCmd.Flags().IntP("port", "p", 3000, "Port for the HTTP endpoint")
	rootCmd.Flags().StringP("bind-addr", "", "0.0.0.0", "the bind addr of the server")
	rootCmd.Flags().StringP("secret", "s", "", "Protect access to the webhook")
	rootCmd.Flags().StringToStringP("env", "e", nil, "Provide environment variables that can be accessed by event handlers")
	rootCmd.Flags().StringP("config-boskos", "", "", "config file for boskos (default is $HOME/.piper.yaml)")
	rootCmd.Flags().StringP("config-piper", "", "", "config file for piper (default is $HOME/.boskos.yaml)")
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
