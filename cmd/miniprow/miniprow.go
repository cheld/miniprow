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
	"strconv"
	"strings"

	boskosServer "github.com/cheld/miniprow/pkg/boskos/server"
	"github.com/cheld/miniprow/pkg/common/config"
	commonServer "github.com/cheld/miniprow/pkg/common/server"
	piperServer "github.com/cheld/miniprow/pkg/piper/server"
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
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Launch a HTTP server to provide a webhook integration endpoint",
	Long: `Launch a HTTP server to provide a webhook integration endpoint.

The following endpoints are available:
http://<localhost:port>/webhook/github
http://<localhost:port>/webhook/gitlab
http://<localhost:port>/webhook/http`,
	Run: func(cmd *cobra.Command, args []string) {

		// read cli flags
		piperCfg, _ := cmd.Flags().GetString("piper-config")
		boskosCfg, _ := cmd.Flags().GetString("boskos-config")
		secret, _ := cmd.Flags().GetString("secret")
		bindaddr, _ := cmd.Flags().GetString("bind-addr")
		port, _ := cmd.Flags().GetInt("port")
		envSettings, _ := cmd.Flags().GetStringToString("env")

		// read environment variables
		for _, entry := range os.Environ() {
			keyValue := strings.Split(entry, "=")
			if keyValue[0] == "PORT" {
				p, _ := strconv.Atoi(keyValue[0])
				if p != 0 {
					port = p
				}
			}
		}

		// find config files
		piperCfg = config.FindFile(piperCfg, "piper.yaml")
		boskosCfg = config.FindFile(boskosCfg, "boskos.yaml")

		// Register http endpoints
		mux := http.NewServeMux()
		boskosServer.Register(mux, boskosCfg)
		piperServer.Register(mux, piperCfg, envSettings, secret)
		commonServer.Register(mux)

		// Start server
		server := &http.Server{
			Handler: mux,
			Addr:    fmt.Sprintf("%s:%d", bindaddr, port),
		}
		err := server.ListenAndServe()
		fmt.Println(err)
	},
}

func command() *cobra.Command {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntP("port", "p", 3000, "Port for the HTTP endpoint")
	serveCmd.Flags().StringP("bind-addr", "", "0.0.0.0", "the bind addr of the server")
	serveCmd.Flags().StringP("secret", "s", "", "Protect access to the webhook")
	serveCmd.Flags().StringToStringP("env", "e", nil, "Provide environment variables that can be accessed by event handlers")
	serveCmd.Flags().StringP("piper-config", "", "", "config file (default is $HOME/.piper.yaml)")
	serveCmd.Flags().StringP("boskos-config", "", "", "config file (default is $HOME/.boskos.yaml)")
	return rootCmd
}
