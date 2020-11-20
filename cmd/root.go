package cmd

import (
	goflag "flag"
	"fmt"
	"os"

	pflag "github.com/spf13/pflag"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cicd-bot",
	Short: "CI/CD bot helps to automate tasks for typical development projects.",
	Long: `CI/CD bot helps to automate tasks for typical development projects. 
Tools are integrated using webhooks, flexible rules and triggers. Examples
for automation tasks are job execution, GitHub/Gitlab policy enforcement,
chat-ops via /foo style commands and Slack notifications.`,
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
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	goflag.Parse()
}
