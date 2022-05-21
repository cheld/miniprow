package main

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
//package boskosctl

// // eventCmd represents the event command
// var eventCmd = &cobra.Command{
// 	Use:   "event [payload]",
// 	Short: "Manually fire an event over command line",
// 	Long: `A longer description that spans multiple lines and likely contains examples
// and usage of using your command. For example:

// Cobra is a CLI library for Go that empowers applications.
// This application is a tool to generate the needed files
// to quickly create a Cobra application.`,
// 	Run: func(cmd *cobra.Command, args []string) {

// 		// parse cli
// 		cfgFile, _ := cmd.Flags().GetString("config")
// 		overrideVariables, _ := cmd.Flags().GetStringToString("env")
// 		eventType, _ := cmd.Flags().GetString("type")

// 		// prepare configuration
// 		payload := strings.Join(args, " ")
// 		payloadFile, _ := cmd.Flags().GetString("file")
// 		if payloadFile == "-" {
// 			payload = readStdIn()
// 		} else if payloadFile != "" {
// 			payload = readFile(payloadFile)
// 		}
// 		cfg, err := config.Load(cfgFile)
// 		if err != nil {
// 			fmt.Println(err)
// 			os.Exit(1)
// 		}
// 		env := config.Environ(overrideVariables)

// 		// execute event
// 		handler := event.NewHandler(cfg, env)
// 		dispatcher := trigger.NewDispatcher(cfg)
// 		dispatcher.Execute(handler.HandleCli(eventType, payload))
// 		fmt.Println()
// 	},
// }

// func readStdIn() string {
// 	reader := bufio.NewReader(os.Stdin)
// 	var chars []rune
// 	for {
// 		c, _, err := reader.ReadRune()
// 		if err != nil && err == io.EOF {
// 			break
// 		}
// 		chars = append(chars, c)
// 	}
// 	return string(chars)
// }

// func readFile(filename string) string {
// 	filecontent, err := ioutil.ReadFile(filename)
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}
// 	return string(filecontent)
// }

// func init() {
// 	rootCmd.AddCommand(eventCmd)
// 	eventCmd.Flags().StringToStringP("env", "e", nil, "provide environment variables that can be accessed by event handlers")
// 	eventCmd.Flags().StringP("config", "c", "", "config file (default is $HOME/.piper.yaml)")
// 	eventCmd.Flags().StringP("file", "f", "", "read event payload from a file (use \"-f -\" for stdin)")
// 	eventCmd.Flags().StringP("type", "t", "", "The type of event. It must match the event configuration")
// 	eventCmd.MarkFlagRequired("type")
// }
