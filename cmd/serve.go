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
package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/cheld/cicd-bot/pkg/piper/config"
	"github.com/cheld/cicd-bot/pkg/webhook"
	"github.com/spf13/cobra"
)

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

		// parse cli
		cfgFile, _ := cmd.Flags().GetString("config")
		secret, _ := cmd.Flags().GetString("secret")
		bindaddr, _ := cmd.Flags().GetString("bind-addr")
		port, _ := cmd.Flags().GetInt("port")
		overrideVariables, _ := cmd.Flags().GetStringToString("env")

		// prepare configuration
		cfg, err := config.Load(cfgFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		env := config.Environ(overrideVariables)
		if env["PORT"] != "" {
			p, _ := strconv.Atoi(env["PORT"])
			if p != 0 {
				port = p
			}
		}
		opts := webhook.Options{
			Port:     port,
			Secret:   secret,
			Bindaddr: bindaddr,
		}

		// start webhook
		webhook.Run(cfg, env, opts)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntP("port", "p", 3000, "Port for the HTTP endpoint")
	serveCmd.Flags().StringP("bind-addr", "", "127.0.0.1", "the bind addr of the server")
	serveCmd.Flags().StringP("secret", "s", "", "Protect access to the webhook")
	serveCmd.Flags().StringToStringP("env", "e", nil, "Provide environment variables that can be accessed by event handlers")
	serveCmd.Flags().StringP("config", "c", "", "config file (default is $HOME/.cicd-bot.yaml)")
}
