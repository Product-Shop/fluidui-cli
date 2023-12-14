/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package flutter

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/spf13/cobra"
)

var (
	projectName string
)

// flutterCmd represents the flutter command
var FlutterCmd = &cobra.Command{
	Use:   "flutter",
	Short: "A brief description of your command",
	Long:  `A Flutter CLI for making dope shit`,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Hugo",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		flutterCLI("--version", "")
	},
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Print the version number of Hugo",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		flutterCLI("create", projectName)
	},
}

func flutterCLI(command string, args string) {
	cmd := exec.Command("flutter", command, args)
	b, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Flutter function failed: %v", err)
	}
	fmt.Printf("%s\n", b)
}

func init() {
	FlutterCmd.AddCommand(versionCmd)
	FlutterCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&projectName, "name", "n", "", "The name of the project")

	if err := createCmd.MarkFlagRequired("name"); err != nil {
		fmt.Println(err)
	}
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// flutterCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// flutterCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
